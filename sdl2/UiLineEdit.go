package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
	"log"
	"time"
)

const (
	ui_LineEdit_offsetx = 3
	ui_LineEdit_offsety = 2
)

type uiLineEditVal struct {
	//編輯框 文本數組
	arrs []rune

	//文本長度
	size int

	//文本字符串
	text string

	//光標位置 [0,len(text)]
	chartBegin int
	chartEnd   int

	//水平滾動
	scroll int
}

//返回 滾動
func (u *uiLineEditVal) getScroll(font *ttf.Font, w int) int {
	w -= ui_LineEdit_offsetx * 2
	if u.chartBegin == u.chartEnd {
		chart := u.chartBegin
		for i := chart - 1; i >= 0; i-- {
			str := string(u.arrs[i:chart])
			if width, _, e := font.SizeUTF8(str); e != nil {
				log.Println(e)
				return 0
			} else if width > w {
				u.scroll = i + 1
				str := string(u.arrs[:i])
				pos, _, e := font.SizeUTF8(str)
				if e != nil {
					log.Println(e)
					return 0
				}
				return pos
			}
		}
		u.scroll = 0
		return 0
	}

	return 0
}
func (u *uiLineEditVal) SetRune(arrs []rune) {
	max := len(u.arrs)
	size := len(arrs)
	if size > max {
		size = max
	}
	copy(u.arrs, arrs[:size])

	u.size = size
	u.chartBegin = size
	u.chartEnd = size
	u.text = string(u.arrs)
}
func (u *uiLineEditVal) GetString() string {
	return u.text
}

//單行編輯框
type UiLineEdit struct {
	UiBase

	//當前文本
	val uiLineEditVal
	//字體
	font *ttf.Font
	//文本顏色
	color sdl.Color
	//背景顏色
	bkColor sdl.Color

	//渲染圖像
	valTexture *sdl.Texture

	lastChart time.Time
	cursor    *sdl.Cursor
	cursorSet bool
	cursorOld *sdl.Cursor

	tmpw, tmph int
}

func NewUiLineEdit(max /*允許的最大文本長度*/ int) (*UiLineEdit, error) {
	return NewUiLineEditFont(max, FONT_DEFAULT_FILE, FONT_DEFAULT_SIZE)
}
func NewUiLineEditFont(max /*允許的最大文本長度*/ int, fFile string, fSize int) (*UiLineEdit, error) {
	font, e := ttf.OpenFont(fFile, fSize)
	if e != nil {
		return nil, e
	}

	ui := &UiLineEdit{font: font}
	ui.val.arrs = make([]rune, max, max)
	ui.bkColor = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ui.color = sdl.Color{R: 255, G: 255, B: 255, A: 255}

	ui.cursor = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM)
	return ui, nil
}

func (u *UiLineEdit) SetValue(text string) {
	arrs := []rune(text)
	u.val.SetRune(arrs)

	if u.valTexture != nil {
		u.valTexture.Destroy()
		u.valTexture = nil
	}
}

func (u *UiLineEdit) initValTexture(renderer *sdl.Renderer) *sdl.Texture {
	if u.valTexture != nil {
		return u.valTexture
	}

	font := u.font
	surface, e := font.RenderUTF8_Blended(u.val.GetString(), u.color)
	if e != nil {
		log.Println(e)
		return nil
	}
	{
		texture, e := renderer.CreateTextureFromSurface(surface)
		if e != nil {
			log.Println(e)
			return nil
		}
		u.valTexture = texture
		u.tmpw = surface.W
		u.tmph = surface.H
	}
	return u.valTexture
	w, h := u.GetSize()
	surfaceTarget, e := sdl.CreateRGBSurface(0,
		int32(w)+ui_LineEdit_offsetx*2,
		int32(h)+ui_LineEdit_offsetx*2,
		32,
		R_MASK,
		G_MASK,
		B_MASK,
		A_MASK,
	)
	if e != nil {
		log.Println(e)
		return nil
	}
	pos := u.val.getScroll(font, w)
	log.Println(pos)
	src := sdl.Rect{X: int32(pos)}
	srcW := surface.W - int32(pos)
	dist := sdl.Rect{X: ui_LineEdit_offsetx, Y: -ui_LineEdit_offsetx}
	if srcW < surfaceTarget.W {
		src.W = srcW
	} else {
		src.W = surface.W
	}
	if surface.H < surfaceTarget.H {
		src.H = surface.H
	} else {
		src.H = surface.H
	}

	dist.W = src.W
	dist.H = src.H

	e = surface.Blit(&src, surfaceTarget, &dist)
	if e != nil {
		log.Println(e)
		return nil
	}
	texture, e := renderer.CreateTextureFromSurface(surfaceTarget)
	if e != nil {
		log.Println(e)
		return nil
	}
	u.valTexture = texture
	return texture
}
func (u *UiLineEdit) Draw(renderer *sdl.Renderer, duration time.Duration) {
	//繪製背景
	x, y := u.GetDrawPos()

	renderer.SetDrawColor(u.bkColor.R, u.bkColor.G, u.bkColor.B, u.bkColor.A)
	renderer.FillRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(u.Width), H: int32(u.Height)})

	//繪製文本
	texture := u.initValTexture(renderer)
	if texture != nil {
		x, y := u.GetDrawPos()
		w, h := u.GetSize()
		renderer.Copy(texture,
			nil,
			&sdl.Rect{X: int32(x), Y: int32(y), W: int32(w), H: int32(h)},
		)
	}

	//绘制插入符
	if getDirector().focus == u {
		u.drawChart(renderer)
	}
}
func (u *UiLineEdit) drawChart(renderer *sdl.Renderer) {
	now := time.Now()
	if u.lastChart.IsZero() {
		u.lastChart = now
	} else if now.After(u.lastChart.Add(time.Millisecond * 700)) {
		if now.After(u.lastChart.Add(time.Millisecond * 1400)) {
			u.lastChart = now
		} else {
			return
		}
	}

	renderer.SetDrawColor(255, 128, 128, 0)
	x, y := u.GetDrawPos()
	_, h := u.GetSize()

	x1 := int32(x + ui_LineEdit_offsetx + 20)
	y1 := int32(y + 1)
	renderer.FillRect(&sdl.Rect{
		X: x1,
		Y: y1,
		W: 3,
		H: int32(h - 2),
	})
}
func (u *UiLineEdit) Destroy() {
	u.UiBase.Destroy()
	if u.cursor != nil {
		sdl.FreeCursor(u.cursor)
	}
	if u.font != nil {
		u.font.Close()
	}
}
func (u *UiLineEdit) SetColor(color sdl.Color) {
	u.color = color
}
func (u *UiLineEdit) SetBackgroundColor(color sdl.Color) {
	u.bkColor = color
}
func (u *UiLineEdit) GetColor() sdl.Color {
	return u.color
}
func (u *UiLineEdit) GetBackgroundColor() sdl.Color {
	return u.bkColor
}
func (u *UiLineEdit) onBtnEvt(t *sdl.MouseButtonEvent) bool {
	if t.Type == sdl.MOUSEBUTTONDOWN {
		if u.posInRect(t.X, t.Y) {
			SetFocus(u)

			return true
		}
	}
	return false
}
func (u *UiLineEdit) OnEvent(evt sdl.Event) bool {
	switch t := evt.(type) {
	case *sdl.MouseButtonEvent:
		u.onBtnEvt(t)
	case *sdl.MouseMotionEvent:
		if u.posInRect(t.X, t.Y) {
			if !u.cursorSet {
				u.cursorOld = sdl.GetCursor()
				sdl.SetCursor(u.cursor)
				u.cursorSet = true
			}
			return true
		} else if u.cursorSet {
			u.cursorSet = false
			if u.cursorOld != nil {
				sdl.SetCursor(u.cursorOld)
			}
		}
	}

	return false
}
