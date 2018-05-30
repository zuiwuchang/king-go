package command

import (
	"reflect"
)

type _CommanderSignal struct {
	signal    chan interface{}
	commnader ICommander
}

// NewCommanderSignal .
func NewCommanderSignal(signal chan interface{}, commnader ICommander) ICommanderSignal {
	return _CommanderSignal{
		signal:    signal,
		commnader: commnader,
	}
}

// Done 如果 command 未註冊 返回 errCommandUnknow 否則執行 ch <- command
func (c _CommanderSignal) Done(command interface{}) (e error) {
	commandType := reflect.TypeOf(command)
	if c.commnader.HanderType(commandType) == nil {
		e = NewErrCommandUnknowType(commandType)
		return
	}

	c.signal <- command
	return
}

// Close close(ch) 此後 不能在 執行 任何 Done 操作
func (c _CommanderSignal) Close() {
	close(c.signal)
	return
}

// Run 運行 commnad<-ch Execute(e) 直到 Execute 返回 錯誤 或 close(ch)
// 如果 Run 返回 錯誤 於是你調用了 Close 此時應該 繼續調用 Run 或 RunNull 以便 在使用 帶緩存的 chan 時 緩存 goroutine 能夠正常退出
func (c _CommanderSignal) Run() (e error) {
	for command := range c.signal {
		e = c.commnader.Execute(command)
		if e != nil {
			break
		}
	}
	return
}

// RunNull 只執行 commnad<-ch 而不執行 Execute(e)
// 通常只在 使用了 帶緩存的 chan 時 Close 後 讓 chan 的 goroutine 退出
func (c _CommanderSignal) RunNull() {
	for range c.signal {
	}
	return
}
