package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

//基礎接口定義
type Object interface {
	//繪製自己
	Draw(renderer *sdl.Renderer, duration time.Duration)

	//處理 事件 返回 true 停止事件傳遞
	OnEvent(evt sdl.Event) bool

	//返回 坐標
	GetPos() (int, int)
	//返回 大小
	GetSize() (int, int)

	//銷毀 元素
	Destroy()
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
	X, Y, Width, Height int
	//z坐標 值越小 越先繪製 越後響應事件
	Z int

	//錨點 [0,100]
	AnchorX, AnchorY int

	//旋轉角度
	Angle float64
	//翻轉
	Flip sdl.RendererFlip

	//子元素
	childs []*Node
	//帶id的子元素
	keyChilds map[string]*Node
}

//返回 坐標
func (n *Node) GetPos() (int, int) {
	return n.X, n.Y
}

//返回 大小
func (n *Node) GetSize() (int, int) {
	return n.Width, n.Height
}

//繪製自己
func (n *Node) Draw(renderer *sdl.Renderer, duration time.Duration) {
	texture := n.Texture
	if texture != nil {
		width := int32(n.Width)
		height := int32(n.Height)
		x := int32(n.X - n.Width*n.AnchorX/100)
		y := int32(n.Height*n.AnchorY/100 - n.Y)
		renderer.CopyEx(texture,
			nil,
			&sdl.Rect{X: x, Y: y, W: width, H: height},
			n.Angle,
			nil,
			n.Flip,
		)
	}
}

//處理 事件 返回 true 停止事件傳遞
func (n *Node) OnEvent(evt sdl.Event) bool {
	return false
}

//銷毀 元素
func (n *Node) Destroy() {

}
