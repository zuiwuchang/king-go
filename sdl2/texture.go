package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	"unsafe"
)

func LoadTexture(file string) (*sdl.Texture, error) {
	renderer := g_director.renderer
	surface, e := img.Load(file)
	if e != nil {
		return nil, e
	}

	//設置透明色
	key := sdl.MapRGB(surface.Format, 0, 0, 0)
	if e = surface.SetColorKey(1, key); e != nil {
		return nil, e
	}

	return renderer.CreateTextureFromSurface(surface)
}
func LoadTextureMem(b []byte) (*sdl.Texture, error) {
	renderer := g_director.renderer

	rw := sdl.RWFromMem(unsafe.Pointer(&b[0]), len(b))
	surface, e := img.Load_RW(rw, true)
	if e != nil {
		return nil, e
	}

	//設置透明色
	key := sdl.MapRGB(surface.Format, 0, 0, 0)
	if e = surface.SetColorKey(1, key); e != nil {
		return nil, e
	}
	return renderer.CreateTextureFromSurface(surface)
}
