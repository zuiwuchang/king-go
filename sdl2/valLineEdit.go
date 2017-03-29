package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
	"king-go/algorithm"
	"log"
	"time"
)

//單行文本 值
type valLineEdit struct {
	//字體
	font *ttf.Font
	//文本顏色
	color sdl.Color

	//文本字符串 緩存
	text string

	//文本紋理
	texture *sdl.Texture

	//光標位置 [0,len(text)]
	chartBegin      int
	chartEnd        int
	chartBeginPixel int
	chartEndPixel   int

	//光標紋理
	chartTexture *sdl.Texture

	//插入符 顏色
	r, g, b   uint8
	lastChart time.Time

	//是否 選擇文本中
	isSelect bool
}

func (v *valLineEdit) GetChartRGB() (r uint8, g uint8, b uint8) {
	return v.r, v.g, v.b
}
func (v *valLineEdit) SetChartRGB(r, g, b uint8) {
	v.r = r
	v.g = g
	v.b = b

	if v.chartTexture != nil {
		v.chartTexture.Destroy()
		v.chartTexture = nil
	}
}

//返回 光標
func (v *valLineEdit) getChart(x int32) (int, int) {
	font := v.font
	arrs := []rune(v.text)
	size := len(arrs)
	pos := 0

	n, e := algorithm.BinarySearch(0, size-1, func(i int) (int, error) {
		w, _, e := font.SizeUTF8(string(arrs[:i+1]))

		if e != nil {
			return 0, e
		}
		if int32(w) < x {
			if i+1 == size {
				pos = w
				return 0, nil
			}
			return -1, nil
		}

		if i == 0 {
			pos = 0
			return 0, nil
		}

		w2, _, e := font.SizeUTF8(string(arrs[:i]))

		if e != nil {
			return 0, e
		}
		if int32(w2) < x {
			pos = w2
			return 0, nil
		}
		return 1, nil
	})
	if e != nil {
		log.Println(e)
		return 0, 0
	}
	return n, pos
}

//返回是否 正在 文本選擇中
func (v *valLineEdit) IsSelect() bool {
	return v.isSelect
}

//選中文本
func (v *valLineEdit) SelectIng(x int32) {
	n, pixel := v.getChart(x)

	v.chartEnd = n
	v.chartEndPixel = pixel
}

//開始 選擇文本
func (v *valLineEdit) SelectStart(x int32) {
	v.isSelect = true

	n, pixel := v.getChart(x)

	v.chartBegin = n
	v.chartBeginPixel = pixel
	v.chartEnd = n
	v.chartEndPixel = pixel
}

//結束 選擇文本
func (v *valLineEdit) SelectStop(x int32) {
	if !v.isSelect {
		return
	}

	v.isSelect = false
}

//銷毀 資源
func (v *valLineEdit) Destroy() {
	if v.font != nil {
		v.font.Close()
	}
	if v.texture != nil {
		v.texture.Destroy()
	}
	if v.chartTexture != nil {
		v.chartTexture.Destroy()
	}
	*v = valLineEdit{}
}

//返回 字符串值
func (v *valLineEdit) GetString() string {
	return v.text
}

//設置 字符串值
func (v *valLineEdit) SetString(str string, width int) {
	if v.texture != nil {
		v.texture.Destroy()
		v.texture = nil
	}

	if str == "" {
		v.text = ""
		//設置chart

		v.chartBegin = 0
		v.chartEnd = 0
		v.chartBeginPixel = 0
		v.chartEndPixel = 0
		return
	}

	//尋找 最大可顯示文本
	font := v.font
	arrs := []rune(str)
	size := len(arrs)
	pos := 0

	n, e := algorithm.BinarySearch(0, size-1, func(i int) (int, error) {
		w, _, e := font.SizeUTF8(string(arrs[:i+1]))

		if e != nil {
			return 0, e
		}
		if w >= width {
			return 1, nil
		}
		if i+1 == size {
			return 0, nil
		}

		w2, _, e := font.SizeUTF8(string(arrs[:i+2]))

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
		log.Println(e)
		return
	}

	v.text = string(arrs[:n+1])
	v.chartBegin = n
	v.chartEnd = n
	v.chartBeginPixel = pos
	v.chartEndPixel = pos
}

func (v *valLineEdit) initTexture(renderer *sdl.Renderer, x, y, w, h int32) {
	if v.texture != nil {
		return
	}

	surface, e := v.font.RenderUTF8_Blended(v.text, v.color)
	if e != nil {
		log.Println(e)
		return
	}
	surfaceTarget, e := sdl.CreateRGBSurface(0,
		w,
		h,
		32,
		R_MASK,
		G_MASK,
		B_MASK,
		A_MASK,
	)
	if e != nil {
		log.Println(e)
		return
	}

	src := sdl.Rect{}
	if surface.W < surfaceTarget.W {
		src.W = surface.W
	} else {
		src.W = surfaceTarget.W
	}
	if surface.H < surfaceTarget.H {
		src.H = surface.H
	} else {
		src.H = surfaceTarget.H
	}

	surface.Blit(&src, surfaceTarget, &src)

	texture, e := renderer.CreateTextureFromSurface(surfaceTarget)
	if e != nil {
		log.Println(e)
		return
	}
	v.texture = texture
}

//繪製文本
func (v *valLineEdit) DrawText(renderer *sdl.Renderer, x, y, w, h int32) {
	if v.text == "" {
		return
	}

	v.initTexture(renderer, x, y, w, h)
	if v.texture != nil {
		renderer.Copy(v.texture,
			nil,
			&sdl.Rect{X: x, Y: y, W: w, H: h},
		)
	}
}
func (v *valLineEdit) DrawChart(renderer *sdl.Renderer, x, y, w, h int32) {
	if v.chartBegin == v.chartEnd {
		now := time.Now()
		if v.lastChart.IsZero() {
			v.lastChart = now
		} else {
			if now.Before(v.lastChart.Add(time.Millisecond * 600)) {
				//draw
			} else if now.Before(v.lastChart.Add(time.Millisecond * 600 * 2)) {
				return
			} else {
				v.lastChart = now
			}
		}

		renderer.SetDrawColor(v.r, v.g, v.b, 255)
		renderer.FillRect(&sdl.Rect{
			X: x + int32(v.chartBeginPixel),
			Y: y,
			W: w,
			H: h,
		})
	} else {
		v.initChartTexture(renderer)

		var posX, w int
		if v.chartBegin < v.chartEnd {
			posX = v.chartBeginPixel
			w = v.chartEndPixel - posX
		} else {
			posX = v.chartEndPixel
			w = v.chartBeginPixel - posX
		}

		renderer.Copy(v.chartTexture,
			nil,
			&sdl.Rect{
				X: x + int32(posX),
				Y: y,
				W: int32(w),
				H: h,
			},
		)

	}
}
func (v *valLineEdit) initChartTexture(renderer *sdl.Renderer) {
	if v.chartTexture != nil {
		return
	}
	surface, e := sdl.CreateRGBSurface(0,
		20,
		20,
		32,
		R_MASK,
		G_MASK,
		B_MASK,
		A_MASK,
	)
	if e != nil {
		log.Println(e)
		return
	}
	surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: 20, H: 20},
		sdl.Color{R: v.r, G: v.g, B: v.b, A: 160}.Uint32(),
	)

	if texture, e := renderer.CreateTextureFromSurface(surface); e != nil {
		log.Println(e)
	} else {
		v.chartTexture = texture
	}
}
