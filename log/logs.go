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
	return &Loggers{
		Info:  NewInfo(out, flags),
		Warn:  NewWarn(out, flags),
		Error: NewError(out, flags),
		Fault: NewFault(out, flags),
	}
}

//初始化 默認 調試 日誌
func NewDebugLoggers() *Loggers {
	flags := log.Ltime | log.Lshortfile
	out := kio.NewSafeWriter(os.Stdout)

	return &Loggers{
		Trace: NewTraceColor(out, flags),
		Debug: NewDebugColor(out, flags),
		Info:  NewInfoColor(out, flags),
		Warn:  NewWarnColor(out, flags),
		Error: NewErrorColor(out, flags),
		Fault: NewFaultColor(out, flags),
	}
}
