package reuse

import (
	"bytes"
	"encoding/binary"
	"errors"
	"net"
	"sync"
	"time"
)

type IConn interface {
	net.Conn
	//返回 復用唯一標識
	Id() uint64
}
type iWrapperConn interface {
	Conn() net.Conn
	Write([]byte) error
	CloseId(uint64, bool)
}

var ErrorConnClosed error = errors.New("Conn already closed")
var ErrorReuseConnClosed error = errors.New("Reuse Conn already closed")

//可復用的 net.Conn
type connImpl struct {
	sync.Mutex
	//復用的 唯一 標識
	id uint64

	wrapper iWrapperConn

	//退出標記
	chExit chan (interface{})

	//recv 緩衝區
	chRead chan (*bytes.Buffer)

	//recv 數據緩存
	readBuffer *bytes.Buffer

	//關閉 標記
	closed       bool
	closedSocket bool

	//遠端是否 關閉
	closeRemote bool
}

func (c *connImpl) Id() uint64 {
	return c.id
}
func (c *connImpl) unsafeClose() {
	if c.closed {
		return
	}

	close(c.chRead)
	close(c.chExit)

	c.closed = true

	c.wrapper.CloseId(c.id, c.closeRemote)
	c.closeRemote = true

	if LogInfo != nil {
		LogInfo.Println("one out", c.Id(), c.RemoteAddr())
	}
}
func (c *connImpl) Close() (e error) {
	c.Lock()
	if c.closedSocket {
		c.Unlock()
		e = ErrorReuseConnClosed
		return
	}
	c.unsafeClose()
	c.closedSocket = true
	c.Unlock()
	return
}
func (c *connImpl) LocalAddr() net.Addr {
	return c.wrapper.Conn().LocalAddr()
}
func (c *connImpl) RemoteAddr() net.Addr {
	return c.wrapper.Conn().RemoteAddr()
}
func (c *connImpl) SetDeadline(t time.Time) error {
	return c.wrapper.Conn().SetDeadline(t)
}
func (c *connImpl) SetReadDeadline(t time.Time) error {
	return c.wrapper.Conn().SetReadDeadline(t)
}
func (c *connImpl) SetWriteDeadline(t time.Time) error {
	return c.wrapper.Conn().SetWriteDeadline(t)
}
func (c *connImpl) Read(b []byte) (n int, e error) {
	if c.readBuffer == nil || c.readBuffer.Len() == 0 {
		//var n int
		c.readBuffer = <-c.chRead
		if c.readBuffer == nil { //已經關閉
			e = ErrorReuseConnClosed
			return
		}
	}
	n, e = c.readBuffer.Read(b)
	return
}
func (c *connImpl) Write(b []byte) (n int, e error) {
	if len(b) == 0 {
		return
	}
	if c.closed { //已經關閉
		e = ErrorReuseConnClosed
		return
	}

	size := len(b) + HeaderSize + 8
	buf := make([]byte, size)
	binary.LittleEndian.PutUint16(buf, uint16(size))
	binary.LittleEndian.PutUint16(buf[2:], commandRead)
	binary.LittleEndian.PutUint64(buf[4:], c.id)
	copy(buf[HeaderSize+8:], b)

	e = c.wrapper.Write(buf)
	if e != nil {
		return
	}
	n = len(b)
	return
}
