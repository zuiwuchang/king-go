package log

import (
	"github.com/fatih/color"
	kio "github.com/zuiwuchang/king-go/io"
	"io"
	"log"
	"os"
	"sync"
)

type Loggers struct {
	Trace, Debug, Info, Warn, Error, Fault *log.Logger
}

//創建默認的 日誌
func NewLoggers(out io.Writer, flags int) *Loggers {
	c := NewCreator()
	c.Flags = flags
	return &Loggers{
		Info:  c.NewInfo(out),
		Warn:  c.NewWarn(out),
		Error: c.NewError(out),
		Fault: c.NewFault(out),
	}
}

//初始化 默認 調試 日誌
func NewDebugLoggers() *Loggers {
	c := NewCreator()
	c.Color = false
	m := &sync.Mutex{}

	return &Loggers{
		Trace: c.NewTrace(NewStdoutColorWriter(color.New(color.FgCyan), m)),
		Debug: c.NewDebug(NewStdoutColorWriter(color.New(color.FgBlue), m)),
		Info:  c.NewInfo(NewStdoutColorWriter(color.New(color.FgGreen), m)),
		Warn:  c.NewWarn(NewStdoutColorWriter(color.New(color.FgYellow), m)),
		Error: c.NewError(NewStdoutColorWriter(color.New(color.FgMagenta), m)),
		Fault: c.NewFault(NewStdoutColorWriter(color.New(color.FgRed), m)),
	}
}

//初始化 默認 調試 日誌
func NewDebugLoggers2(tag string) *Loggers {
	c := NewCreator()
	c.Tag = tag
	out := kio.NewSafeWriter(os.Stdout)

	return &Loggers{
		Trace: c.NewTrace(out),
		Debug: c.NewDebug(out),
		Info:  c.NewInfo(out),
		Warn:  c.NewWarn(out),
		Error: c.NewError(out),
		Fault: c.NewFault(out),
	}
}
