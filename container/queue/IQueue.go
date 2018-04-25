package queue

import (
	"errors"
)

// ErrQueueCap 初始化 容量太小
var ErrQueueCap = errors.New("queue capacity must large than 0")

// ErrQueueOverflow 容量容量已慢
var ErrQueueOverflow = errors.New("queue capacity overflow")

// ErrQueueEmpty .
var ErrQueueEmpty = errors.New("queue empty")

// IQueue 隊列 定義
type IQueue interface {
	// PushBack 壓入 隊列 尾 失敗 通常返回 ErrQueueOverflow
	PushBack(val interface{}) (e error)
	// PushFront 壓入 隊列 頭 失敗 通常返回 ErrQueueOverflow
	PushFront(val interface{}) (e error)

	// PopBack 從 隊列 尾 出棧 如果為空 返回 nil,ErrQueueEmpty
	PopBack() (val interface{}, e error)
	// PopFront 從 隊列 頭 出棧 如果為空 返回 nil,ErrQueueEmpty
	PopFront() (val interface{}, e error)

	// Back 返回 隊列 尾 如果為空 返回 nil,ErrQueueEmpty
	Back() (val interface{}, e error)
	// Front 返回 隊列 頭 如果為空 返回 nil,ErrQueueEmpty
	Front() (val interface{}, e error)

	// Cap 返回 隊列 容量
	Cap() int
	// Len 返回 隊列 大小
	Len() int

	// Reset 重置 隊列
	Reset()
}
