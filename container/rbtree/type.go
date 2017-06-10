//紅黑樹
package rbtree

//key 定義
type IKey interface {
	// <
	Less(k IKey) bool
	// ==
	Equal(k IKey) bool
}

//value 定義
type IValue interface{}

//節點定義
type IElement interface {
	Key() IKey
	Value() IValue
}
