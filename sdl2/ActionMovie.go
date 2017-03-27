package sdl2

import (
	"github.com/veandco/go-sdl2/sdl"
	"time"
)

//播放動畫
type ActionMovie struct {
	//是否循環播放
	loop bool
	//動畫幀
	textures []*sdl.Texture
	//原紋理
	texture *sdl.Texture
	//是否 保存 原紋理
	textureOk bool

	//當前播放 幀
	pos float64

	//花費時間
	duration time.Duration

	//速度
	speed   float64
	speedOk bool

	//動作結束回調
	callback ActionCallBack
	params   interface{}
}

//n 預計紋理數量 duration 花費時間
func NewActionMovie(n int, duration time.Duration) *ActionMovie {
	return &ActionMovie{textures: make([]*sdl.Texture, 0, n),
		pos:      0,
		duration: duration,
	}
}

//執行動作
func (a *ActionMovie) DoAction(node Object, duration time.Duration) {
	a.calculateSpeed()

	size := len(a.textures)

	//ok
	pos := int(a.pos)

	if pos >= size {
		//不循環 移除 動作
		if !a.loop {
			node.RemoveAction(a)
			//還原紋理
			node.SetTexture(a.texture)
		}
		//call
		if a.callback != nil {
			a.callback(node, a, a.params)
		}

		if a.loop {
			a.pos = 0
		} else {
			return
		}
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
	action.speedOk = false
	return &action
}

//增加一個 動畫幀
func (a *ActionMovie) PushFrame(texture *sdl.Texture) {
	if texture != nil {
		a.textures = append(a.textures, texture)
	}
}

//計算播放速度
func (a *ActionMovie) calculateSpeed() {
	if a.speedOk {
		return
	}

	size := len(a.textures)
	if size == 0 {
		return
	}
	a.speed = float64(size) / float64(a.duration/time.Millisecond)

	a.speedOk = true
}

//設置 action 完成 通知
func (a *ActionMovie) SetCallBack(callback ActionCallBack, params interface{}) Action {
	a.callback = callback
	a.params = params
	return a
}

//返回 action 完成 通知
func (a *ActionMovie) GetCallBack() (ActionCallBack, interface{}) {
	return a.callback, a.params
}

//返回 是否 循環執行
func (a *ActionMovie) GetLoop() bool {
	return a.loop
}

//設置 是否 循環執行
func (a *ActionMovie) SetLoop(yes bool) Action {
	a.loop = yes
	return a
}
