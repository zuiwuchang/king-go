package rbtree

//紅黑樹 容器
type RedBlackTree struct {
	//緩存 節點數量
	length int
	//緩存 根節點
	root *_Element
}

//創建一個 紅黑樹
func New() *RedBlackTree {
	return &RedBlackTree{root: _ElementNil}
}

//返回 節點 數量
func (r *RedBlackTree) Len() int {
	return r.length
}

//插入 節點
//返回 true 新增節點 false 替換掉已有節點
func (r *RedBlackTree) Insert(k IKey, v IValue) bool {
	var add bool
	r.root, add = insert(r.root, k, v)
	if add {
		r.length++
	}
	return add
}

//正向遍歷 所有節點 返回 false 停止 遍歷
func (r *RedBlackTree) Do(callback func(k IKey, v IValue) bool) {
	r.do(r.root, callback)
}
func (r *RedBlackTree) do(x *_Element, callback func(k IKey, v IValue) bool) bool {
	if x != _ElementNil {
		return r.do(x.L, callback) &&
			callback(x.K, x.V) &&
			r.do(x.R, callback)
	}
	return true
}

//逆向遍歷 所有節點 返回 false 停止 遍歷
func (r *RedBlackTree) DoReverse(callback func(k IKey, v IValue) bool) {
	r.doReverse(r.root, callback)
}
func (r *RedBlackTree) doReverse(x *_Element, callback func(k IKey, v IValue) bool) bool {
	if x != _ElementNil {
		return r.do(x.R, callback) &&
			callback(x.K, x.V) &&
			r.do(x.L, callback)
	}
	return true
}
