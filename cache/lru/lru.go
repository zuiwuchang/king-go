package lru

type IKey interface{}
type IValue interface{}

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

	//釋放 緩存並返回 Len()
	//
	//執行後 緩存容量將 <= Cap() * percentage
	Resize(percentage float64) int
}
