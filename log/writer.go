package log

import (
	"github.com/fatih/color"
	kio "github.com/zuiwuchang/king-go/io"
	"github.com/zuiwuchang/king-go/strings"
	"io"
	"sync"
)

//創建一個 帶顏色 的 stdout 輸出
func NewStdoutColorWriter(color *color.Color, mutex *sync.Mutex) io.Writer {
	return &stdoutColorWriter{
		Color: color,
		Mutex: mutex,
	}
}

type stdoutColorWriter struct {
	Mutex *sync.Mutex
	Color *color.Color
}

func (s *stdoutColorWriter) Write(b []byte) (n int, e error) {
	str := strings.BytesToString(b)
	if s.Mutex == nil {
		n, e = s.Color.Print(str)
	} else {
		s.Mutex.Lock()
		n, e = s.Color.Print(str)
		s.Mutex.Unlock()
	}
	return
}

type colorWriter struct {
	Writer io.Writer
	Color  *color.Color
}

func (w *colorWriter) Write(b []byte) (n int, e error) {
	str := w.Color.Sprint(strings.BytesToString(b))
	b = strings.StringToBytes(str)
	e = kio.WriteAll(w.Writer, b)
	if e != nil {
		return
	}
	n = len(b)
	return
}
func (w *colorWriter) ReadFrom(src io.Reader) (n int64, e error) {
	n, e = io.Copy(w.Writer, src)
	return
}
