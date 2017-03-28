package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	"sort"
	"time"
)

const (
	BUFFER_INIT_COUNT = 10
)

//基礎接口定義
type Object interface {
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
	//設置 坐標
	SetPos(x float64, y float64)
	//返回 大小
	GetSize() (int, int)
	//設置 大小
	SetSize(w int, h int)

	//返回 z 坐標
	GetZ() int
	//返回 z 坐標
	SetZ(z int)
	//返回 子節點中的 最大 z 坐標 必須要sort後才會返回 正確值
	GetMaxZ() int
	//返回 子節點中的 最小 z 坐標 必須要sort後才會返回 正確值
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
	Add(obj Object)
	//將子元素 按z坐標排序
	SortZ()
	//增加一個 子元素 並按z坐標排序
	AddSortZ(obj Object)

	//刪除一個 子元素
	Remove(obj Object)
	//刪除一個 指定id的 子元素
	RemoveById(id string)
	//刪除 指定 tag 的元素
	RemoveByTag(tag string)

	//返回 父節點
	GetParent() Object
	//設置 父節點
	SetParent(parent Object)

	//綁定一個 動作 多次 bind 的 動作 同時被執行
	BindAction(a Action)
	//移除一個 動作
	RemoveAction(a Action)

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
	X, Y          float64
	Width, Height int
	//z坐標 值越小 越先繪製 越後響應事件
	Z int

	//錨點 [0,100]
	AnchorX, AnchorY float64

	//旋轉角度
	Angle float64
	//翻轉
	Flip sdl.RendererFlip

	//子元素
	childs []Object

	id  string
	tag string

	//父節點
	parent Object

	//是否不可見
	hide bool

	//動作集合
	actions map[Action]bool
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

//坐標轉 屏幕坐標
func (n *Node) ToScreenPos(x, y float64) (float64, float64) {
	for node := n.GetParent(); node != nil; node = node.GetParent() {
		tx, ty := node.GetPos()

		w, h := node.GetSize()
		anchorX, anchorY := node.GetAnchor()
		tx = tx - float64(w)*anchorX
		ty = ty - float64(h)*anchorY

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
	n.Z = z
}

//返回 子節點中的 最大 z 坐標 必須要sort後才會返回 正確值
func (n *Node) GetMaxZ() int {
	size := len(n.childs)
	if size == 0 {
		return 0
	}
	return n.childs[size-1].GetZ()
}

//返回 子節點中的 最小 z 坐標 必須要sort後才會返回 正確值
func (n *Node) GetMinZ() int {
	size := len(n.childs)
	if size == 0 {
		return 0
	}
	return n.childs[0].GetZ()
}

//繪製自己
func (n *Node) Draw(renderer *sdl.Renderer, duration time.Duration) {
	//執行動作
	n.OnAction(duration)
	for i := 0; i < len(n.childs); i++ {
		n.childs[i].OnAction(duration)
	}

	if !n.IsVisible() {
		//不可見 直接返回
		return
	}

	//繪製自己
	n.draw(renderer, duration)
	//繪製子元素
	for i := 0; i < len(n.childs); i++ {
		n.childs[i].Draw(renderer, duration)
	}
}

//返回繪製坐標
func (n *Node) GetDrawPos() (float64, float64) {
	x := n.X - float64(n.Width)*n.AnchorX
	y := n.Y - float64(n.Height)*n.AnchorY
	return n.ToScreenPos(x, y)
}

//繪製自己
func (n *Node) draw(renderer *sdl.Renderer, duration time.Duration) {
	//繪製自己
	texture := n.Texture
	if texture != nil {
		x, y := n.GetDrawPos()
		renderer.CopyEx(texture,
			nil,
			&sdl.Rect{X: int32(x), Y: int32(y), W: int32(n.Width), H: int32(n.Height)},
			n.Angle,
			nil,
			n.Flip,
		)
	}
}
func (n *Node) OnAction(duration time.Duration) {
	if n.actions == nil {
		return
	}

	for a, _ := range n.actions {
		a.DoAction(n, duration)
	}
}

//處理 事件 返回 true 停止事件傳遞
func (n *Node) OnEvent(evt sdl.Event) bool {
	//詢問 子元素
	for i := len(n.childs) - 1; i > -1; i-- {
		if n.childs[i].IsVisible() && n.childs[i].OnEvent(evt) {
			return true
		}
	}
	return false
}

//銷毀 元素
func (n *Node) Destroy() {
	for i := 0; i < len(n.childs); i++ {
		n.childs[i].Destroy()
	}
	n.childs = nil
	for a, _ := range n.actions {
		if a.GetAutoDestory() {
			a.Destory()
		}
	}
	n.actions = nil
}

//初始化 子元素 slice
func (n *Node) initChilds() {
	if n.childs == nil {
		n.childs = make([]Object, 0, BUFFER_INIT_COUNT)
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
func (n *Node) Add(obj Object) {
	//如果存在父節點 從父節點中刪除
	parent := obj.GetParent()
	if parent != nil {
		parent.Remove(obj)
	}

	//添加到當前節點
	n.initChilds()
	n.childs = append(n.childs, obj)
	obj.SetParent(n)
}

type SortChilds []Object

func (s SortChilds) Len() int {
	return len(s)
}
func (s SortChilds) Less(i, j int) bool {
	return s[i].GetZ() < s[j].GetZ()
}
func (s SortChilds) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

//將子元素 按z坐標排序
func (n *Node) SortZ() {
	sort.Sort(SortChilds(n.childs))
}

//增加一個 子元素 並按z坐標排序
func (n *Node) AddSortZ(obj Object) {
	n.Add(obj)
	n.SortZ()
}

//刪除一個 子元素
func (n *Node) Remove(obj Object) {
	for i := 0; i < len(n.childs); i++ {
		if n.childs[i] == obj {
			obj.SetParent(nil)

			n.childs = append(n.childs[:i], n.childs[i+1:]...)
			break
		}
	}
}

//刪除一個 指定id的 子元素
func (n *Node) RemoveById(id string) {
	for i := 0; i < len(n.childs); i++ {
		if n.childs[i].GetId() == id {
			n.childs[i].SetParent(nil)

			n.childs = append(n.childs[:i], n.childs[i+1:]...)
			break
		}
	}
}

//刪除 指定 tag 的元素
func (n *Node) RemoveByTag(tag string) {
	for i := 0; i < len(n.childs); i++ {
		if n.childs[i].GetTag() == tag {
			n.childs[i].SetParent(nil)

			n.childs = append(n.childs[:i], n.childs[i+1:]...)
			break
		}
	}
}

//返回 父節點
func (n *Node) GetParent() Object {
	return n.parent
}

//設置 父節點
func (n *Node) SetParent(obj Object) {
	n.parent = obj
}

//綁定一個 動作 多次 bind 的 動作 同時被執行
func (n *Node) BindAction(a Action) {
	if n.actions == nil {
		n.actions = make(map[Action]bool)
	}
	n.actions[a] = true
}

//移除一個 動作
func (n *Node) RemoveAction(a Action) {
	if n.actions != nil {
		delete(n.actions, a)
	}
	if a.GetAutoDestory() {
		a.Destory()
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
