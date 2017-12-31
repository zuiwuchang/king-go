package io

import (
	"fmt"
	"io"
)

//將 b 全部寫入到 w
func WriteAll(w io.Writer, b []byte) (e error) {
	if len(b) == 0 {
		return
	}

	var n, pos int
	for pos != len(b) {
		n, e = w.Write(b[pos:])
		if n > 0 {
			pos += n
		}
		if e != nil && e != io.ErrShortWrite {
			break
		}
	}
	return
}

//從 r中 讀取數據 直到 b被填滿
func ReadAll(r io.Reader, b []byte) (e error) {
	if len(b) == 0 {
		return
	}

	var n, pos int
	for pos != len(b) {
		n, e = r.Read(b[pos:])
		if n > 0 {
			pos += n
		}
		if e != nil {
			break
		}
	}
	return
}

//從 r中 讀取數據 直到 b被填滿
func ReadAllEx(r io.Reader, b []byte) (pos int, e error) {
	if len(b) == 0 {
		return
	}

	var n int
	for pos != len(b) {
		n, e = r.Read(b[pos:])
		if n > 0 {
			pos += n
		}
		if e != nil {
			break
		}
	}
	return
}

//將 b 全部寫入到 w
func WriteAllEx(w io.Writer, b []byte) (pos int, e error) {
	if len(b) == 0 {
		return
	}

	var n int
	for pos != len(b) {
		n, e = w.Write(b[pos:])
		fmt.Println(n, e)
		if n > 0 {
			pos += n
		}
		if e != nil && e != io.ErrShortWrite {
			break
		}
	}
	return
}
