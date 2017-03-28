package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

type UiBtuuon struct {
	UiBase

	normalTexture, downTexture *sdl.Texture
	isLBtnDown                 bool
	isRBtnDown                 bool
	isMBtnDown                 bool

	lastLClice time.Time
	lastRClice time.Time
	lastMClice time.Time
}

func NewUiBtuuon(normal, down *sdl.Texture) *UiBtuuon {
	ui := &UiBtuuon{
		normalTexture: normal,
		downTexture:   down,
	}
	ui.SetTexture(normal)
	return ui
}
func (u *UiBtuuon) onLBtnEvt(t *sdl.MouseButtonEvent) bool {
	if t.Type == sdl.MOUSEBUTTONDOWN {
		if u.posInRect(t.X, t.Y) {
			u.isLBtnDown = true
			u.SetTexture(u.downTexture)

			SetFocus(u)
			return true
		}
	} else if t.Type == sdl.MOUSEBUTTONUP {
		if u.isLBtnDown {
			u.isLBtnDown = false
			if !u.isRBtnDown && !u.isMBtnDown {
				u.SetTexture(u.normalTexture)
			}

			if u.posInRect(t.X, t.Y) {
				callback := u.GetEventCallback(UI_EVT_LBTN_CLICK)
				if callback != nil {
					callback(u, nil)
				}

				now := time.Now()
				if now.Before(u.lastLClice.Add(time.Second)) {
					u.lastLClice = time.Time{}
					callback := u.GetEventCallback(UI_EVT_LBTN_DCLICK)
					if callback != nil {
						callback(u, nil)
					}
				} else {
					u.lastLClice = now
				}
			}
		}
	}

	return false
}
func (u *UiBtuuon) onRBtnEvt(t *sdl.MouseButtonEvent) bool {
	if t.Type == sdl.MOUSEBUTTONDOWN {
		if u.posInRect(t.X, t.Y) {
			u.isRBtnDown = true
			u.SetTexture(u.downTexture)

			SetFocus(u)
			return true
		}
	} else if t.Type == sdl.MOUSEBUTTONUP {
		if u.isRBtnDown {
			u.isRBtnDown = false
			if !u.isLBtnDown && !u.isMBtnDown {
				u.SetTexture(u.normalTexture)
			}
			if u.posInRect(t.X, t.Y) {
				callback := u.GetEventCallback(UI_EVT_RBTN_CLICK)
				if callback != nil {
					callback(u, nil)
				}

				now := time.Now()
				if now.Before(u.lastRClice.Add(time.Second)) {
					u.lastRClice = time.Time{}
					callback := u.GetEventCallback(UI_EVT_RBTN_DCLICK)
					if callback != nil {
						callback(u, nil)
					}
				} else {
					u.lastRClice = now
				}
			}
		}
	}

	return false
}
func (u *UiBtuuon) onMBtnEvt(t *sdl.MouseButtonEvent) bool {
	if t.Type == sdl.MOUSEBUTTONDOWN {
		if u.posInRect(t.X, t.Y) {
			u.isMBtnDown = true
			u.SetTexture(u.downTexture)

			SetFocus(u)
			return true
		}
	} else if t.Type == sdl.MOUSEBUTTONUP {
		if u.isMBtnDown {
			u.isMBtnDown = false
			if !u.isLBtnDown && !u.isRBtnDown {
				u.SetTexture(u.normalTexture)
			}
			if u.posInRect(t.X, t.Y) {
				callback := u.GetEventCallback(UI_EVT_MBTN_CLICK)
				if callback != nil {
					callback(u, nil)
				}

				now := time.Now()
				if now.Before(u.lastMClice.Add(time.Second)) {
					u.lastMClice = time.Time{}
					callback := u.GetEventCallback(UI_EVT_MBTN_DCLICK)
					if callback != nil {
						callback(u, nil)
					}
				} else {
					u.lastMClice = now
				}
			}
		}
	}

	return false
}
func (u *UiBtuuon) OnEvent(evt sdl.Event) bool {
	switch t := evt.(type) {
	case *sdl.MouseButtonEvent:
		switch t.Button {
		case sdl.BUTTON_LEFT:
			return u.onLBtnEvt(t)
		case sdl.BUTTON_RIGHT:
			return u.onRBtnEvt(t)
		case sdl.BUTTON_MIDDLE:
			return u.onMBtnEvt(t)
		}
	}
	return false
}
