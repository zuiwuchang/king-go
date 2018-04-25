package queue

// IQueue 隊列 定義
type IQueue interface {
	// PushBack 壓入 隊列 尾
	PushBack(val interface{}) (e error)
	// PushFront 壓入 隊列 頭
	PushFront(val interface{}) (e error)

	// PopBack 從 隊列 尾 出棧 如果為空 返回 nil
	PopBack() (val interface{}, e error)
	// PopFront 從 隊列 頭 出棧 如果為空 返回 nil
	PopFront() (val interface{}, e error)

	// Cap 返回 隊列 容量
	Cap() int
	// Len 返回 隊列 大小
	Len() int

	// Reset 重置 隊列
	Reset()
}
