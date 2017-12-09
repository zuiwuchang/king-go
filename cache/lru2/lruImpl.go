package lru2

import (
	"sync"
	"time"
)

type _Element struct {
	Next *_Element
	Pre  *_Element

	Key   IKey
	Value IValue

	Timer *time.Timer
}
type lruImpl struct {
	Mutex sync.Mutex

	Expired time.Duration
	Max     int

	//緩存 的節點
	Keys  map[IKey]*_Element
	Front *_Element
	Back  *_Element
}

func newLRUImpl(expired time.Duration, maxElementSize int) *lruImpl {
	if maxElementSize < 1 {
		maxElementSize = 1
	}
	return &lruImpl{
		Expired: expired,
		Max:     maxElementSize,

		Keys: make(map[IKey]*_Element),
	}
}

//返回 當前 緩存 量
func (this *lruImpl) Len() (n int) {
	this.Mutex.Lock()
	n = len(this.Keys)
	this.Mutex.Unlock()
	return
}

//返回 緩存 最高容量
func (this *lruImpl) Cap() (n int) {
	this.Mutex.Lock()
	n = this.Max
	this.Mutex.Unlock()
	return
}

//刪除 所有 緩存
func (this *lruImpl) Clear() {
	this.Mutex.Lock()
	for key, ele := range this.Keys {
		//停止 定時器
		if ele.Timer != nil {
			ele.Timer.Stop()
			ele.Timer = nil
		}

		//刪除 map
		delete(this.Keys, key)
	}
	this.Front = nil
	this.Back = nil
	this.Mutex.Unlock()
}

//刪除 指定緩存
func (this *lruImpl) Delete(key IKey) {
	this.Mutex.Lock()
	this.unsafeDelete(key)
	this.Mutex.Unlock()
}
func (this *lruImpl) unsafeDelete(key IKey) {
	ele, ok := this.Keys[key]
	//緩存 不存在 直接 返回
	if !ok {
		return
	}

	//停止 定時器
	if ele.Timer != nil {
		ele.Timer.Stop()
		ele.Timer = nil
	}
	//刪除 map
	delete(this.Keys, key)

	//刪除 鏈表
	this.unsafeRemoveList(ele)
}
func (this *lruImpl) unsafeRemoveList(ele *_Element) {
	if ele.Next == nil {
		this.Back = ele.Pre
	} else { //需要 設置 next
		ele.Next.Pre = ele.Pre
		if ele.Pre == nil {
			this.Front = ele.Next
		} else {
			ele.Pre.Next = ele.Next
		}
	}
	if ele.Pre == nil {
		this.Front = ele.Next
	} else { //需要 設置 pre
		ele.Pre.Next = ele.Next
		if ele.Next == nil {
			this.Back = ele.Pre
		} else {
			ele.Next.Pre = ele.Pre
		}
	}
}

//返回 是否存在 緩存 不會更新超時時間
func (this *lruImpl) Ok(key IKey) (ok bool) {
	this.Mutex.Lock()
	_, ok = this.Keys[key]
	this.Mutex.Unlock()
	return
}

//返回 緩存值 不存在 返回 nil
func (this *lruImpl) Get(key IKey) IValue {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	ele, ok := this.Keys[key]
	//緩存 不存在 直接 返回
	if !ok {
		return nil
	}
	//重置 超時 定時器
	if ele.Timer != nil {
		ele.Timer.Stop()
	}
	ele.Timer = time.NewTimer(this.Expired)
	go this.onExpired(ele, ele.Timer)

	//移動到 Back
	this.unsafeToBack(ele)

	return ele.Value
}
func (this *lruImpl) unsafeToBack(ele *_Element) {
	if ele.Next == nil {
		//本來就是 back 節點 直接返回
		return
	}
	ele.Next.Pre = ele.Pre
	if ele.Pre != nil {
		ele.Pre.Next = ele.Next
	}
}

func (this *lruImpl) onExpired(ele *_Element, timer *time.Timer) {
	//等待超時
	<-timer.C

	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	//驗證 定時器 是否 已經 過期
	if ele.Timer != timer {
		return
	}

	//刪除 緩存
	this.unsafeDelete(ele.Key)
}

//創建 一個 緩存
func (this *lruImpl) Set(key IKey, val IValue) {
	this.Set2(key, val, this.Expired)
}

//創建 一個 緩存 同時指定 超時 時間
func (this *lruImpl) Set2(key IKey, val IValue, expired time.Duration) {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	//驗證 存在
	ele, ok := this.Keys[key]
	if ok {
		//更新
		ele.Value = val
		this.unsafeUpdate(ele, expired)
	} else {
		//創建 新緩存
		if len(this.Keys) == this.Max &&
			this.Front != nil {

			//停止 定時器
			if this.Front.Timer != nil {
				this.Front.Timer.Stop()
				this.Front.Timer = nil
			}

			//刪除 front
			delete(this.Keys, this.Front.Key)

			this.unsafeRemoveList(this.Front)
		}
		//創建
		this.unsafeNew(key, val, expired)
	}
}
func (this *lruImpl) unsafeUpdate(ele *_Element, expired time.Duration) {
	//重置 超時 定時器
	if ele.Timer != nil {
		ele.Timer.Stop()
	}
	ele.Timer = time.NewTimer(expired)
	go this.onExpired(ele, ele.Timer)

	//移動到 Back
	this.unsafeToBack(ele)
}
func (this *lruImpl) unsafeNew(key IKey, val IValue, expired time.Duration) {
	ele := &_Element{
		Pre:  this.Back,
		Next: nil,

		Key:   key,
		Value: val,

		Timer: time.NewTimer(expired),
	}

	this.Keys[key] = ele
	if this.Back == nil {
		this.Front = ele
		this.Back = ele
	} else {
		this.Back.Next = ele
	}

	go this.onExpired(ele, ele.Timer)
}
