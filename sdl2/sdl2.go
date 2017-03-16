//孤 使用 sdl2 封裝的一個 2d 遊戲引擎
package sdl2

import (
	"container/list"
	"github.com/veandco/go-sdl2/sdl"
	img "github.com/veandco/go-sdl2/sdl_image"
	"time"
)

//導演 保存了遊戲 環境 運行狀況
type Director struct {
	//最大fps
	fps int
	//遊戲窗口
	window *sdl.Window
	//窗口 renderer
	renderer *sdl.Renderer

	//場景鏈 僅最外層場景被渲染
	scenes *list.List
}

//唯一的 導演單例
var g_director *Director

//壓入 場景
func PushScene(scene *Scene) {
	g_director.scenes.PushBack(scene)
}

//彈出 場景
func PopScene() *Scene {
	element := g_director.scenes.Back()
	if element == nil {
		return nil
	}

	scene := element.Value.(*Scene)
	g_director.scenes.Remove(element)
	return scene
}

//替換當前 場景
func Replace(scene *Scene) {
	PopScene()
	PushScene(scene)
}

//返回當前 場景
func GetScene() *Scene {
	element := g_director.scenes.Back()
	if element == nil {
		return nil
	}

	return element.Value.(*Scene)
}

//初始化 引擎
func InitEngine() error {
	if g_director != nil {
		return nil
	}
	g_director = &Director{scenes: list.New()}

	img.Init(img.INIT_JPG | img.INIT_PNG | img.INIT_TIF | img.INIT_WEBP)
	return nil
}
func getDirector() *Director {
	return g_director
}

//銷毀 引擎資源
func DestoryEngine() {
	for element := g_director.scenes.Back(); element != nil; element = element.Prev() {
		scene := element.Value.(*Scene)
		scene.Destroy()
	}

	DestoryWindow()
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

//渲染遊戲 運行邏輯
func Run() {
	director := getDirector()
	fps := director.fps
	renderer := director.renderer
	//主邏輯循環
	last := time.Now()
	for {
		//獲取事件
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
		if now.Before(last.Add(time.Second / time.Duration(fps))) {
			continue
		}
		duration := now.Sub(last)

		last = now

		/***  繪製遊戲  ***/
		//清空背景
		renderer.Clear()

		//繪製 元素到 renderer
		if scene := GetScene(); scene != nil {
			scene.Draw(renderer, duration)
		}

		//將 renderer 更新到屏幕
		renderer.Present()
	}
END:
}
