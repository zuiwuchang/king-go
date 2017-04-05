package sdl2

import (
	"time"
)

//勻速 改變 Alpha 到指定值
type ActionAlphaTo struct {
	ActionBase

	//目的 Alpha
	alpha int

	//花費時間
	duration time.Duration

	//是否得到速度
	speedOk bool
	//速度
	speed float64
}

//alpha 目標alpha duration 花費時間
func NewActionAlphaTo(alpha uint8, duration time.Duration) *ActionAlphaTo {
	return &ActionAlphaTo{
		alpha:    int(alpha),
		duration: duration,
	}
}

//計算速度
func (a *ActionAlphaTo) calculateSpeed(node Object) {
	if a.speedOk {
		return
	}
	alpha := int(node.GetAlpha())
	size := a.alpha - alpha

	a.speed = float64(size) / float64(a.duration/time.Millisecond)

	a.speedOk = true
}

//執行動作
func (a *ActionAlphaTo) DoAction(node Object, duration time.Duration) {
	a.calculateSpeed(node)

	//ok
	pos := int(node.GetAlpha())

	if pos == a.alpha {
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
	pos += int(a.speed * float64(duration/time.Millisecond))
	if a.speed < 0 {
		if pos < a.alpha {
			pos = a.alpha
		}
	} else if a.speed > 0 {
		if pos > a.alpha {
			pos = a.alpha
		}
	}
	node.SetAlpha(uint8(pos))
}

//釋放 動作
func (a *ActionAlphaTo) Destory() {
	*a = ActionAlphaTo{}
}

//返回一個動作副本
func (a *ActionAlphaTo) Clone() Action {
	action := *a
	action.speedOk = false
	return &action
}
