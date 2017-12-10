package lru2

import (
	"fmt"
	"sync"
	"testing"
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
	this.unsafeClear()
	this.Mutex.Unlock()
}
func (this *lruImpl) unsafeClear() {
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
	if ele.Pre == nil {
		this.Front = ele.Next
	} else {
		ele.Pre.Next = ele.Next
	}

	ele.Next = nil
	ele.Pre = this.Back
	this.Back.Next = ele
	this.Back = ele
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
	} else {
		this.Back.Next = ele
	}
	this.Back = ele

	go this.onExpired(ele, ele.Timer)
}

//釋放 緩存並返回 Len()
//
//執行後 緩存容量將 <= Cap() * percentage
func (this *lruImpl) Resize(percentage float64) int {
	this.Mutex.Lock()
	defer this.Mutex.Unlock()

	max := (int)((float64)(this.Max) * percentage)
	if max == this.Max {
		return len(this.Keys)
	} else if max == 0 {
		if len(this.Keys) != 0 {
			this.unsafeClear()
		}
		return 0
	} else if max == 1 {
		this.unsafeOnlyOne()
		return len(this.Keys)
	}

	for len(this.Keys) > max {
		this.unsafePopFront()
	}
	return len(this.Keys)
}
func (this *lruImpl) unsafeOnlyOne() {
	ele := this.Back
	if ele == nil { //沒有節點
		return
	} else if ele == this.Front { //只有一個節點
		return
	}

	ele.Pre = nil
	this.Front = ele
	for k, v := range this.Keys {
		if k != ele.Key {
			delete(this.Keys, k)

			//停止 定時器
			if v.Timer != nil {
				v.Timer.Stop()
				v.Timer = nil
			}
		}
	}
}
func (this *lruImpl) unsafePopFront() {
	ele := this.Front
	this.Front = ele.Next
	if ele.Next == nil {
		this.Back = nil
	} else {
		ele.Next.Pre = nil
	}
	delete(this.Keys, ele.Key)
	//停止 定時器
	if ele.Timer != nil {
		ele.Timer.Stop()
		ele.Timer = nil
	}
}
func (this *lruImpl) debugPrint(t *testing.T) {
	node := this.Front
	if t == nil {
		fmt.Print("[")
	}
	sum := 0
	for node != nil {
		sum++
		if t == nil {
			fmt.Printf("%v=%v ,", node.Key, node.Value)
		}
		next := node.Next
		if next != nil {
			if next.Pre != node {
				t.Fatal("bad Pre")
			}
		}

		node = next
	}
	if t == nil {
		fmt.Print("]")
	}
	if len(this.Keys) != sum {
		t.Fatalf(" len(%v %v) ", len(this.Keys), sum)
	}
	if this.Back != nil {
		if t == nil {
			fmt.Printf(" %v=%v \n", this.Back.Key, this.Back.Value)
		}
	}
}
