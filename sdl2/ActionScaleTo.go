package sdl2

import (
	"time"
)

//勻速 改變 scale x y 到指定值
type ActionScaleTo struct {
	ActionBase

	//目的 scale
	scaleX float64
	scaleY float64

	//花費時間
	duration time.Duration

	//是否得到速度
	speedOk bool
	//速度
	speedX float64
	speedY float64
}

//scaleX scaleY 目標scale duration 花費時間
func NewActionScaleTo(scaleX, scaleY float64, duration time.Duration) *ActionScaleTo {
	return &ActionScaleTo{
		scaleX:   scaleX,
		scaleY:   scaleY,
		duration: duration,
	}
}

//計算速度
func (a *ActionScaleTo) calculateSpeed(node IObject) {
	if a.speedOk {
		return
	}
	scaleX, scaleY := node.GetScale()

	sizeX := a.scaleX - scaleX
	sizeY := a.scaleY - scaleY

	a.speedX = sizeX / float64(a.duration/time.Millisecond)
	a.speedY = sizeY / float64(a.duration/time.Millisecond)

	a.speedOk = true
}

func (a *ActionScaleTo) isOk(scaleX, scaleY float64) bool {
	return (a.speedX == 0 || a.scaleX == scaleX) &&
		(a.speedY == 0 || a.scaleY == scaleY)
}

//執行動作
func (a *ActionScaleTo) DoAction(node IObject, duration time.Duration) {
	a.calculateSpeed(node)

	//ok
	scaleX, scaleY := node.GetScale()

	if a.isOk(scaleX, scaleY) {
		//不循環 移除 動作
		if !a.loop {
			node.RemoveAction(a)
		}
		//call
		if a.callback != nil {
			a.callback(node, a, a.params)
		}
		return
	}
	if a.speedX != 0 {
		scaleX += a.speedX * float64(duration/time.Millisecond)
		if a.speedX < 0 {
			if scaleX < a.scaleX {
				scaleX = a.scaleX
			}
		} else if a.speedX > 0 {
			if scaleX > a.scaleX {
				scaleX = a.scaleX
			}
		}
		node.SetScaleX(scaleX)
	}
	if a.speedY != 0 {
		scaleY += a.speedY * float64(duration/time.Millisecond)
		if a.speedY < 0 {
			if scaleY < a.scaleY {
				scaleY = a.scaleY
			}
		} else if a.speedY > 0 {
			if scaleY > a.scaleY {
				scaleY = a.scaleY
			}
		}
		node.SetScaleY(scaleY)
	}
}

//釋放 動作
func (a *ActionScaleTo) Destory() {
	*a = ActionScaleTo{}
}

//返回一個動作副本
func (a *ActionScaleTo) Clone() IAction {
	action := *a
	action.speedOk = false
	return &action
}
