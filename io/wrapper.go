package io

import (
	"io"
)

//將 b 全部寫入到 w
func WriteAll(w io.Writer, b []byte) (e error) {
	var n, pos int
	for pos != len(b) {
		n, e = w.Write(b[pos:])
		if e != nil {
			return e
		}
		pos += n
	}
	return
}

//從 r中 讀取數據 直到 b被填滿
func ReadAll(r io.Reader, b []byte) (e error) {
	var n, pos int
	for pos != len(b) {
		n, e = r.Read(b[pos:])
		if e != nil {
			return
		}
		pos += n
	}
	return
}
