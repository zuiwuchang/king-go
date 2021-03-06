package sdl2

import (
	"bytes"
	"errors"
	"github.com/veandco/go-sdl2/sdl"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
	"king-go/algorithm"
	"strings"
	"time"
)

//單行文本 值
type valLineEdit struct {
	//字體
	font *ttf.Font
	//文本顏色
	color sdl.Color

	//文本字符串 緩存
	str []rune

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

	//最大允許輸入字符數
	max int

	//密碼框 顯示數據
	c     rune
	isPwd bool
}

func (v *valLineEdit) resetStr() {
	v.str = make([]rune, 0)
}
func (v *valLineEdit) renderUTF8(arrs []rune) (*sdl.Surface, error) {
	if v.isPwd {
		var buf bytes.Buffer
		for i := 0; i < len(arrs); i++ {
			buf.WriteRune(v.c)
		}
		return v.font.RenderUTF8_Blended(buf.String(), v.color)
	}

	return v.font.RenderUTF8_Blended(string(arrs), v.color)
}
func (v *valLineEdit) sizeUTF8(arrs []rune) (int, error) {
	if v.isPwd {
		size := len(arrs)
		arrs = make([]rune, size, size)
		for i := 0; i < size; i++ {
			arrs[i] = v.c
		}
	}
	w, _, e := v.font.SizeUTF8(string(arrs))
	if e != nil {
		return 0, e
	}
	return w, nil
}
func (v *valLineEdit) IsPwd() bool {
	return v.isPwd
}
func (v *valLineEdit) SetPwd(yes bool) {
	if v.isPwd == yes {
		return
	}
	v.isPwd = yes
	v.destroyTexture()
	v.chartBegin = 0
	v.chartEnd = 0
	v.chartBeginPixel = 0
	v.chartEndPixel = 0
}
func (v *valLineEdit) SetPwdChar(c string) {
	str := strings.TrimSpace(c)
	if str != "" {
		v.c = []rune(str)[0]
	}

}
func (v *valLineEdit) GetPwdChar() string {
	return string(v.c)
}
func (v *valLineEdit) GetMax() int {
	return v.max
}
func (v *valLineEdit) SetMax(max int) {
	v.max = max
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
	arrs := v.str
	size := len(arrs)
	pos := 0
	offset := 0
	n, e := algorithm.BinarySearch(0, size-1, func(i int) (int, error) {
		w, e := v.sizeUTF8(arrs[:i+1])

		if e != nil {
			return 0, e
		}
		if int32(w) < x {
			if i+1 == size {
				pos = w
				offset = 1
				return 0, nil
			}
			return -1, nil
		}

		if i == 0 {
			pos = 0
			return 0, nil
		}

		w2, e := v.sizeUTF8(arrs[:i])

		if e != nil {
			return 0, e
		}
		if int32(w2) < x {
			pos = w2
			return 0, nil
		}
		return 1, nil
	})
	if e != nil && e.Error() != "not found" {
		g_log.Println(e)
		return 0, 0
	}
	return n + offset, pos
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

	status := sdl.GetKeyboardState()
	if status[sdl.SCANCODE_LSHIFT] != 0 ||
		status[sdl.SCANCODE_RSHIFT] != 0 {
		v.chartEnd = n
		v.chartEndPixel = pixel
	} else {

		v.chartBegin = n
		v.chartBeginPixel = pixel
		v.chartEnd = n
		v.chartEndPixel = pixel
	}
}

//結束 選擇文本
func (v *valLineEdit) SelectStop(x int32) {
	if !v.isSelect {
		return
	}

	v.isSelect = false
}

func (v *valLineEdit) getPos(n int) int {
	if n == 0 {
		return 0
	}
	arrs := v.str
	size := len(arrs)
	if n > size {
		n = size
	}
	w, e := v.sizeUTF8(arrs[:n])
	if e != nil {
		g_log.Println(e)
	}
	return w
}
func (v *valLineEdit) SelectLeft() {
	status := sdl.GetKeyboardState()
	if status[sdl.SCANCODE_LSHIFT] != 0 ||
		status[sdl.SCANCODE_RSHIFT] != 0 {
		if v.chartEnd > 0 {
			v.chartEnd--

			pos := v.getPos(v.chartEnd)
			v.chartEndPixel = pos
		}
	} else {
		if v.chartBegin == v.chartEnd {
			if v.chartBegin > 0 {
				v.chartBegin--
				v.chartEnd--

				pos := v.getPos(v.chartBegin)
				v.chartBeginPixel = pos
				v.chartEndPixel = pos
			}
		} else if v.chartBegin < v.chartEnd {
			v.chartEnd = v.chartBegin
			v.chartEndPixel = v.chartBeginPixel
		} else {
			v.chartBegin = v.chartEnd
			v.chartBeginPixel = v.chartEndPixel
		}
	}
}
func (v *valLineEdit) SelectRight() {
	status := sdl.GetKeyboardState()
	if status[sdl.SCANCODE_LSHIFT] != 0 ||
		status[sdl.SCANCODE_RSHIFT] != 0 {
		arrs := v.str
		size := len(arrs)
		if v.chartEnd < size {
			v.chartEnd++

			pos := v.getPos(v.chartEnd)
			v.chartEndPixel = pos
		}
	} else {
		if v.chartBegin == v.chartEnd {
			arrs := v.str
			size := len(arrs)
			if v.chartBegin < size {
				v.chartBegin++
				v.chartEnd++

				pos := v.getPos(v.chartBegin)
				v.chartBeginPixel = pos
				v.chartEndPixel = pos
			}
		} else if v.chartBegin > v.chartEnd {
			v.chartEnd = v.chartBegin
			v.chartEndPixel = v.chartBeginPixel
		} else {
			v.chartBegin = v.chartEnd
			v.chartBeginPixel = v.chartEndPixel
		}
	}
}
func (v *valLineEdit) destroyTexture() {
	if v.texture != nil {
		v.texture.Destroy()
		v.texture = nil
	}
}

//銷毀 資源
func (v *valLineEdit) Destroy() {
	if v.font != nil {
		v.font.Close()
	}
	v.destroyTexture()
	if v.chartTexture != nil {
		v.chartTexture.Destroy()
	}
	*v = valLineEdit{}
}

//返回 字符串值
func (v *valLineEdit) GetString() string {
	return string(v.str)
}

//設置 字符串值
func (v *valLineEdit) SetString(str string, width int) {
	if v.max > 0 {
		arrs := []rune(str)
		if len(arrs) > v.max {
			str = string(arrs[:v.max])
		}
	}

	if v.texture != nil {
		v.texture.Destroy()
		v.texture = nil
	}

	if str == "" {
		v.resetStr()
		//設置chart

		v.chartBegin = 0
		v.chartEnd = 0
		v.chartBeginPixel = 0
		v.chartEndPixel = 0
		return
	}

	//尋找 最大可顯示文本
	arrs := []rune(str)
	size := len(arrs)
	pos := 0

	n, e := algorithm.BinarySearch(0, size-1, func(i int) (int, error) {
		w, e := v.sizeUTF8(arrs[:i+1])

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

		w2, e := v.sizeUTF8(arrs[:i+2])

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
		g_log.Println(e)
		return
	}

	v.str = arrs[:n+1]
	v.chartBegin = n + 1
	v.chartEnd = n + 1
	v.chartBeginPixel = pos
	v.chartEndPixel = pos
}

func (v *valLineEdit) initTexture(renderer *sdl.Renderer, x, y, w, h int32) {
	if v.texture != nil {
		return
	}

	surface, e := v.renderUTF8(v.str)
	if e != nil {
		g_log.Println(e)
		return
	}
	defer surface.Free()
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
		g_log.Println(e)
		return
	}
	defer surfaceTarget.Free()

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
		g_log.Println(e)
		return
	}
	v.texture = texture

}

//繪製文本
func (v *valLineEdit) DrawText(renderer *sdl.Renderer, x, y, w, h int32) {
	if len(v.str) == 0 {
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
		g_log.Println(e)
		return
	}
	defer surface.Free()
	surface.FillRect(&sdl.Rect{X: 0, Y: 0, W: 20, H: 20},
		sdl.Color{R: v.r, G: v.g, B: v.b, A: 160}.Uint32(),
	)

	if texture, e := renderer.CreateTextureFromSurface(surface); e != nil {
		g_log.Println(e)
	} else {
		v.chartTexture = texture
	}
}

//將 選中項 以 字符串 替換
func (v *valLineEdit) ReplaceStr(str string, width int) error {
	arrs := []rune(str)
	return v.ReplaceRune(arrs, width)
}

//將 選中項 以 []rune 替換
func (v *valLineEdit) ReplaceRune(arrs []rune, width int) error {
	nSize := len(arrs)
	if v.chartBegin == v.chartEnd && nSize == 0 {
		return nil
	}

	var begin, end int
	if v.chartBegin < v.chartEnd {
		begin = v.chartBegin
		end = v.chartEnd
	} else {
		begin = v.chartEnd
		end = v.chartBegin
	}
	pos := begin + nSize

	old := v.str
	oSize := len(old)
	size := pos + (oSize - end)
	if v.max > 0 && size > v.max {
		return errors.New("text more max length")
	}
	nRune := make([]rune, size, size)

	copy(nRune, old[:begin])
	copy(nRune[begin:], arrs)
	copy(nRune[pos:], old[end:])

	w, e := v.sizeUTF8(nRune)
	if e != nil {
		return e
	}
	if w > width {
		return errors.New("text more max length")
	}
	w, e = v.sizeUTF8(nRune[:pos])
	if e != nil {
		return e
	}

	v.chartBegin = pos
	v.chartEnd = pos
	v.chartBeginPixel = w
	v.chartEndPixel = w

	v.str = nRune
	v.destroyTexture()
	return nil
}

//刪除
func (v *valLineEdit) Backspace() {
	arrs := v.str
	size := len(arrs)
	if size == 0 {
		return
	}

	if v.chartBegin == v.chartEnd {
		if v.chartBegin == 0 {
			return
		}
		n := v.chartBegin

		w, e := v.sizeUTF8(arrs[:n-1])
		if e != nil {
			g_log.Println(e)
			return
		}

		copy(arrs[n-1:], arrs[n:])
		v.str = arrs[:size-1]
		v.destroyTexture()
		v.chartBegin--
		v.chartEnd--
		v.chartBeginPixel = w
		v.chartEndPixel = w
	} else {
		var begin, end int
		if v.chartBegin < v.chartEnd {
			begin = v.chartBegin
			end = v.chartEnd
		} else {
			begin = v.chartEnd
			end = v.chartBegin
		}
		pos := begin
		copy(arrs[pos:], arrs[end:])
		arrs := arrs[:pos+size-end]

		w, e := v.sizeUTF8(arrs[:pos])
		if e != nil {
			g_log.Println(e)
			return
		}

		v.chartBegin = pos
		v.chartEnd = pos
		v.chartBeginPixel = w
		v.chartEndPixel = w

		v.str = arrs
		v.destroyTexture()
	}
}

//返回 選中文本
func (v *valLineEdit) GetSelectStr() (str string) {
	if v.chartBegin == v.chartEnd {
		return
	}

	var begin, end int
	if v.chartBegin < v.chartEnd {
		begin = v.chartBegin
		end = v.chartEnd
	} else {
		begin = v.chartEnd
		end = v.chartBegin
	}

	arrs := v.str
	str = string(arrs[begin:end])
	return
}

//光標 選擇
func (v *valLineEdit) Select(begin, end int) {
	size := len(v.str)

	if size == 0 {
		return
	}

	if begin > size {
		begin = size
	}
	if end > size {
		end = size
	}

	if begin > end {
		begin, end = end, begin
	}

	v.chartBegin = begin
	v.chartBeginPixel = v.getPos(begin)

	if begin == end {
		v.chartEnd = v.chartBegin
		v.chartEndPixel = v.chartBeginPixel
		return
	}
	v.chartEnd = end
	v.chartEndPixel = v.getPos(end)
}
