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

//重置 紅黑樹
func (r *RedBlackTree) Reset() {
	r.root = _ElementNil
	r.length = 0
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

//返回最小 節點
func (r *RedBlackTree) Min() IElement {
	x := r.root
	for x != _ElementNil && x.L != _ElementNil {
		x = x.L
	}

	if x == _ElementNil {
		return nil
	}
	return x
}

//返回最大 節點
func (r *RedBlackTree) Max() IElement {
	x := r.root
	for x != _ElementNil && x.R != _ElementNil {
		x = x.R
	}

	if x == _ElementNil {
		return nil
	}
	return x
}

//正向遍歷 所有節點 返回 false 停止 遍歷
func (r *RedBlackTree) Do(callback func(ele IElement) bool) {
	r.do(r.root, callback)
}

//正向遍歷 指定節點樹 返回 false 停止 遍歷
func (r *RedBlackTree) DoTree(root IElement, callback func(ele IElement) bool) {
	if root == nil {
		return
	}
	r.do(root.(*_Element), callback)
}

func (r *RedBlackTree) do(x *_Element, callback func(ele IElement) bool) bool {
	if x != _ElementNil {
		return r.do(x.L, callback) &&
			callback(x) &&
			r.do(x.R, callback)
	}
	return true
}

//逆向遍歷 所有節點 返回 false 停止 遍歷
func (r *RedBlackTree) DoReverse(callback func(ele IElement) bool) {
	r.doReverse(r.root, callback)
}

//逆向遍歷 指定節點樹 返回 false 停止 遍歷
func (r *RedBlackTree) DoTreeReverse(root IElement, callback func(ele IElement) bool) {
	r.doReverse(root.(*_Element), callback)
}
func (r *RedBlackTree) doReverse(x *_Element, callback func(ele IElement) bool) bool {
	if x != _ElementNil {
		return r.doReverse(x.R, callback) &&
			callback(x) &&
			r.doReverse(x.L, callback)
	}
	return true
}

//返回 指定節點/nil
func (r *RedBlackTree) Get(k IKey) IElement {
	x := search(r.root, k)
	if x == _ElementNil {
		return nil
	}
	return x
}

//返回 指定節點的value/nil
func (r *RedBlackTree) GetValue(k IKey) IValue {
	x := search(r.root, k)
	if x == _ElementNil {
		return nil
	}
	return x.V
}

//刪除 指定 節點
func (r *RedBlackTree) Erase(k IKey) {
	if k == nil {
		return
	}

	x := search(r.root, k)
	if x == _ElementNil {
		return
	}
	r.length--
	r.root = erase(r.root, x)
}
