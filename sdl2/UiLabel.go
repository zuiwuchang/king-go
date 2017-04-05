package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
	"king-go/algorithm"
	"strings"
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
	pading int32
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
	ui.SetAlpha(255)
	ui.SetScale(1, 1)
	return ui, nil
}
func (u *UiLabel) GetPadding() int32 {
	return u.pading
}
func (u *UiLabel) SetPadding(pading int32) {
	u.pading = pading
}
func (u *UiLabel) destoryTexture() {
	u.restoryTexture()
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
	if u.val != val {
		u.val = val
		u.restoryTexture()
	}
}
func (u *UiLabel) SetColor(color sdl.Color) {
	u.color = color
}

func (u *UiLabel) GetColor() sdl.Color {
	return u.color
}
func (u *UiLabel) restoryTexture() {
	if u.valTexture != nil {
		u.valTexture.Destroy()
		u.valTexture = nil
	}
}
func (u *UiLabel) sizeUTF8(arrs []rune) (int, error) {
	w, _, e := u.font.SizeUTF8(string(arrs))
	if e != nil {
		return 0, e
	}
	return w, nil
}
func (u *UiLabel) renderUTF8(arrs []rune) (*sdl.Surface, error) {
	return u.font.RenderUTF8_Blended(string(arrs), u.color)
}

//將 字符串 繪製到 surface 返回 是否 可以繼續 繪製
//yptr 當前 已被佔用 高度
func (u *UiLabel) blitStr(arrs []rune, surfaceTarget *sdl.Surface, yptr *int32) bool {
	y := *yptr
	//可用 高度
	h := surfaceTarget.H
	if h <= y {
		//高度全被佔用 忽略 繪製 內容
		return false
	}
	h -= y

	//返回 可被繪製 的 最短 長度
	//尋找 最大可顯示文本
	size := len(arrs)
	pos := 0
	width := int(surfaceTarget.W)
	n, e := algorithm.BinarySearch(0, size-1, func(i int) (int, error) {
		w, e := u.sizeUTF8(arrs[:i+1])

		if e != nil {
			return 0, e
		}
		if w >= width {
			return 1, nil
		}
		if i+1 == size {
			pos = w
			return 0, nil
		}

		w2, e := u.sizeUTF8(arrs[:i+2])

		if e != nil {
			return 0, e
		}
		if w2 >= width {
			pos = w
			return 0, nil
		}
		return -1, nil
	})
	if e != nil {
		return false
	}
	n++
	surface, e := u.renderUTF8(arrs[:n])
	if e != nil {
		g_log.Println(e)
		return false
	}

	e = surface.Blit(nil,
		surfaceTarget,
		&sdl.Rect{
			X: 0,
			Y: y,
			W: surface.W,
			H: surface.H,
		},
	)
	if e != nil {
		g_log.Println(e)
		return false
	}
	*yptr = y + surface.H
	if size != n {

		return u.blitStr(arrs[n:], surfaceTarget, yptr)
	}
	return true
}
func (u *UiLabel) initTexture(renderer *sdl.Renderer) {
	if u.valTexture != nil {
		return
	}
	if u.val == "" {
		return
	}

	//創建 surface
	w, h := u.GetSize()
	surfaceTarget, e := sdl.CreateRGBSurface(0,
		int32(w)-u.pading*2,
		int32(h)-u.pading*2,
		32,
		R_MASK,
		G_MASK,
		B_MASK,
		A_MASK,
	)
	if e != nil {
		g_log.Println(e)
		return
	}
	defer surfaceTarget.Free()

	//分隔換行
	str := u.val
	strs := strings.Split(str, "\n")
	y := int32(0)
	for _, str := range strs {
		if !u.blitStr([]rune(str), surfaceTarget, &y) {
			break
		}
	}

	//創建紋理
	texture, e := renderer.CreateTextureFromSurface(surfaceTarget)
	if e != nil {
		g_log.Println(e)
		return
	}
	u.valTexture = texture
}
func (u *UiLabel) Draw(renderer *sdl.Renderer, duration time.Duration) {
	//繪製背景
	x, y := u.GetDrawPos()
	w, h := u.GetDrawSize()

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
	u.initTexture(renderer)
	if u.valTexture != nil {
		renderer.Copy(u.valTexture,
			nil,
			&sdl.Rect{X: int32(x) + u.pading,
				Y: int32(y) + u.pading,
				W: int32(w),
				H: int32(h),
			},
		)
	}
}
