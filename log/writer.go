package log

import (
	"github.com/fatih/color"
	kio "github.com/zuiwuchang/king-go/io"
	"github.com/zuiwuchang/king-go/strings"
	"io"
	"os"
)

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

//返回一個 goroutine Safe 的 stdout
func NewStdSafeWriter() io.Writer {
	return kio.NewSafeWriter(os.Stdout)
}
