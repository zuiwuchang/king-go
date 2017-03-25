//動作橋段
package sdl2

import (
	"time"
)

//動作執行完一個週期後的 回調
type ActionCallBack func(node Object, a Action)

//施加到 演員的 動作定義
type Action interface {
	//執行動作
	DoAction(node Object, duration time.Duration)
	//釋放 動作
	Destory()
	//是否自動 釋放
	Auto() bool
	//返回一個動作副本
	Clone() Action
}
