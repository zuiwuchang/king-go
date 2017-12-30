package log

import (
	kio "github.com/zuiwuchang/king-go/io"
	"io"
	"log"
	"os"
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
