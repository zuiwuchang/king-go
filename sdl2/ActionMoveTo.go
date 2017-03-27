package sdl2

import (
	"time"
)

//勻速 移動到 指定 坐標
type ActionMoveTo struct {
	//目的坐標
	x, y float64

	//花費時間
	duration time.Duration
	//是否得到速度
	speedOk bool

	//速度
	speedX, speedY float64

	//動作結束回調
	callback ActionCallBack
	params   interface{}

	//是否循環執行
	loop bool
}

//x,y 當前坐標 targetX, targetY 目標坐標 duration 花費時間
func NewActionMoveTo(x, y float64, duration time.Duration) *ActionMoveTo {
	return &ActionMoveTo{x: x,
		y:        y,
		duration: duration,
	}
}

//計算速度
func (a *ActionMoveTo) calculateSpeed(node Object) {
	if a.speedOk {
		return
	}
	x, y := node.GetPos()

	x -= a.x
	y -= a.y
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}

	n := float64(a.duration / time.Millisecond)

	a.speedX = x / n
	a.speedY = y / n
	a.speedOk = true
}

//執行動作
func (a *ActionMoveTo) DoAction(node Object, duration time.Duration) {
	a.calculateSpeed(node)

	n := float64(duration / time.Millisecond)
	nx := n * a.speedX
	ny := n * a.speedY
	x, y := node.GetPos()
	if x == a.x && y == a.y {
		//不循環 移除 動作
		if !a.loop {
			node.RemoveAction(a)
		}

		if a.callback != nil {
			a.callback(node, a, a.params)
		}
		return
	}

	if x < a.x {
		x += nx
		if x > a.x {
			x = a.x
		}
	} else if x > a.x {
		x -= nx
		if x < a.x {
			x = a.x
		}
	}
	if y < a.y {
		y += ny
		if y > a.y {
			y = a.y
		}
	} else if y > a.y {
		y -= ny
		if y < a.y {
			y = a.y
		}
	}
	node.SetPos(x, y)
}

//釋放 動作
func (a *ActionMoveTo) Destory() {
	*a = ActionMoveTo{}
}

//是否自動 釋放
func (a *ActionMoveTo) Auto() bool {
	return false
}

//返回一個動作副本
func (a *ActionMoveTo) Clone() Action {
	action := *a
	action.speedOk = false
	return &action
}

//設置 action 完成 通知
func (a *ActionMoveTo) SetCallBack(callback ActionCallBack, params interface{}) Action {
	a.callback = callback
	a.params = params
	return a
}

//返回 action 完成 通知
func (a *ActionMoveTo) GetCallBack() (ActionCallBack, interface{}) {
	return a.callback, a.params
}

//返回 是否 循環執行
func (a *ActionMoveTo) GetLoop() bool {
	return a.loop
}

//設置 是否 循環執行
func (a *ActionMoveTo) SetLoop(yes bool) Action {
	a.loop = yes
	return a
}
