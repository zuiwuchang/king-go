package command

import (
	"fmt"
	"reflect"
)

// IsUnknow 返回 錯誤是否是 errCommandUnknow
func IsUnknow(e error) bool {
	return reflect.TypeOf(e) == _errCommandUnknowType
}

var _errCommandUnknowType = reflect.TypeOf(errCommandUnknow{})

type errCommandUnknow struct {
	errMessage string
}

func (e errCommandUnknow) Error() string {
	return e.errMessage
}

// NewErrCommandUnknowType .
func NewErrCommandUnknowType(commandType reflect.Type) (e error) {
	e = errCommandUnknow{
		errMessage: fmt.Sprint("command not registered : ", commandType),
	}
	return
}
