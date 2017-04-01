package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
)

//對話框
//對於模式對話框 需要自己保證 其 z 值爲最大
type UiDialog struct {
	UiBase

	//是否可移動
	canMove bool
	//是否爲 模式對話框
	isModel bool

	isLDown bool
}

func (u *UiDialog) SetMove(yes bool) {
	u.canMove = yes
}
func (u *UiDialog) IsMove() bool {
	return u.canMove
}
func (u *UiDialog) SetModel(yes bool) {
	u.isModel = yes
}
func (u *UiDialog) IsModel() bool {
	return u.isModel
}
func NewUiDialog(move bool, model bool) *UiDialog {
	ui := &UiDialog{
		canMove: move,
		isModel: model,
	}
	ui.SetAlpha(255)
	ui.SetVisible(false)
	return ui
}
func (u *UiDialog) onButton(t *sdl.MouseButtonEvent, ok bool) bool {
	switch t.Button {
	case sdl.BUTTON_LEFT:
		if t.Type == sdl.MOUSEBUTTONDOWN {
			if u.posInRect(t.X, t.Y) {
				if u.canMove {
					u.isLDown = true
				}
				SetFocus(u)
				return true
			}
		} else {
			if u.isLDown {
				u.isLDown = false
			}
			return false
		}
	case sdl.BUTTON_RIGHT:
		if t.Type == sdl.MOUSEBUTTONDOWN {
			if u.posInRect(t.X, t.Y) {
				SetFocus(u)
				return true
			}
		} else {
			return false
		}
	case sdl.BUTTON_MIDDLE:
		if t.Type == sdl.MOUSEBUTTONDOWN {
			if u.posInRect(t.X, t.Y) {
				SetFocus(u)
				return true
			}
		} else {
			return false
		}
	}
	return ok
}
func (u *UiDialog) doMove(x, y int32) {
	nx, ny := u.GetPos()
	nx += float64(x)
	ny += float64(y)
	u.SetPos(nx, ny)
}
func (u *UiDialog) OnEvent(evt sdl.Event) bool {
	//詢問 子元素
	for i := len(u.childs) - 1; i > -1; i-- {
		if u.childs[i].IsVisible() && u.childs[i].OnEvent(evt) {
			return true
		}
	}
	ok := u.isModel

	switch t := evt.(type) {
	case *sdl.MouseButtonEvent:
		return u.onButton(t, ok)
	case *sdl.MouseWheelEvent:
		if u.posInRect(t.X, t.Y) {
			return true
		}
		return ok
	case *sdl.MouseMotionEvent:
		if u.canMove && u.isLDown {
			u.doMove(t.XRel, t.YRel)
		}
		if u.posInRect(t.X, t.Y) {
			return true
		}
		return ok
	case *sdl.KeyDownEvent:
		return ok
	case *sdl.KeyUpEvent:
		return ok
	case *sdl.TextInputEvent:
		return ok
	}

	return false
}
