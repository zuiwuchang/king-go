package rbtree

//爲 int 實現的 IKey 接口
type IKeyInt int

func (i IKeyInt) Less(k IKey) bool {
	v := k.(IKeyInt)
	return i < v
}
func (i IKeyInt) Equal(k IKey) bool {
	v := k.(IKeyInt)
	return i == v
}
