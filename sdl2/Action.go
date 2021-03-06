//動作橋段
package sdl2

import (
	"time"
)

//動作執行完一個週期後的 回調
type ActionCallBack func(node IObject, a IAction, params interface{})

//施加到 演員的 動作定義
type IAction interface {
	//執行動作
	DoAction(node IObject, duration time.Duration)
	//釋放 動作
	Destroy()
	//是否自動 釋放
	//返回 true 移除action時 自動調用 a.Destory()
	GetAutoDestroy() bool
	SetAutoDestory(yes bool)
	//返回一個動作副本
	Clone() IAction

	//設置 action 完成 通知
	SetCallBack(callback ActionCallBack, params interface{})
	//返回 action 完成 通知
	GetCallBack() (ActionCallBack, interface{})

	//返回 是否 循環執行
	GetLoop() bool
	//設置 是否 循環執行
	SetLoop(yes bool)
}
type ActionBase struct {
	//是否 重複 action
	loop bool

	//是否自動 destory
	auto bool

	//動作結束回調
	callback ActionCallBack
	params   interface{}
}

func (a *ActionBase) SetAutoDestory(yes bool) {
	a.auto = yes
}
func (a *ActionBase) SetCallBack(callback ActionCallBack, params interface{}) {
	a.callback = callback
	a.params = params

}
func (a *ActionBase) SetLoop(yes bool) {
	a.loop = yes
}

//釋放 動作
func (a *ActionBase) Destroy() {
}

//是否自動 釋放
//返回 true 移除action時 自動調用 a.Destroy()
func (a *ActionBase) GetAutoDestroy() bool {
	return a.auto
}

//返回 action 完成 通知
func (a *ActionBase) GetCallBack() (ActionCallBack, interface{}) {
	return a.callback, a.params
}

//返回 是否 循環執行
func (a *ActionBase) GetLoop() bool {
	return a.loop
}
