package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
)

const (
	UI_BTN_NORMAL = 0
	UI_BTN_NODOWN = 1
)

type UiBtuuon struct {
	Node

	normalTexture, downTexture *sdl.Texture
}

func NewUiBtuuon(normal, down *sdl.Texture) *UiBtuuon {
	ui := &UiBtuuon{
		normalTexture: normal,
		downTexture:   down,
	}
	ui.SetTexture(normal)
	return ui
}
