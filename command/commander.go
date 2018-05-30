package command

import (
	"fmt"
	"reflect"
	"strings"
)

type _Commander struct {
	keys map[reflect.Type]CommanderHander
}

// New 創建一個 ICommander
func New() ICommander {
	return _Commander{
		keys: make(map[reflect.Type]CommanderHander),
	}
}

// RegisterCommander 自動將 handerST 中的 處理函數 註冊到 commander
//
// 處理函數 必須是名稱以 prefix 開始 的導出函數 func (xxx)[Prefix]XXX(commandType)(error)
//
// handerST 必須是一個 struct 或 *struct 否則將 panics
func RegisterCommander(commander ICommander, handerST interface{}, prefix string) {
	t := reflect.TypeOf(handerST)
	prefix = strings.TrimSpace(prefix)

	// 獲取 所有 導出 函數
	for i := 0; i < t.NumMethod(); i++ {
		m := t.Method(i)
		mt := m.Type

		// 驗證 函數簽名
		if mt.NumIn() != 2 ||
			mt.NumOut() != 1 ||
			mt.Out(0) != errInterface ||
			!strings.HasPrefix(m.Name, prefix) {
			continue
		}

		registerCommander(commander, handerST, m)
	}
}
func registerCommander(commander ICommander, handerST interface{}, m reflect.Method) {
	commander.RegisterType(m.Type.In(1), func(command interface{}) (e error) {
		rs := m.Func.Call(
			[]reflect.Value{
				reflect.ValueOf(handerST),
				reflect.ValueOf(command),
			},
		)[0].Interface()
		if rs != nil {
			e = rs.(error)
		}
		return
	})
}

func (c _Commander) String() string {
	return fmt.Sprint(c.keys)
}

// Done 執行一個 命令
// 如果 command 未註冊 返回 errCommandUnknow
func (c _Commander) Done(command interface{}) (e error) {
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

// NumHander 返回 Hander 數量
func (c _Commander) NumHander() int {
	return len(c.keys)
}
