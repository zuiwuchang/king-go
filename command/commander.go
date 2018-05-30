package command

import (
	"reflect"
)

type _Commander struct {
	keys map[reflect.Type]CommanderHander
}

// NewCommander 創建一個 ICommander
func NewCommander() ICommander {
	return _Commander{
		keys: make(map[reflect.Type]CommanderHander),
	}
}

// Execute 執行一個 命令
// 如果 command 未註冊 返回 errCommandUnknow
func (c _Commander) Execute(command interface{}) (e error) {
	commandType := reflect.TypeOf(command)
	if f, _ := c.keys[commandType]; f == nil {
		e = NewErrCommandUnknowType(commandType)
	} else {
		e = f(command)
	}
	return
}

// HanderType 返回 commandType 的處理 Hander 如果不存在 返回 nil
func (c _Commander) HanderType(commandType reflect.Type) (f CommanderHander) {
	f, _ = c.keys[commandType]
	return
}

// Hander HanderType(reflect.TypeOf(command)) 的語法糖
func (c _Commander) Hander(command interface{}) (f CommanderHander) {
	commandType := reflect.TypeOf(command)
	f, _ = c.keys[commandType]
	return
}

// Register 註冊 一個 命令
// 如果 f == nil 則 註銷此 處理器
func (c _Commander) RegisterType(commandType reflect.Type, f CommanderHander) {
	if f == nil {
		delete(c.keys, commandType)
	} else {
		c.keys[commandType] = f
	}
}

// Register RegisterType(reflect.TypeOf(command),f) 的語法糖
func (c _Commander) Register(command interface{}, f CommanderHander) {
	commandType := reflect.TypeOf(command)
	if f == nil {
		delete(c.keys, commandType)
	} else {
		c.keys[commandType] = f
	}
}
