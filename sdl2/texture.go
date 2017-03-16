package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
)

func LoadTexture(file string) (*sdl.Texture, error) {
	renderer := g_director.renderer
	return img.LoadTexture(renderer, file)
}
