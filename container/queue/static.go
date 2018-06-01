package queue

// _StaticQueue 容量 固定的 隊列
type _StaticQueue struct {
	// 隊列 數據
	data []interface{}
	// 當前 大小
	length int
	// 起點位置
	pos int
}

// NewStatic 返回一個 容量固定的 隊列
func NewStatic(capacity int) (q IQueue, e error) {
	if capacity < 1 {
		e = ErrQueueCap
	} else {
		q = &_StaticQueue{
			data: make([]interface{}, capacity),
		}
	}
	return
}

// PushBack 壓入 隊列 尾 失敗 通常返回 ErrQueueOverflow
func (q *_StaticQueue) PushBack(val interface{}) (e error) {
	if q.length == len(q.data) {
		e = ErrQueueOverflow
		return
	}

	i := q.pos + q.length
	if i >= len(q.data) {
		i -= len(q.data)
	}
	q.data[i] = val
	q.length++
	return
}

// PushFront 壓入 隊列 頭 失敗 通常返回 ErrQueueOverflow
func (q *_StaticQueue) PushFront(val interface{}) (e error) {
	if q.length == len(q.data) {
		e = ErrQueueOverflow
		return
	}
	q.length++
	if q.pos == 0 {
		q.pos = len(q.data) - 1
	} else {
		q.pos--
	}
	q.data[q.pos] = val
	return
}

// PopBack 從 隊列 尾 出棧 如果為空 返回 nil,ErrQueueEmpty
func (q *_StaticQueue) PopBack() (val interface{}, e error) {
	if q.length == 0 {
		e = ErrQueueEmpty
		return
	}

	i := q.pos + q.length - 1
	if i >= len(q.data) {
		i -= len(q.data)
	}
	val = q.data[i]
	q.data[i] = nil

	if q.length == 1 {
		q.length = 0
		q.pos = 0
	} else {
		q.length--
	}
	return
}

// PopFront 從 隊列 頭 出棧 如果為空 返回 nil,ErrQueueEmpty
func (q *_StaticQueue) PopFront() (val interface{}, e error) {
	if q.length == 0 {
		e = ErrQueueEmpty
		return
	}
	val = q.data[q.pos]
	q.data[q.pos] = nil

	if q.length == 1 {
		q.length = 0
		q.pos = 0
	} else {
		q.length--
		q.pos++
		if q.pos == len(q.data) {
			q.pos = 0
		}
	}
	return
}

// Back 返回 隊列 尾 如果為空 返回 nil,ErrQueueEmpty
func (q *_StaticQueue) Back() (val interface{}, e error) {
	if q.length == 0 {
		e = ErrQueueEmpty
		return
	}

	i := q.pos + q.length - 1
	if i >= len(q.data) {
		i -= len(q.data)
	}
	val = q.data[i]

	return
}

// Front 返回 隊列 頭 如果為空 返回 nil,ErrQueueEmpty
func (q *_StaticQueue) Front() (val interface{}, e error) {
	if q.length == 0 {
		e = ErrQueueEmpty
		return
	}
	val = q.data[q.pos]

	return
}

// Cap 返回 隊列 容量
func (q *_StaticQueue) Cap() int {
	return len(q.data)
}

// Len 返回 隊列 大小
func (q *_StaticQueue) Len() int {
	return q.length
}

// Reset 重置 隊列
func (q *_StaticQueue) Reset() {
	if q.length == 0 {
		return
	}

	for i := 0; i < q.length; i++ {
		q.data[q.pos] = nil
		q.pos++
		if q.pos == len(q.data) {
			q.pos = 0
		}
	}
	q.pos = 0
	q.length = 0
}
