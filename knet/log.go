package knet

import (
	"github.com/zuiwuchang/king-go/log"
)

var kLog, kLogNil *log.Loggers

func init() {
	kLogNil = &log.Loggers{}
	kLog = kLogNil
}

//啓用所有默認 日誌 以便調試
func InitDebugLoggers() {
	kLog = log.NewDebugLoggers2("[knet] ")
}

//爲庫 設置自定義的 日誌
func SetLoggers(l *log.Loggers) {
	if l == nil {
		kLog = kLogNil
	} else {
		kLog = l
	}
}
