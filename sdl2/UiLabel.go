package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
	"time"
)

//靜態文本框
type UiLabel struct {
	UiBase

	//字體
	font *ttf.Font

	//文本
	val string
	//文本 紋理
	valTexture *sdl.Texture
	//文本顏色
	color sdl.Color

	//邊距
	pading int
}

func NewUiLabel() (*UiLabel, error) {
	return NewUiLabelFont(FONT_DEFAULT_FILE, FONT_DEFAULT_SIZE)
}
func NewUiLabelFont(fFile string, fSize int) (*UiLabel, error) {
	font, e := ttf.OpenFont(fFile, fSize)
	if e != nil {
		return nil, e
	}

	ui := &UiLabel{
		font:   font,
		pading: 3,
		color:  sdl.Color{R: 255, G: 255, B: 255, A: 255},
	}
	return ui, nil
}
func (u *UiLabel) destoryTexture() {
	if u.valTexture != nil {
		u.valTexture.Destroy()
		u.valTexture = nil
	}
}
func (u *UiLabel) Destroy() {
	if u.font != nil {
		u.font.Close()
		u.font = nil
	}

	u.destoryTexture()

	u.UiBase.Destroy()
}
func (u *UiLabel) GetValue() string {
	return u.val
}
func (u *UiLabel) SetValue(val string) {
	u.val = val
}
func (u *UiLabel) SetColor(color sdl.Color) {
	u.color = color
}

func (u *UiLabel) GetColor() sdl.Color {
	return u.color
}
func (u *UiLabel) initTexture() {
	if u.valTexture != nil {
		return
	}

	//
}
func (u *UiLabel) Draw(renderer *sdl.Renderer, duration time.Duration) {
	//繪製背景
	x, y := u.GetDrawPos()
	w, h := u.GetSize()

	if texture := u.GetTexture(); texture != nil {
		renderer.Copy(texture,
			nil,
			&sdl.Rect{X: int32(x),
				Y: int32(y),
				W: int32(w),
				H: int32(h),
			},
		)
	}

	if u.val == "" {
		return
	}
	//繪製文本
	u.initTexture()
	renderer.Copy(u.valTexture,
		nil,
		&sdl.Rect{X: int32(x),
			Y: int32(y),
			W: int32(w),
			H: int32(h),
		},
	)
}
