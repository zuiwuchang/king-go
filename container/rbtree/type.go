//紅黑樹
package rbtree

//key 定義
type IKey interface {
	// -1(<) 0(==) 1(>)
	Compare(k IKey) int
}

//value 定義
type IValue interface{}

//節點定義
type IElement interface {
	//返回 key
	Key() IKey
	//返回 value
	Value() IValue
}
