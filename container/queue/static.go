package queue

// staticQueue 容量 固定的 隊列
type staticQueue struct {
	// 隊列 數據
	data []interface{}
	// 當前 大小
	size int
	// 起點位置
	pos int
}

// NewStatic 返回一個 容量固定的 隊列
func NewStatic(capacity int) (q IQueue, e error) {
	if capacity < 1 {
		e = ErrQueueCap
	} else {
		q = &staticQueue{
			data: make([]interface{}, capacity),
		}
	}
	return
}

// PushBack 壓入 隊列 尾 失敗 通常返回 ErrQueueOverflow
func (q *staticQueue) PushBack(val interface{}) (e error) {
	if q.size == len(q.data) {
		e = ErrQueueOverflow
		return
	}

	i := q.pos + q.size
	if i >= len(q.data) {
		i -= len(q.data)
	}
	q.data[i] = val
	q.size++
	return
}

// PushFront 壓入 隊列 頭 失敗 通常返回 ErrQueueOverflow
func (q *staticQueue) PushFront(val interface{}) (e error) {
	if q.size == len(q.data) {
		e = ErrQueueOverflow
		return
	}
	q.size++
	if q.pos == 0 {
		q.pos = len(q.data) - 1
	} else {
		q.pos--
	}
	q.data[q.pos] = val
	return
}

// PopBack 從 隊列 尾 出棧 如果為空 返回 nil,ErrQueueEmpty
func (q *staticQueue) PopBack() (val interface{}, e error) {
	if q.size == 0 {
		e = ErrQueueEmpty
		return
	}

	i := q.pos + q.size - 1
	if i >= len(q.data) {
		i -= len(q.data)
	}
	val = q.data[i]

	if q.size == 1 {
		q.size = 0
		q.pos = 0
	} else {
		q.size--
	}
	return
}

// PopFront 從 隊列 頭 出棧 如果為空 返回 nil,ErrQueueEmpty
func (q *staticQueue) PopFront() (val interface{}, e error) {
	if q.size == 0 {
		e = ErrQueueEmpty
		return
	}
	val = q.data[q.pos]

	if q.size == 1 {
		q.size = 0
		q.pos = 0
	} else {
		q.size--
		q.pos++
		if q.pos == len(q.data) {
			q.pos = 0
		}
	}
	return
}

// Back 返回 隊列 尾 如果為空 返回 nil,ErrQueueEmpty
func (q *staticQueue) Back() (val interface{}, e error) {
	if q.size == 0 {
		e = ErrQueueEmpty
		return
	}

	i := q.pos + q.size - 1
	if i >= len(q.data) {
		i -= len(q.data)
	}
	val = q.data[i]

	return
}

// Front 返回 隊列 頭 如果為空 返回 nil,ErrQueueEmpty
func (q *staticQueue) Front() (val interface{}, e error) {
	if q.size == 0 {
		e = ErrQueueEmpty
		return
	}
	val = q.data[q.pos]

	return
}

// Cap 返回 隊列 容量
func (q *staticQueue) Cap() int {
	return len(q.data)
}

// Len 返回 隊列 大小
func (q *staticQueue) Len() int {
	return q.size
}

// Reset 重置 隊列
func (q *staticQueue) Reset() {
	if q.size == 0 {
		return
	}

	for i := 0; i < q.size; i++ {
		q.data[q.pos] = nil
		q.pos++
		if q.pos == len(q.data) {
			q.pos = 0
		}
	}
	q.pos = 0
	q.size = 0
}
