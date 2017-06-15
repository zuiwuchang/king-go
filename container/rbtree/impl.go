package rbtree

//定義一個 樹 結構
type _Element struct {
	K IKey
	V IValue

	//顏色定義
	Red bool

	//父節點
	P *_Element

	//子節點
	L, R *_Element
}

//定義葉節點
var _ElementNil *_Element = &_Element{Red: false}

//實現 IElement 接口
func (ele *_Element) Key() IKey {
	return ele.K
}

//實現 IElement 接口
func (ele *_Element) Value() IValue {
	return ele.V
}
func (ele *_Element) SetValue(v interface{}) {
	ele.V = v
}

//左旋轉
func leftRotate(root, x *_Element) *_Element {
	//緩存 y
	y := x.R

	//移動 β
	x.R = y.L
	if y.L != _ElementNil {
		y.L.P = x
	}

	//y 佔據 原 x位置
	y.P = x.P
	if x.P == _ElementNil { //更新root
		root = y
	} else if x == x.P.L {
		x.P.L = y
	} else {
		x.P.R = y
	}

	//x 成爲 y 孩子
	y.L = x
	x.P = y
	return root
}

//右旋轉
func rightRotate(root, y *_Element) *_Element {
	//緩存 x
	x := y.L

	//移動 β
	y.L = x.R
	if x.R != _ElementNil {
		x.R.P = y
	}

	//x 佔據 原 y位置
	x.P = y.P
	if y.P == _ElementNil { //更新root
		root = x
	} else if y == y.P.L {
		y.P.L = x
	} else {
		y.P.R = x
	}

	//y 成爲 x 孩子
	x.R = y
	y.P = x
	return root
}

//插入 並返回新的 root 節點 是否新增節點
func insert(root *_Element, k IKey, v IValue) (*_Element, bool) {
	//創建 新 節點
	z := &_Element{K: k, V: v, Red: true, L: _ElementNil, R: _ElementNil}

	//記錄 z 的 父節點
	y := _ElementNil
	//當前 位置
	x := root
	for x != _ElementNil {
		compare := x.K.Compare(k)
		if compare == 0 {
			//已經存在key 直接替換
			x.V = v
			return root, false
		}

		y = x
		if compare > 0 {
			//比當前小 比對左子樹
			x = x.L
		} else {
			//比當前大 比對右子樹
			x = x.R
		}
	}

	//設置 父節點
	z.P = y
	if y == _ElementNil { //設置新的 root
		root = z
	} else if k.Compare(y.K) < 0 {
		y.L = z
	} else {
		y.R = z
	}
	return insertFixup(root, z), true
}

//修復插入引起的 紅黑樹性質變化
func insertFixup(root, z *_Element) *_Element {
	for z.P.Red {
		if z.P == z.P.P.L {
			y := z.P.P.R
			//case 1
			if y.Red {
				z.P.Red = false
				y.Red = false
				z.P.P.Red = true
				z = z.P.P
				continue
			}

			//case 2
			if z == z.P.R {

				z = z.P
				root = leftRotate(root, z)
			}
			//case 3
			z.P.Red = false
			z.P.P.Red = true
			root = rightRotate(root, z.P.P)
		} else {
			y := z.P.P.L
			//case 1
			if y.Red {
				z.P.Red = false
				y.Red = false
				z.P.P.Red = true
				z = z.P.P

				continue
			}
			//case 2
			if z == z.P.L {
				z = z.P
				root = rightRotate(root, z)
			}

			//case 3
			z.P.Red = false
			z.P.P.Red = true
			root = leftRotate(root, z.P.P)
		}

	}
	root.Red = false
	return root
}

//返回 最小 節點
func min(x *_Element) *_Element {
	for x != _ElementNil && x.L != _ElementNil {
		x = x.L
	}
	return x
}

//以子樹 v 替換 子樹 u 並返回新的 root
func transplant(root, u, v *_Element) *_Element {
	if u.P == _ElementNil { //u爲root節點
		root = v
	} else if u == u.P.L { //u是 左孩子
		u.P.L = v
	} else /*if u==u.P.R*/ { //u是 右孩子
		u.P.R = v
	}

	//更新 v的 父節點
	v.P = u.P

	return root
}

//刪除指定節點 並返回新的 root
func erase(root, z *_Element) *_Element {
	if z == _ElementNil {
		return root
	}

	var x *_Element
	y := z
	red := y.Red
	if z.L == _ElementNil {
		x = z.R
		root = transplant(root, z, z.R)
	} else if z.R == _ElementNil {
		x = z.L
		root = transplant(root, z, z.L)
	} else {
		y = min(z.R)
		red = y.Red
		x = y.R
		if y.P == z {
			x.P = y
		} else {
			root = transplant(root, y, y.R)
			y.R = z.R
			y.R.P = y
		}
		root = transplant(root, z, y)
		y.L = z.L
		y.L.P = y
		y.Red = z.Red
	}

	if !red {
		root = eraseFixup(root, x)
	}

	return root
}
func eraseFixup(root, x *_Element) *_Element {
	for x != root && !x.Red {
		if x == x.P.L {
			w := x.P.R
			if w.Red { //case 1
				w.Red = false
				x.P.Red = true
				root = leftRotate(root, x.P)

				w = x.P.R
			}

			if !w.L.Red && !w.R.Red { //case2
				w.Red = true
				x = x.P
			} else {
				if !w.R.Red {
					w.L.Red = false
					w.Red = true
					root = rightRotate(root, w)
					w = x.P.R
				}

				w.Red = x.P.Red
				x.P.Red = false
				w.R.Red = false
				root = leftRotate(root, x.P)
				x = root
			}
		} else {
			w := x.P.L

			if w.Red { //case 1
				w.Red = false
				x.P.Red = true
				root = rightRotate(root, x.P)
				w = x.P.L

			}
			if !w.R.Red && !w.L.Red { //case2
				w.Red = true
				x = x.P
			} else {
				if !w.L.Red {
					w.R.Red = false
					w.Red = true
					root = leftRotate(root, w)
					w = x.P.L
				}

				w.Red = x.P.Red
				x.P.Red = false
				w.L.Red = false
				root = rightRotate(root, x.P)
				x = root
			}
		}
	}
	x.Red = false
	return root
}

//查找節點
func search(x *_Element, k IKey) *_Element {
	for x != _ElementNil {
		compare := x.K.Compare(k)
		if compare == 0 {
			return x
		}

		if compare > 0 {
			x = x.L
		} else {
			x = x.R
		}
	}
	return x
}
