package log

import (
	"github.com/fatih/color"
	kio "github.com/zuiwuchang/king-go/io"
	"io"
	"log"
	"os"
)

var Trace *log.Logger
var Debug *log.Logger
var Info *log.Logger
var Warn *log.Logger
var Error *log.Logger
var Fault *log.Logger

func NewTrace(out io.Writer, flags int) *log.Logger {
	return log.New(
		out,
		"[Trace] ",
		flags,
	)
}
func NewDebug(out io.Writer, flags int) *log.Logger {
	return log.New(
		out,
		"[Debug] ",
		flags,
	)
}
func NewInfo(out io.Writer, flags int) *log.Logger {
	return log.New(
		out,
		"[Info] ",
		flags,
	)
}
func NewWarn(out io.Writer, flags int) *log.Logger {
	return log.New(
		out,
		"[Warn] ",
		flags,
	)
}
func NewError(out io.Writer, flags int) *log.Logger {
	return log.New(
		out,
		"[Error] ",
		flags,
	)
}
func NewFault(out io.Writer, flags int) *log.Logger {
	return log.New(
		out,
		"[Fault] ",
		flags,
	)
}

func NewTraceColor(out io.Writer, flags int) *log.Logger {
	return log.New(
		&colorWriter{
			Color:  color.New(color.FgCyan),
			Writer: out,
		},
		"[Trace] ",
		flags,
	)
}
func NewDebugColor(out io.Writer, flags int) *log.Logger {
	return log.New(
		&colorWriter{
			Color:  color.New(color.FgBlue),
			Writer: out,
		},
		"[Debug] ",
		flags,
	)
}
func NewInfoColor(out io.Writer, flags int) *log.Logger {
	return log.New(
		&colorWriter{
			Color:  color.New(color.FgGreen),
			Writer: out,
		},
		"[Info] ",
		flags,
	)
}
func NewWarnColor(out io.Writer, flags int) *log.Logger {
	return log.New(
		&colorWriter{
			Color:  color.New(color.FgYellow),
			Writer: out,
		},
		"[Warn] ",
		flags,
	)
}
func NewErrorColor(out io.Writer, flags int) *log.Logger {
	return log.New(
		&colorWriter{
			Color:  color.New(color.FgMagenta),
			Writer: out,
		},
		"[Error] ",
		flags,
	)
}
func NewFaultColor(out io.Writer, flags int) *log.Logger {
	return log.New(
		&colorWriter{
			Color:  color.New(color.FgRed),
			Writer: out,
		},
		"[Fault] ",
		flags,
	)
}

//初始化 所有 全局 日誌 以便調試
func InitDebugLoggers() {
	flags := log.Ltime | log.Lshortfile
	out := kio.NewSafeWriter(os.Stdout)

	Trace = NewTraceColor(out, flags)
	Debug = NewDebugColor(out, flags)
	Info = NewInfoColor(out, flags)
	Warn = NewWarnColor(out, flags)
	Error = NewErrorColor(out, flags)
	Fault = NewFaultColor(out, flags)
}
