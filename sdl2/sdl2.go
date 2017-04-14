//孤 使用 sdl2 封裝的一個 2d 遊戲引擎
package sdl2

import (
	"container/list"
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	ttf "github.com/veandco/go-sdl2/sdl_ttf"
	"log"
	"sync"
	"time"
)

const (
	FONT_DEFAULT_FILE = "NotoSansCJKtc-Regular.otf"
	FONT_DEFAULT_SIZE = 16
)

//導演 保存了遊戲 環境 運行狀況
type Director struct {
	//最大fps
	fps int
	//遊戲窗口
	window *sdl.Window
	//窗口 renderer
	renderer *sdl.Renderer
	//預設字體
	font *ttf.Font

	//場景鏈 僅最外層場景被渲染 響應事件
	scenes *list.List

	//事件列表
	events *list.List
	mutex  sync.Mutex

	//焦點 元素
	focus Object
}

func (d *Director) pushEvent(evt interface{}) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	d.events.PushBack(evt)
}

func (d *Director) pollEvent() sdl.Event {
	d.mutex.Lock()
	defer d.mutex.Unlock()

	element := d.events.Front()
	if element == nil {
		return nil
	}
	d.events.Remove(element)
	return element.Value
}

//唯一的 導演單例
var g_director *Director

//壓入 場景
func PushScene(scene Object) {
	g_director.scenes.PushBack(scene)
}

//彈出 場景
func PopScene() Object {
	element := g_director.scenes.Back()
	if element == nil {
		return nil
	}

	scene := element.Value.(Object)
	g_director.scenes.Remove(element)
	return scene
}

//替換當前 場景
func Replace(scene Object) Object {
	old := PopScene()
	PushScene(scene)
	return old
}

//返回當前 場景
func GetScene() Object {
	element := g_director.scenes.Back()
	if element == nil {
		return nil
	}

	return element.Value.(Object)
}

//初始化 引擎
func InitEngine() error {
	if g_director != nil {
		return nil
	}
	g_director = &Director{scenes: list.New(),
		events: list.New(),
	}

	img.Init(img.INIT_JPG | img.INIT_PNG | img.INIT_TIF | img.INIT_WEBP)
	if e := initTTF(); e != nil {
		return e
	}

	return nil
}
func initTTF() error {
	if e := ttf.Init(); e != nil {
		return e
	}

	font, e := ttf.OpenFont(FONT_DEFAULT_FILE, FONT_DEFAULT_SIZE)
	if e != nil {
		return e
	}
	getDirector().font = font
	return nil
}
func getDirector() *Director {
	return g_director
}

//銷毀 引擎資源
func DestoryEngine() {
	director := getDirector()

	for element := director.scenes.Back(); element != nil; element = element.Prev() {
		scene := element.Value.(Object)
		scene.Destroy()
	}

	DestoryWindow()
	director.font.Close()
	g_director = nil
}

//創建一個 遊戲窗口
func CreateWindow(title string, w, h, fps int) error {
	//創建 窗口
	window, e := sdl.CreateWindow(title,
		sdl.WINDOWPOS_UNDEFINED,
		sdl.WINDOWPOS_UNDEFINED,
		w,
		h,
		0)
	if e != nil {
		return e
	}
	director := getDirector()
	director.fps = fps

	//創建 renderer
	renderer, e := sdl.CreateRenderer(window, -1, 0)
	if e != nil {
		window.Destroy()
		return e
	}

	director.window = window
	director.renderer = renderer

	return nil
}

//銷毀 遊戲 窗口
func DestoryWindow() {
	director := getDirector()
	if director.renderer != nil {
		director.renderer.Destroy()
		director.renderer = nil
	}

	if director.window != nil {
		director.window.Destroy()
		director.window = nil
	}
}

//顯示 測試 數據
func drawShow(fps int) {
	director := getDirector()
	font := director.font
	surface, e := font.RenderUTF8_Blended(
		fmt.Sprint("fps : ", fps),
		sdl.Color{R: 128, G: 128, B: 128},
	)
	if e != nil {
		log.Println(e)
		return
	}
	w, _ := director.window.GetSize()
	rect := sdl.Rect{X: int32(w) - 7 - surface.W, Y: 7, W: surface.W, H: surface.H}

	renderer := director.renderer
	texture, e := renderer.CreateTextureFromSurface(surface)
	if e != nil {
		log.Println(e)
		return
	}
	defer texture.Destroy()

	renderer.Copy(
		texture,
		nil,
		&rect,
	)
}

//渲染遊戲 運行邏輯
func Run(r, g, b, a uint8, show bool /*是否顯示調試信息*/) {
	director := getDirector()
	fps := director.fps
	renderer := director.renderer

	//主邏輯循環
	last := time.Now()
	var wait time.Duration
	if fps > 0 {
		wait = time.Second / time.Duration(fps)
	}
	nowFPS := 0
	lastFPS := last
	for {
		//獲取 擴展 事件
		for evt := director.pollEvent(); evt != nil; evt = director.pollEvent() {
			//處理 sdl 事件
			switch evt.(type) {
			case *sdl.QuitEvent:
				goto END
			default:
				if scene := GetScene(); scene != nil {
					scene.OnEvent(evt)
				}
			}
		}
		//獲取 sdl 事件
		for evt := sdl.PollEvent(); evt != nil; evt = sdl.PollEvent() {
			//處理 sdl 事件
			switch evt.(type) {
			case *sdl.QuitEvent:
				goto END
			default:
				if scene := GetScene(); scene != nil {
					scene.OnEvent(evt)
				}
			}
		}

		//保證 fps 省略掉多餘繪製
		now := time.Now()
		if wait > 0 && now.Before(last.Add(wait)) {
			continue
		}
		duration := now.Sub(last)

		last = now

		/***  繪製遊戲  ***/
		renderer.SetDrawColor(r, g, b, a)
		//清空背景
		renderer.Clear()

		//繪製 元素到 renderer
		if scene := GetScene(); scene != nil {
			scene.OnAction(duration)
			scene.Draw(renderer, duration)
		}

		if show {
			nowFPS++
			drawShow(fps)
			if now.After(lastFPS.Add(time.Second * 1)) {
				lastFPS = now
				fps = nowFPS
				nowFPS = 0
			}
		}
		//將 renderer 更新到屏幕
		renderer.Present()

	}
END:
}

//設置獲取焦點元素
func SetFocus(obj Object) {
	director := g_director
	//元素已經得到焦點 無需額外操作
	if obj == director.focus {
		return
	}

	//保存原焦點元素
	old := director.focus

	//設置新焦點
	director.focus = obj

	//發送 失去焦點事件
	if old != nil {
		director.pushEvent(&FocusOutEvent{Obj: old})
	}

	//發送 獲取焦點事件
	if obj != nil {
		director.pushEvent(&FocusInEvent{Obj: obj})
	}

}

//返回 遊戲 窗口
func GetWindow() *sdl.Window {
	if g_director == nil {
		return nil
	}
	return g_director.window
}

//壓入一個 事件
func PushEvent(evt interface{}) {
	g_director.pushEvent(evt)
}
