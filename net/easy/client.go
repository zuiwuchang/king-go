package easy

import (
	"net"
	"time"
)

// 超時 讀寫 信號
type rwSignal struct {
}

// _Client IClient 的一個實現
type _Client struct {
	// socket
	Conn net.Conn

	// recv 緩衝區
	Buffer []byte
	// buffer 中數據 大小
	BufferSize int
	// buffer 中數據 偏移
	BufferPos int

	// 解包器
	Analyze IAnalyze

	// 帶超時寫入 的 信號
	ReadSignal chan rwSignal

	// 帶超時寫入 的 信號
	WriteSignal chan rwSignal
}

// NewClient 將一個 net.Conn 包裝爲 IClient
func NewClient(c net.Conn, bufferSize int, analyze IAnalyze) IClient {

	// 啓動 read goroutine
	// go read
	return &_Client{
		Conn:   c,
		Buffer: make([]byte, bufferSize),

		Analyze: analyze,

		ReadSignal:  make(chan rwSignal),
		WriteSignal: make(chan rwSignal),
	}
}

func (c *_Client) Read() (msg []byte, e error) {

	return
}

func (c *_Client) ReadTimeout(timeout time.Duration) (msg []byte, e error) {

	return
}
func (c *_Client) WaitRead() (msg []byte, e error) {
	return
}
func (c *_Client) Write(msg []byte) (n int64, e error) {
	return
}

func (c *_Client) WriteTimeout(msg []byte, timeout time.Duration) (n int64, e error) {
	return
}

func (c *_Client) WaitWrite() (n int64, e error) {
	return
}

func (c *_Client) Close() (e error) {
	e = c.Conn.Close()
	return
}
