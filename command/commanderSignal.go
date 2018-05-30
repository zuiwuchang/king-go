package command

import (
	"reflect"
)

type _CommanderSignal struct {
	signal    chan interface{}
	commnader ICommander
}

// NewSignal .
func NewSignal(signal chan interface{}, commnader ICommander) ICommanderSignal {
	if signal == nil {
		signal = make(chan interface{})
	}
	if commnader == nil {
		commnader = New()
	}
	return _CommanderSignal{
		signal:    signal,
		commnader: commnader,
	}
}

// Commander 返回 綁定的 ICommander
func (c _CommanderSignal) Commander() ICommander {
	return c.commnader
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

// Run 運行 commnad<-ch Done(e) 直到 Done 返回 錯誤 或 close(ch)
// 如果 Run 返回 錯誤 於是你調用了 Close 此時應該 繼續調用 Run 或 RunNull 以便 在使用 帶緩存的 chan 時 緩存 goroutine 能夠正常退出
func (c _CommanderSignal) Run() (e error) {
	for command := range c.signal {
		e = c.commnader.Done(command)
		if e != nil {
			break
		}
	}
	return
}

// RunNull 只執行 commnad<-ch 而不執行 Done(e)
// 通常只在 使用了 帶緩存的 chan 時 Close 後 讓 chan 的 goroutine 退出
func (c _CommanderSignal) RunNull() {
	for range c.signal {
	}
	return
}
