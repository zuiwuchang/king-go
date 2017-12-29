package io

import (
	"io"
	"sync"
)

type safeWriter struct {
	sync.Mutex
	Writer io.Writer
}

func NewSafeWriter(w io.Writer) io.Writer {
	return &safeWriter{
		Writer: w,
	}
}
func (w *safeWriter) Write(b []byte) (n int, e error) {
	w.Lock()
	n, e = w.Writer.Write(b)
	w.Unlock()
	return
}
