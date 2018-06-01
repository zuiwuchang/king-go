package queue

import (
	"container/list"
)

// _StepQueue 按照 指定 步長增長的 隊列
type _StepQueue struct {
	// 緩存 節點
	cacheStep IQueue
	items     *list.List
	step      int

	length, capacity int
}

// NewStepQueue 返回一個 按照 指定 步長增長的 隊列
func NewStepQueue(step int, cache int) (q IQueue, e error) {
	if step < 1 {
		e = ErrQueueCap
	} else {
		var cacheStep IQueue
		if cache > 0 {
			cacheStep, e = NewStatic(cache)
			if e != nil {
				return
			}
		}
		q = &_StepQueue{
			cacheStep: cacheStep,
			step:      step,
			items:     list.New(),
		}
	}
	return
}
func (q *_StepQueue) createElement() (node IQueue, e error) {
	if q.cacheStep == nil {
		//沒有 緩存 創建 新節點
		node, e = NewStatic(q.step)
	} else {
		if q.cacheStep.Len() == 0 {
			//沒有 緩存 創建 新節點
			node, e = NewStatic(q.step)
		} else {
			// 返回 緩存
			var val interface{}
			val, e = node.PopBack()
			if e == nil {
				node = val.(IQueue)
			}
		}
	}
	return
}

// PushBack 壓入 隊列 尾 失敗 通常返回 ErrQueueOverflow
func (q *_StepQueue) PushBack(val interface{}) (e error) {
	// 獲取 幾點
	var node IQueue
	ele := q.items.Back()
	if ele == nil {
		// 沒有 節點 創建 新 節點
		node, e = q.createElement()
		if e != nil {
			return
		}
		e = node.PushBack(val)
		if e != nil {
			return
		}
		// 壓入 新節點
		q.items.PushBack(node)
		q.capacity += node.Cap()
	} else {
		node = ele.Value.(IQueue)
		e = node.PushBack(val)
		if e == ErrQueueOverflow {
			// 隊列 已滿 創建 新 節點
			node, e = q.createElement()
			if e != nil {
				return
			}
			e = node.PushBack(val)
			if e != nil {
				return
			}
			// 壓入 新節點
			q.items.PushBack(node)
			q.capacity += node.Cap()
		}
	}
	q.length++
	return
}

// PushFront 壓入 隊列 頭 失敗 通常返回 ErrQueueOverflow
func (q *_StepQueue) PushFront(val interface{}) (e error) {
	// 獲取 幾點
	var node IQueue
	ele := q.items.Front()
	if ele == nil {
		// 沒有 節點 創建 新 節點
		node, e = q.createElement()
		if e != nil {
			return
		}
		e = node.PushFront(val)
		if e != nil {
			return
		}
		// 壓入 新節點
		q.items.PushFront(node)
		q.capacity += node.Cap()
	} else {
		node = ele.Value.(IQueue)
		e = node.PushFront(val)
		if e == ErrQueueOverflow {
			// 隊列 已滿 創建 新 節點
			node, e = q.createElement()
			if e != nil {
				return
			}
			e = node.PushFront(val)
			if e != nil {
				return
			}
			// 壓入 新節點
			q.items.PushFront(node)
			q.capacity += node.Cap()
		}
	}
	q.length++
	return
}

// PopBack 從 隊列 尾 出棧 如果為空 返回 nil,ErrQueueEmpty
func (q *_StepQueue) PopBack() (val interface{}, e error) {
	if q.length == 0 {
		e = ErrQueueEmpty
		return
	}
	ele := q.items.Back()
	node := ele.Value.(IQueue)
	val, _ = node.PopBack()
	q.length--
	if node.Len() == 0 {
		q.items.Remove(ele)
		q.freeNode(node)
	}
	return
}

// PopFront 從 隊列 頭 出棧 如果為空 返回 nil,ErrQueueEmpty
func (q *_StepQueue) PopFront() (val interface{}, e error) {
	if q.length == 0 {
		e = ErrQueueEmpty
		return
	}
	ele := q.items.Front()
	node := ele.Value.(IQueue)
	val, _ = node.PopFront()
	q.length--
	if node.Len() == 0 {
		q.items.Remove(ele)
		q.freeNode(node)
	}
	return
}
func (q *_StepQueue) freeNode(node IQueue) {
	cacheStep := q.cacheStep
	if cacheStep == nil {
		return
	} else if cacheStep.Len() < cacheStep.Cap() {
		cacheStep.PushBack(node)
	}
}

// Back 返回 隊列 尾 如果為空 返回 nil,ErrQueueEmpty
func (q *_StepQueue) Back() (val interface{}, e error) {
	if q.length == 0 {
		e = ErrQueueEmpty
		return
	}
	ele := q.items.Back()
	node := ele.Value.(IQueue)
	val, _ = node.Back()
	return
}

// Front 返回 隊列 頭 如果為空 返回 nil,ErrQueueEmpty
func (q *_StepQueue) Front() (val interface{}, e error) {
	if q.length == 0 {
		e = ErrQueueEmpty
		return
	}
	ele := q.items.Front()
	node := ele.Value.(IQueue)
	val, _ = node.Front()
	return
}

// Cap 返回 隊列 容量
func (q *_StepQueue) Cap() int {
	return q.capacity
}

// Len 返回 隊列 大小
func (q *_StepQueue) Len() int {
	return q.length
}

// Reset 重置 隊列
func (q *_StepQueue) Reset() {
	cacheStep := q.cacheStep
	if cacheStep != nil && cacheStep.Len() != cacheStep.Cap() {
		for ele := q.items.Front(); ele != nil && cacheStep.Len() != cacheStep.Cap(); ele = ele.Next() {
			node := ele.Value.(IQueue)
			node.Reset()
			cacheStep.PushBack(node)
		}
	}
	q.items.Init()
	q.length = 0
	q.capacity = 0
	return
}
