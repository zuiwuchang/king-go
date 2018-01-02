package log

import (
	"github.com/fatih/color"
	"log"
	"sync"
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
	c.Color = false
	m := &sync.Mutex{}
	Trace = c.NewTrace(NewStdoutColorWriter(color.New(color.FgCyan), m))
	Debug = c.NewDebug(NewStdoutColorWriter(color.New(color.FgBlue), m))
	Info = c.NewInfo(NewStdoutColorWriter(color.New(color.FgGreen), m))
	Warn = c.NewWarn(NewStdoutColorWriter(color.New(color.FgYellow), m))
	Error = c.NewError(NewStdoutColorWriter(color.New(color.FgMagenta), m))
	Fault = c.NewFault(NewStdoutColorWriter(color.New(color.FgRed), m))
}
