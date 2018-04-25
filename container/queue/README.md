# queue
queue 提供了一些 常用的 隊列

# IQueue 定義了 隊列的一般行為
```Go
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

	// Cap 返回 隊列 容量
	Cap() int
	// Len 返回 隊列 大小
	Len() int

	// Reset 重置 隊列
	Reset()
}
```

# NewStatic
NewStatic 用數組實現了一個 容量固定的 隊列
```Go
func NewStatic(capacity int) (q IQueue, e error) 
```
```Go
package main

import (
	"github.com/zuiwuchang/king-go/container/queue"
	"log"
)

func main() {
	q, e := queue.NewStatic(3)
	if e != nil {
		log.Fatalln(e)
	}
	q.PushBack(2)
	q.PushBack(3)
	q.PushFront(1)

	for q.Len() != 0 {
		log.Println(q.PopFront())
	}
}
```
