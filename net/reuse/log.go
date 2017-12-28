package reuse

import (
	"github.com/fatih/color"
	"github.com/zuiwuchang/king-go/strings"
	"log"
	"sync"
)

var LogTrace *log.Logger
var LogDebug *log.Logger
var LogInfo *log.Logger
var LogWarn *log.Logger
var LogError *log.Logger
var LogFault *log.Logger

type logWriter struct {
	sync.Mutex
	Printf func(string, ...interface{})
}

func newLogWriter(f func(string, ...interface{})) *logWriter {
	return &logWriter{
		Printf: f,
	}
}
func (l *logWriter) Write(b []byte) (n int, e error) {
	l.Lock()
	l.Printf(strings.BytesToString(b))
	l.Unlock()
	return len(b), nil
}
func InitDebugLog() {
	logFlags := log.Ltime | log.Lshortfile
	LogTrace = log.New(
		newLogWriter(color.Cyan),
		"[reuse Trace] ",
		logFlags,
	)

	LogDebug = log.New(
		newLogWriter(color.Blue),
		"[reuse Debug] ",
		logFlags,
	)

	LogInfo = log.New(
		newLogWriter(color.Green),
		"[reuse Info] ",
		logFlags,
	)

	LogWarn = log.New(
		newLogWriter(color.Yellow),
		"[reuse Warn] ",
		logFlags,
	)
	LogError = log.New(
		newLogWriter(color.Magenta),
		"[reuse Error] ",
		logFlags,
	)
	LogFault = log.New(
		newLogWriter(color.Red),
		"[reuse Fault] ",
		logFlags,
	)
}
