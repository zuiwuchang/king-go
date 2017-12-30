package log

import (
	kio "github.com/zuiwuchang/king-go/io"
	"log"
	"os"
)

var Trace *log.Logger
var Debug *log.Logger
var Info *log.Logger
var Warn *log.Logger
var Error *log.Logger
var Fault *log.Logger

//初始化 所有 全局 日誌 以便調試
func InitDebugLoggers() {
	c := NewCreator()
	out := kio.NewSafeWriter(os.Stdout)

	Trace = c.NewTrace(out)
	Debug = c.NewDebug(out)
	Info = c.NewInfo(out)
	Warn = c.NewWarn(out)
	Error = c.NewError(out)
	Fault = c.NewFault(out)
}
