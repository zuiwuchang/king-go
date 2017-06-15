package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	"king-go/container/rbtree"
	"time"
)

const (
	BUFFER_INIT_COUNT = 10
)

//基礎接口定義
type IObject interface {
	//繪製自己
	Draw(renderer *sdl.Renderer, duration time.Duration)

	//執行動作
	OnAction(duration time.Duration)

	//處理 事件 返回 true 停止事件傳遞
	OnEvent(evt sdl.Event) bool

	//是否可見
	IsVisible() bool
	//設置是否可見
	SetVisible(yes bool)

	//返回 坐標 (相對父節點)
	GetPos() (float64, float64)
	//坐標轉 屏幕坐標
	ToScreenPos(x, y float64) (float64, float64)
	//坐標轉 相對坐標
	ToPos(x, y float64) (float64, float64)

	//設置 坐標
	SetPos(x float64, y float64)
	//返回 大小
	GetSize() (int, int)
	//設置 大小
	SetSize(w int, h int)

	//返回繪製坐標
	GetDrawPos() (float64, float64)
	//返回繪製 大小
	GetDrawSize() (float64, float64)

	//返回 z 坐標
	GetZ() int
	//返回 z 坐標
	SetZ(z int)
	//返回 子節點中的 最大 z 坐標
	GetMaxZ() int
	//返回 子節點中的 最小 z 坐標
	GetMinZ() int

	//銷毀 元素
	Destroy()

	//設置元素 id
	SetId(id string)
	//設置元素 Tag
	SetTag(tag string)
	//返回元素 id
	GetId() string
	//返回元素 Tag
	GetTag() string

	//增加一個 子元素
	Add(obj IObject)

	//刪除一個 子元素
	Remove(obj IObject)
	//刪除一個 指定id的 子元素
	RemoveById(id string)
	//刪除 指定 tag 的元素
	RemoveByTag(tag string)

	//返回 子元素 數量
	LenChilds() int

	//返回 父節點
	GetParent() IObject
	//設置 父節點
	SetParent(parent IObject)

	//綁定一個 動作 多次 bind 的 動作 同時被執行
	BindAction(a IAction)
	//移除一個 動作
	RemoveAction(a IAction)

	//設置紋理
	SetTexture(texture *sdl.Texture)
	//返回當前 紋理
	GetTexture() *sdl.Texture

	//返回 錨點
	GetAnchor() (anchorX, anchorY float64)
	//設置 錨點
	SetAnchor(anchorX, anchorY float64)
	SetAnchorX(anchor float64)
	SetAnchorY(anchor float64)

	//設置 透明度
	SetAlpha(a float64)
	//返回 透明度
	GetAlpha() float64

	//設置 縮放
	SetScale(scaleX, scaleY float64)
	SetScaleX(scaleX float64)
	SetScaleY(scaleY float64)

	//返回 縮放
	GetScale() (scaleX float64, scaleY float64)
	GetScaleX() float64
	GetScaleY() float64
}

type RendererFlip sdl.RendererFlip

const (
	FLIP_NONE       = sdl.FLIP_NONE
	FLIP_HORIZONTAL = sdl.FLIP_HORIZONTAL
	FLIP_VERTICAL   = sdl.FLIP_VERTICAL
)

//基礎 接口 實現
type Node struct {
	//紋理
	Texture *sdl.Texture

	//坐標 大小
	X, Y float64

	//放大縮小 [0,+...)
	ScaleX, ScaleY float64
	Width, Height  int
	//z坐標 值越小 越先繪製 越後響應事件
	Z int

	//錨點 [0,1]
	AnchorX, AnchorY float64

	//旋轉角度
	Angle float64
	//翻轉
	Flip sdl.RendererFlip

	//子元素
	childs    *rbtree.RedBlackTree
	lenChilds int

	id  string
	tag string

	//父節點
	parent IObject

	//是否不可見
	hide bool

	//透明
	Alpha float64

	//動作集合
	actions map[IAction]bool
}

func NewNode(x, y float64, z, w, h int, texture *sdl.Texture) *Node {
	return &Node{
		X:       x,
		Y:       y,
		Z:       z,
		Width:   w,
		Height:  h,
		ScaleX:  1,
		ScaleY:  1,
		Texture: texture,
		Alpha:   255,
	}
}

//返回 子節點數量
func (n *Node) LenChilds() int {
	return n.lenChilds
}

//是否可見
func (n *Node) IsVisible() bool {
	return !n.hide
}

//設置是否可見
func (n *Node) SetVisible(ok bool) {
	n.hide = !ok
}

//返回 坐標
func (n *Node) GetPos() (float64, float64) {
	return n.X, n.Y
}

//坐標轉 相對坐標
func (n *Node) ToPos(x, y float64) (float64, float64) {
	dx, dy := n.GetDrawPos()
	return x - dx, y - dy
}

//坐標轉 屏幕坐標
func (n *Node) ToScreenPos(x, y float64) (float64, float64) {
	for node := n.GetParent(); node != nil; node = node.GetParent() {
		tx, ty := node.GetPos()

		w, h := node.GetSize()
		anchorX, anchorY := node.GetAnchor()
		scaleX, scaleY := node.GetScale()
		tx -= float64(w) * scaleX * anchorX
		ty -= float64(h) * scaleY * anchorY

		x *= scaleX
		y *= scaleY
		x += tx
		y += ty

	}
	return x, y
}

//設置 坐標
func (n *Node) SetPos(x float64, y float64) {
	n.X = x
	n.Y = y
}

//返回 大小
func (n *Node) GetSize() (int, int) {
	return n.Width, n.Height
}

//設置 大小
func (n *Node) SetSize(w int, h int) {
	n.Width = w
	n.Height = h
}

//返回 z 坐標
func (n *Node) GetZ() int {
	return n.Z
}

//返回 z 坐標
func (n *Node) SetZ(z int) {
	if n.Z == z {
		return
	}
	n.Z = z
	p := n.parent
	if p == nil {
		return
	}
	p.Remove(n)
	p.Add(n)
}

//返回 子節點中的 最大 z 坐標
func (n *Node) GetMaxZ() int {
	ele := n.childs.Max()
	if ele == nil {
		return 0
	}
	return int(ele.Key().(rbtree.IKeyInt))
}

//返回 子節點中的 最小 z 坐標
func (n *Node) GetMinZ() int {
	ele := n.childs.Min()
	if ele == nil {
		return 0
	}
	return int(ele.Key().(rbtree.IKeyInt))
}

//繪製自己
func (n *Node) Draw(renderer *sdl.Renderer, duration time.Duration) {
	if !n.IsVisible() {
		//不可見 直接返回
		return
	}

	//繪製自己
	n.draw(renderer, duration)
	//繪製子元素
	if n.childs != nil {
		n.childs.Do(func(ele rbtree.IElement) bool {
			arrs := ele.Value().([]IObject)
			for i := 0; i < len(arrs); i++ {
				arrs[i].Draw(renderer, duration)
			}
			return true
		})
	}
}
func (n *Node) GetDrawScale() (float64, float64) {
	x, y := n.GetScale()

	for node := n.GetParent(); node != nil; node = node.GetParent() {
		tx, ty := node.GetScale()
		if x != 0 {
			x *= tx
		}
		if y != 0 {
			y *= ty
		}
		if x == 0 && y == 0 {
			break
		}
	}

	return x, y
}

//返回繪製 大小
func (n *Node) GetDrawSize() (float64, float64) {
	scaleX, scaleY := n.GetDrawScale()
	return float64(n.Width) * scaleX, float64(n.Height) * scaleY
}

//返回繪製坐標
func (n *Node) GetDrawPos() (float64, float64) {
	w, h := n.GetSize()
	x, y := n.GetPos()

	scaleX, scaleY := n.GetScale()

	x -= float64(w) * scaleX * n.AnchorX
	y -= float64(h) * scaleY * n.AnchorY

	return n.ToScreenPos(x, y)
}

//繪製自己
func (n *Node) draw(renderer *sdl.Renderer, duration time.Duration) {
	alpha := n.GetAlpha()
	if alpha == 0 {
		return
	}
	//繪製自己
	texture := n.Texture
	if texture != nil {
		texture.SetAlphaMod(uint8(alpha))

		x, y := n.GetDrawPos()
		w, h := n.GetDrawSize()
		renderer.CopyEx(texture,
			nil,
			&sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)},
			n.Angle,
			nil,
			n.Flip,
		)
	}
}
func (n *Node) OnAction(duration time.Duration) {
	if n.actions != nil {
		for a, _ := range n.actions {
			a.DoAction(n, duration)
		}
	}

	//執行動作
	if n.childs != nil {
		n.childs.Do(func(ele rbtree.IElement) bool {
			arrs := ele.Value().([]IObject)
			for i := 0; i < len(arrs); i++ {
				arrs[i].OnAction(duration)
			}
			return true
		})
	}
}

//處理 事件 返回 true 停止事件傳遞
func (n *Node) OnEvent(evt sdl.Event) bool {
	size := n.LenChilds()
	if size == 0 {
		return false
	}
	//備份子節點 防止 在 OnEvent 中 改變子節點
	childs := make([]IObject, size, size)
	index := 0
	n.childs.Do(func(ele rbtree.IElement) bool {
		arrs := ele.Value().([]IObject)
		for i := 0; i < len(arrs); i++ {
			childs[index] = arrs[i]
			index++
		}
		return true
	})

	//詢問 子元素
	for i := size - 1; i > -1; i-- {
		if childs[i].IsVisible() &&
			childs[i].OnEvent(evt) {
			return true
		}
	}
	return false
}

//銷毀 元素
func (n *Node) Destroy() {
	if n.childs != nil {
		n.childs.Do(func(ele rbtree.IElement) bool {
			arrs := ele.Value().([]IObject)
			for i := 0; i < len(arrs); i++ {
				arrs[i].Destroy()
			}
			return true
		})
		n.childs.Reset()
	}

	if n.actions != nil {
		for a, _ := range n.actions {
			if a.GetAutoDestroy() {
				a.Destroy()
			}
		}
		n.actions = nil
	}
}

//初始化 子元素 slice
func (n *Node) initChilds() {
	if n.childs == nil {
		n.childs = rbtree.New()
	}
}

//設置元素 id
func (n *Node) SetId(id string) {
	n.id = id
}

//設置元素 Tag
func (n *Node) SetTag(tag string) {
	n.tag = tag
}

//返回元素 id
func (n *Node) GetId() string {
	return n.id
}

//返回元素 Tag
func (n *Node) GetTag() string {
	return n.tag
}

//增加一個 子元素
func (n *Node) Add(obj IObject) {
	//如果存在父節點 從父節點中刪除
	parent := obj.GetParent()
	if parent != nil {
		parent.Remove(obj)
	}

	//添加到當前節點
	n.initChilds()

	t := n.childs
	z := obj.GetZ()
	ele := t.Get(rbtree.IKeyInt(z))
	var arrs []IObject
	if ele == nil {
		arrs = make([]IObject, 1, BUFFER_INIT_COUNT)
		arrs[0] = obj
		t.Insert(rbtree.IKeyInt(z), arrs)
	} else {
		arrs = ele.Value().([]IObject)
		arrs = append(arrs, obj)
		ele.SetValue(arrs)
	}
	obj.SetParent(n)

	n.lenChilds++

}

//刪除一個 子元素
func (n *Node) Remove(obj IObject) {
	n.childs.Do(func(ele rbtree.IElement) bool {
		arrs := ele.Value().([]IObject)
		for i := 0; i < len(arrs); i++ {
			if arrs[i] == obj {
				arrs = append(arrs[:i], arrs[i+1:]...)

				arrs = append(arrs[:i], arrs[i+1:]...)
				if len(arrs) == 0 {
					n.childs.Erase(ele)
				} else {
					ele.SetValue(arrs)
				}

				n.lenChilds--
				return false
			}
		}
		return true
	})
}

//刪除一個 指定id的 子元素
func (n *Node) RemoveById(id string) {
	n.childs.Do(func(ele rbtree.IElement) bool {
		arrs := ele.Value().([]IObject)
		for i := 0; i < len(arrs); i++ {
			if arrs[i].GetId() == id {
				arrs = append(arrs[:i], arrs[i+1:]...)

				if len(arrs) == 0 {
					n.childs.Erase(ele)
				} else {
					ele.SetValue(arrs)
				}

				n.lenChilds--
				return false
			}
		}
		return true
	})
}

//刪除 指定 tag 的元素
func (n *Node) RemoveByTag(tag string) {
	n.childs.Do(func(ele rbtree.IElement) bool {
		arrs := ele.Value().([]IObject)
		for i := 0; i < len(arrs); i++ {
			if arrs[i].GetTag() == tag {
				arrs[i].SetParent(nil)

				arrs = append(arrs[:i], arrs[i+1:]...)
				if len(arrs) == 0 {
					n.childs.Erase(ele)
				} else {
					ele.SetValue(arrs)
				}

				n.lenChilds--
				return false
			}
		}
		return true
	})
}

//返回 父節點
func (n *Node) GetParent() IObject {
	return n.parent
}

//設置 父節點 不要調用 此函數 除非你的確知道在做什麼
//要 添加子節點 調用 父節點的 Add 方法
func (n *Node) SetParent(obj IObject) {
	n.parent = obj
}

//綁定一個 動作 多次 bind 的 動作 同時被執行
func (n *Node) BindAction(a IAction) {
	if n.actions == nil {
		n.actions = make(map[IAction]bool)
	}
	n.actions[a] = true
}

//移除一個 動作
func (n *Node) RemoveAction(a IAction) {
	if n.actions != nil {
		delete(n.actions, a)
	}
	if a.GetAutoDestroy() {
		a.Destroy()
	}
}

//設置紋理
func (n *Node) SetTexture(texture *sdl.Texture) {
	n.Texture = texture
}

//返回當前 紋理
func (n *Node) GetTexture() *sdl.Texture {
	return n.Texture
}

//返回 錨點
func (n *Node) GetAnchor() (anchorX, anchorY float64) {
	return n.AnchorX, n.AnchorY
}

//設置 錨點
func (n *Node) SetAnchor(anchorX, anchorY float64) {
	n.AnchorX = anchorX
	n.AnchorY = anchorY
}
func (n *Node) SetAnchorX(anchor float64) {
	n.AnchorX = anchor
}
func (n *Node) SetAnchorY(anchor float64) {
	n.AnchorY = anchor
}

//設置 透明度
func (n *Node) SetAlpha(a float64) {
	n.Alpha = a
}

//返回 透明度
func (n *Node) GetAlpha() float64 {
	return n.Alpha
}

//設置 縮放
func (n *Node) SetScale(scaleX, scaleY float64) {
	if scaleX < 0 {
		scaleX = 0
	}
	if scaleY < 0 {
		scaleY = 0
	}
	n.ScaleX = scaleX
	n.ScaleY = scaleY
}
func (n *Node) SetScaleX(scale float64) {
	if scale < 0 {
		n.ScaleX = 0
		return
	}
	n.ScaleX = scale
}
func (n *Node) SetScaleY(scale float64) {
	if scale < 0 {
		n.ScaleY = 0
		return
	}
	n.ScaleY = scale
}

//返回 縮放
func (n *Node) GetScale() (float64, float64) {
	return n.ScaleX, n.ScaleY
}
func (n *Node) GetScaleX() float64 {
	return n.ScaleX
}
func (n *Node) GetScaleY() float64 {
	return n.ScaleY
}
