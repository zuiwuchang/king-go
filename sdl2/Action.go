//動作橋段
package sdl2

import (
	"time"
)

//動作執行完一個週期後的 回調
type ActionCallBack func(node Object, a Action, params interface{})

//施加到 演員的 動作定義
type Action interface {
	//執行動作
	DoAction(node Object, duration time.Duration)
	//釋放 動作
	Destory()
	//是否自動 釋放
	//返回 true 移除action時 自動調用 a.Destory()
	Auto() bool
	//返回一個動作副本
	Clone() Action

	//設置 action 完成 通知
	SetCallBack(callback ActionCallBack, params interface{}) Action
	//返回 action 完成 通知
	GetCallBack() (ActionCallBack, interface{})

	//返回 是否 循環執行
	GetLoop() bool
	//設置 是否 循環執行
	SetLoop(yes bool) Action
}
