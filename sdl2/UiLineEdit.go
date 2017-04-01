package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
	"time"
)

const (
	ui_LineEdit_offsetx = 4
	ui_LineEdit_offsety = 2
)

//單行編輯框
type UiLineEdit struct {
	UiBase

	//當前文本
	val valLineEdit

	//背景顏色
	bkColor sdl.Color

	//光標
	cursor    *sdl.Cursor
	cursorSet bool
	cursorOld *sdl.Cursor
}

func NewUiLineEdit() (*UiLineEdit, error) {
	return NewUiLineEditFont(FONT_DEFAULT_FILE, FONT_DEFAULT_SIZE)
}
func NewUiLineEditFont(fFile string, fSize int) (*UiLineEdit, error) {
	font, e := ttf.OpenFont(fFile, fSize)
	if e != nil {
		return nil, e
	}

	ui := &UiLineEdit{val: valLineEdit{
		font:  font,
		color: sdl.Color{R: 255, G: 255, B: 255, A: 255},
	}}
	ui.val.resetStr()
	ui.val.SetPwdChar("*")
	ui.val.SetChartRGB(128, 128, 128)
	ui.bkColor = sdl.Color{R: 0, G: 0, B: 0, A: 255}
	ui.cursor = sdl.CreateSystemCursor(sdl.SYSTEM_CURSOR_IBEAM)
	return ui, nil
}
func (u *UiLineEdit) GetValue() string {
	return u.val.GetString()
}
func (u *UiLineEdit) SetValue(text string) {
	w, _ := u.GetSize()
	u.val.SetString(text, int(w-ui_LineEdit_offsetx*2))
}

func (u *UiLineEdit) Draw(renderer *sdl.Renderer, duration time.Duration) {
	//繪製背景
	x, y := u.GetDrawPos()
	w, h := u.GetSize()

	renderer.SetDrawColor(u.bkColor.R, u.bkColor.G, u.bkColor.B, u.bkColor.A)
	renderer.FillRect(&sdl.Rect{X: int32(x), Y: int32(y), W: int32(u.Width), H: int32(u.Height)})

	//繪製文本
	u.val.DrawText(renderer,
		int32(x)+ui_LineEdit_offsetx,
		int32(y)+ui_LineEdit_offsety,
		int32(w)-ui_LineEdit_offsetx*2,
		int32(h)-ui_LineEdit_offsety*2,
	)

	//绘制插入符
	if getDirector().focus == u {
		u.val.DrawChart(renderer, int32(x+ui_LineEdit_offsetx), int32(y+ui_LineEdit_offsety), 3, int32(h-ui_LineEdit_offsety*2))
	}
}

func (u *UiLineEdit) Destroy() {

	if u.cursor != nil {
		sdl.FreeCursor(u.cursor)
	}
	u.val.Destroy()

	u.UiBase.Destroy()
}
func (u *UiLineEdit) SetChartRGB(r, g, b uint8) {
	u.val.SetChartRGB(r, g, b)
}
func (u *UiLineEdit) GetChartRGB() (r uint8, g uint8, b uint8) {
	return u.val.GetChartRGB()
}
func (u *UiLineEdit) SetColor(color sdl.Color) {
	u.val.color = color
}
func (u *UiLineEdit) SetBackgroundColor(color sdl.Color) {
	u.bkColor = color
}
func (u *UiLineEdit) GetColor() sdl.Color {
	return u.val.color
}
func (u *UiLineEdit) GetBackgroundColor() sdl.Color {
	return u.bkColor
}
func (u *UiLineEdit) onBtnEvt(t *sdl.MouseButtonEvent) bool {
	if t.Type == sdl.MOUSEBUTTONDOWN {
		if u.posInRect(t.X, t.Y) {
			SetFocus(u)

			if t.Button == sdl.BUTTON_LEFT {
				x, _ := u.GetDrawPos()
				//開始 選擇
				u.val.SelectStart(t.X - int32(x))
			}
			return true
		}
	} else if t.Type == sdl.MOUSEBUTTONUP {
		if t.Button == sdl.BUTTON_LEFT {
			x, _ := u.GetDrawPos()
			//結束選擇
			u.val.SelectStop(t.X - int32(x))
		}
	}
	return false
}
func (u *UiLineEdit) OnEvent(evt sdl.Event) bool {
	switch t := evt.(type) {
	case *sdl.KeyDownEvent:
		if getDirector().focus == u {
			if t.Keysym.Sym == sdl.K_LEFT {
				u.val.SelectLeft()
			} else if t.Keysym.Sym == sdl.K_RIGHT {
				u.val.SelectRight()
			} else if t.Keysym.Sym == sdl.K_BACKSPACE {
				u.val.Backspace()
			} else if t.Keysym.Sym == sdl.K_a {
				state := sdl.GetKeyboardState()
				if state[sdl.SCANCODE_LCTRL] == 1 ||
					state[sdl.SCANCODE_RCTRL] == 1 {
					//select all
					u.SelectAll()
				}
			} else if t.Keysym.Sym == sdl.K_c &&
				!u.val.IsPwd() {
				state := sdl.GetKeyboardState()
				if state[sdl.SCANCODE_LCTRL] == 1 ||
					state[sdl.SCANCODE_RCTRL] == 1 {
					//copy
					str := u.val.GetSelectStr()
					if str != "" {
						e := sdl.SetClipboardText(str)
						if e != nil {
							g_log.Println(e)
						}
					}
				}
			} else if t.Keysym.Sym == sdl.K_x &&
				!u.val.IsPwd() {
				state := sdl.GetKeyboardState()
				if state[sdl.SCANCODE_LCTRL] == 1 ||
					state[sdl.SCANCODE_RCTRL] == 1 {
					//cut
					str := u.val.GetSelectStr()
					if str != "" {
						e := sdl.SetClipboardText(str)
						if e == nil {
							w, _ := u.GetSize()
							u.val.ReplaceStr("", int(w))
						} else {
							g_log.Println(e)
						}
					}
				}
			} else if t.Keysym.Sym == sdl.K_v {
				state := sdl.GetKeyboardState()
				if state[sdl.SCANCODE_LCTRL] == 1 ||
					state[sdl.SCANCODE_RCTRL] == 1 {
					//paste
					str, e := sdl.GetClipboardText()
					if e == nil {
						if str != "" {
							u.ReplaceStr(str)
						}
					} else {
						g_log.Println(e)
					}
				}
			}

			callback := u.GetEventCallback(UI_EVT_KEY_DOWM)
			if callback != nil {
				callback(u, t)
			}
			return true
		}
	case *sdl.TextInputEvent:
		if getDirector().focus == u {
			size := len(t.Text)
			b := make([]byte, 0, size)
			for i := 0; i < size; i++ {
				if t.Text[i] == 0 {
					break
				} else {
					b = append(b, t.Text[i])
				}
			}
			u.ReplaceStr(string(b))
			return true
		}
	case *sdl.MouseButtonEvent:
		u.onBtnEvt(t)
	case *sdl.MouseMotionEvent:
		if u.val.IsSelect() {
			x, _ := u.GetDrawPos()
			u.val.SelectIng(t.X - int32(x))
		}

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

//將 選中項 以 字符串 替換
func (u *UiLineEdit) ReplaceStr(str string) error {
	w, _ := u.GetSize()
	w -= ui_LineEdit_offsetx * 2

	arrs := []rune(str)
	e := u.val.ReplaceRune(arrs, int(w))
	if e != nil {
		g_log.Println(e)
	}
	return e
}

//將 選中項 以 []rune 替換
func (u *UiLineEdit) ReplaceRune(arrs []rune) error {
	w, _ := u.GetSize()
	w -= ui_LineEdit_offsetx * 2

	e := u.val.ReplaceRune(arrs, int(w))
	if e != nil {
		g_log.Println(e)
	}
	return e
}

//設置 允許輸入的最大字符數
func (u *UiLineEdit) GetMax() int {
	return u.val.GetMax()
}

//返回 允許輸入的最大字符數
func (u *UiLineEdit) SetMax(max int) {
	u.val.SetMax(max)
}

//光標 選擇
func (u *UiLineEdit) Select(begin, end int) {
	u.val.Select(begin, end)
}

//光標 全選
func (u *UiLineEdit) SelectAll() {
	str := u.val.GetString()
	if str == "" {
		return
	}

	u.val.Select(0, len([]rune(str)))
}

//返回 是否 是密碼框
func (u *UiLineEdit) IsPwd() bool {
	return u.val.isPwd
}

//設置 是否 是密碼框
func (u *UiLineEdit) SetPwd(yes bool) {
	u.val.SetPwd(yes)
}

//設置 密碼框 顯示 文本
func (u *UiLineEdit) SetPwdChar(c string) {
	u.val.SetPwdChar(c)
}

//返回 密碼框 顯示 文本
func (u *UiLineEdit) GetPwdChar() string {
	return u.val.GetPwdChar()
}
