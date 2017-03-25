package sdl2

import (
	"time"
)

//勻速 移動到 指定 坐標
type ActionMoveTo struct {
	//目的坐標
	x, y float64

	//速度
	speedX, speedY float64

	//動作結束回調
	callback ActionCallBack
	params   interface{}
}

//x,y 當前坐標 targetX, targetY 目標坐標 duration 花費時間
func NewActionMoveTo(x, y, targetX, targetY float64, duration time.Duration) *ActionMoveTo {
	x -= targetX
	y -= targetY
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}

	n := float64(duration / time.Millisecond)
	return &ActionMoveTo{x: targetX,
		y:      targetY,
		speedX: x / n,
		speedY: y / n,
	}
}

//執行動作
func (a *ActionMoveTo) DoAction(node Object, duration time.Duration) {
	n := float64(duration / time.Millisecond)
	nx := n * a.speedX
	ny := n * a.speedY
	x, y := node.GetPos()
	if x == a.x && y == a.y {
		node.RemoveAction(a)

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

}

//是否自動 釋放
func (a *ActionMoveTo) Auto() bool {
	return false
}

//返回一個動作副本
func (a *ActionMoveTo) Clone() Action {
	action := *a
	return &action
}

//設置 action 完成 通知
func (a *ActionMoveTo) SetCallBack(callback ActionCallBack, params interface{}) {
	a.callback = callback
	a.params = params
}

//返回 action 完成 通知
func (a *ActionMoveTo) GetCallBack() (ActionCallBack, interface{}) {
	return a.callback, a.params
}
