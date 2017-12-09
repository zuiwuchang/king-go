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
	//返回 是否存在 緩存 不會更新超時時間
	Ok(key IKey) bool
	//返回 緩存值 不存在 返回 nil
	Get(key IKey) IValue
	//創建 一個 緩存
	Set(key IKey, val IValue)
	//創建 一個 緩存 同時指定 超時 時間
	Set2(key IKey, val IValue, expired time.Duration)
}
