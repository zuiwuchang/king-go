# lru2
lru2 用法基本同 lru 請參考 [king-go/cache/lru](https://github.com/zuiwuchang/king-go/tree/master/cache/lru) 說明

lru2 只是 爲lru 的 每個 緩存的項目 增加了個 超時時間 一旦 超時 緩存將被自動 移除
```Go
package lru2

import (
	"time"
)

type IKey interface{}
type IValue interface{}

//創建 一個 內存 緩存
func NewLRU(maxElementSize int) ILRU {
	return NewLRU2(time.Second*60*60, maxElementSize)
}

//創建 一個 內存 緩存
func NewLRU2(expired time.Duration, maxElementSize int) ILRU {
	return newLRUImpl(expired, maxElementSize)
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
	//返回 是否存在 緩存 不會 移動節點 更新超時時間
	Ok(key IKey) bool
	//返回 緩存值 不存在 返回 nil
	Get(key IKey) IValue
	//創建 一個 緩存
	Set(key IKey, val IValue)
	//創建 一個 緩存 同時指定 超時 時間
	Set2(key IKey, val IValue, expired time.Duration)
}

```

# 注意
請儘量使用 lru 而非 lru2 

lru2 爲每個 項目 創建了一個 超時 計時器 而且 一旦 get 需要 更新 計時器 孤不得不講 讀寫鎖 換成 互斥鎖 無論內存利用率還是效率 lru2 都不如 lru
