package lru

import (
	"sync"
)

type _Element struct {
	Next *_Element
	Pre  *_Element

	Key   IKey
	Value IValue
}
type lruImpl struct {
	RW sync.RWMutex

	Max int

	//緩存 的節點
	Keys  map[IKey]*_Element
	Front *_Element
	Back  *_Element
}

func newLRUImpl(maxElementSize int) *lruImpl {
	if maxElementSize < 1 {
		maxElementSize = 1
	}
	return &lruImpl{
		Max: maxElementSize,

		Keys: make(map[IKey]*_Element),
	}
}

//返回 當前 緩存 量
func (this *lruImpl) Len() (n int) {
	this.RW.RLock()
	n = len(this.Keys)
	this.RW.RUnlock()
	return
}

//返回 緩存 最高容量
func (this *lruImpl) Cap() (n int) {
	this.RW.RLock()
	n = this.Max
	this.RW.RUnlock()
	return
}

//刪除 所有 緩存
func (this *lruImpl) Clear() {
	this.RW.Lock()
	for key, _ := range this.Keys {
		//刪除 map
		delete(this.Keys, key)
	}
	this.Front = nil
	this.Back = nil
	this.RW.Unlock()
}

//刪除 指定緩存
func (this *lruImpl) Delete(key IKey) {
	this.RW.Lock()
	this.unsafeDelete(key)
	this.RW.Unlock()
}
func (this *lruImpl) unsafeDelete(key IKey) {
	ele, ok := this.Keys[key]
	//緩存 不存在 直接 返回
	if !ok {
		return
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

//返回 是否存在 緩存
func (this *lruImpl) Ok(key IKey) (ok bool) {
	this.RW.RLock()
	_, ok = this.Keys[key]
	this.RW.RUnlock()
	return
}

//返回 緩存值 不存在 返回 nil
func (this *lruImpl) Get(key IKey) IValue {
	this.RW.RLock()
	defer this.RW.RUnlock()

	ele, ok := this.Keys[key]
	//緩存 不存在 直接 返回
	if !ok {
		return nil
	}

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

//創建 一個 緩存
func (this *lruImpl) Set(key IKey, val IValue) {
	this.RW.Lock()
	defer this.RW.Unlock()

	//驗證 存在
	ele, ok := this.Keys[key]
	if ok {
		//更新
		ele.Value = val

		//移動到 Back
		this.unsafeToBack(ele)
	} else {
		//創建 新緩存
		if len(this.Keys) == this.Max &&
			this.Front != nil {

			//刪除 front
			delete(this.Keys, this.Front.Key)

			this.unsafeRemoveList(this.Front)
		}
		//創建
		this.unsafeNew(key, val)
	}
}

func (this *lruImpl) unsafeNew(key IKey, val IValue) {
	ele := &_Element{
		Pre:  this.Back,
		Next: nil,

		Key:   key,
		Value: val,
	}

	this.Keys[key] = ele
	if this.Back == nil {
		this.Front = ele
		this.Back = ele
	} else {
		this.Back.Next = ele
	}
}
