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
		if x.K.Equal(k) {
			//已經存在key 直接替換
			x.V = v
			return root, false
		}

		y = x
		if k.Less(x.K) {
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
	} else if k.Less(y.K) {
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
