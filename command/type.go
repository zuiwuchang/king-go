package command

import (
	"reflect"
)

// CommanderHander 命令處理 函數
type CommanderHander func(command interface{}) (e error)

// ICommander 命令執行器
type ICommander interface {
	// Execute 執行一個 命令
	// 如果 command 未註冊 返回 errCommandUnknow
	Execute(command interface{}) (e error)

	// HanderType 返回 commandType 的處理 Hander 如果不存在 返回 nil
	HanderType(commandType reflect.Type) (f CommanderHander)
	// Hander HanderType(reflect.TypeOf(command)) 的語法糖
	Hander(command interface{}) (f CommanderHander)

	// Register 註冊 一個 命令
	// 如果 f == nil 則 註銷此 處理器
	RegisterType(commandType reflect.Type, f CommanderHander)

	// Register RegisterType(reflect.TypeOf(command),f) 的語法糖
	Register(command interface{}, f CommanderHander)
}

// ICommanderSignal 一個 帶有 ICommander 的 chan
type ICommanderSignal interface {
	// Done 如果 command 未註冊 返回 errCommandUnknow 否則執行 ch <- command
	Done(command interface{}) (e error)

	// Close close(ch) 此後 不能在 執行 任何 Done 操作
	Close()

	// Run 運行 commnad<-ch Execute(e) 直到 Execute 返回 錯誤 或 close(ch)
	// 如果 Run 返回 錯誤 於是你調用了 Close 此時應該 繼續調用 Run 或 RunNull 以便 在使用 帶緩存的 chan 時 緩存 goroutine 能夠正常退出
	Run() (e error)

	// RunNull 只執行 commnad<-ch 而不執行 Execute(e)
	// 通常只在 使用了 帶緩存的 chan 時 Close 後 讓 chan 的 goroutine 退出
	RunNull()
}
