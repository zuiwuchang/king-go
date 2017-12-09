# lru
lru 是使用 lru 算法 實現的 一個線程安全的 內存緩存
對於緩存的 項目 key 和 value 都是 一個 interface{}
```Go
type IKey interface{}
type IValue interface{}
```

# 如何使用
調用 NewLRU 創建一個 ILRU 接口

使用 ILRU 接口提供 的 Get/Set 方法 獲取/設置 緩存
```Go
//創建 一個 內存 緩存
func NewLRU(maxElementSize int) ILRU {
	return newLRUImpl(maxElementSize)
}

//緩存接口定義
type ILRU interface {
	//返回 當前 緩存 量
	Len() int
	//返回 緩存 最高容量
	Cap() int

	//刪除 所有 緩存
	Clear()
	//刪除 指定緩存
	Delete(key IKey)
	//返回 是否存在 緩存 不會移動緩存
	Ok(key IKey) bool
	//返回 緩存值 不存在 返回 nil
	Get(key IKey) IValue
	//創建 一個 緩存
	Set(key IKey, val IValue)
}
```

# 注意
lru 緩存 不會自動過期 只有當 緩存 達到最大容量時 繼續增加緩存 才會 刪除最少使用的 緩存

如果你想讓 某個 緩存 過期 請使用 ILRU.Delete 方法

如果你期望 緩存的某個 key 能夠 自動 超時 過期 可以考慮使用 [king-go/cache/lru2](https://github.com/zuiwuchang/king-go/tree/master/cache/lru2)

# Example
```Go
package main

import (
	"fmt"
	"github.com/zuiwuchang/king-go/cache/lru"
)

func main() {
	//創建 緩存 容量爲 3
	cache := lru.NewLRU(3)

	//添加 緩存
	for i := 0; i < 3; i++ {
		cache.Set(i, 100+i)
	}

	//返回緩存 100
	fmt.Println(cache.Get(0))
	//因爲 容量 已滿 刪除 最少使用/最後使用的
	//也就是元素 1 101
	//因爲 Get(0) 100 比 101 更多的被使用到
	cache.Set(3, 104)

	//100 102 103
	for i := 0; i < 4; i++ {
		v := cache.Get(i)
		if v != nil {
			fmt.Print(v, ",")
		}
	}
	fmt.Println()
}
```
