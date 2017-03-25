package sdl2

import (
	"fmt"
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

//播放動畫
type ActionMovie struct {
	//是否循環播放
	Loop bool
	//動畫幀
	textures []*sdl.Texture
	//原紋理
	texture *sdl.Texture
	//是否 保存 原紋理
	textureOk bool

	//當前播放 幀
	pos float64

	//速度
	speed   float64
	speedOk bool

	//動作結束回調
	callback ActionCallBack
	params   interface{}
}

//n 預計紋理數量 duration 花費時間
func NewActionMovie(n int) *ActionMovie {
	return &ActionMovie{textures: make([]*sdl.Texture, n),
		pos: 0,
	}
}

//執行動作
func (a *ActionMovie) DoAction(node Object, duration time.Duration) {
	size := len(a.textures)

	//ok
	pos := int(a.pos)

	if pos >= size {
		//不循環 移除 動作
		if !a.Loop {
			node.RemoveAction(a)
			//還原紋理
			node.SetTexture(a.texture)
			return
		}
		//call
		if a.callback != nil {
			a.callback(node, a, a.params)
		}
		a.pos = 0
	}

	//沒保存 原紋理
	if !a.textureOk {
		a.textureOk = true
		a.texture = node.GetTexture()
	}
	//fmt.Println(a.speed * float64(duration/time.Millisecond))
	a.pos += a.speed * float64(duration/time.Millisecond)
	pos = int(a.pos)
	if pos >= size {
		pos = size - 1
	}
	node.SetTexture(a.textures[pos])
}

//釋放 動作
func (a *ActionMovie) Destory() {

}

//是否自動 釋放
func (a *ActionMovie) Auto() bool {
	return false
}

//返回一個動作副本
func (a *ActionMovie) Clone() Action {
	action := *a
	action.pos = 0
	action.texture = nil
	action.textureOk = false
	return &action
}

//增加一個 動畫幀
func (a *ActionMovie) PushFrame(texture *sdl.Texture) {
	if texture != nil {
		a.textures = append(a.textures, texture)
	}
}

//計算播放速度
func (a *ActionMovie) calculateSpeed(duration time.Duration) {
	size := len(a.textures)
	if size == 0 {
		return
	}
	a.speed = float64(size) / float64(duration/time.Millisecond)
	fmt.Println(a.speed)
}

//設置 action 完成 通知
func (a *ActionMovie) SetCallBack(callback ActionCallBack, params interface{}) {
	a.callback = callback
	a.params = params
}

//返回 action 完成 通知
func (a *ActionMovie) GetCallBack() (ActionCallBack, interface{}) {
	return a.callback, a.params
}
