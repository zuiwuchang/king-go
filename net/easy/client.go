package easy

import (
	"net"
	"time"
)

// 超時 寫 信號
type rSignal struct {
	Message []byte
	Error   error
}

// 超時 寫 信號
type wSignal struct {
	N     int
	Error error
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
	ReadSignal chan rSignal
	// 超時 讀取 定時器
	ReadTimer *time.Timer

	// 帶超時寫入 的 信號
	WriteSignal chan wSignal
	// 超時 寫入 定時器
	WriteTimer *time.Timer
}

// NewClient 將一個 net.Conn 包裝爲 IClient
func NewClient(c net.Conn, bufferSize int, analyze IAnalyze) IClient {
	// 啓動 read goroutine
	// go read
	return &_Client{
		Conn:   c,
		Buffer: make([]byte, bufferSize),

		Analyze: analyze,

		ReadSignal:  make(chan rSignal),
		WriteSignal: make(chan wSignal),
	}
}
func (c *_Client) read() (e error) {
	buffer := c.Buffer[c.BufferPos+c.BufferSize:]
	if len(buffer) == 0 {
		// buufer 已經讀到 尾 將緩存 移到最前端
		copy(c.Buffer, c.Buffer[c.BufferPos:c.BufferPos+c.BufferSize])
		c.BufferPos = 0
		buffer = c.Buffer[c.BufferSize:]
	}

	var n int
	n, e = c.Conn.Read(buffer)
	if e == nil {
		c.BufferSize += n
	}
	return
}
func (c *_Client) readMessage(size int) (msg []byte) {
	msg = make([]byte, size)
	copy(msg, c.Buffer[c.BufferPos:size])
	c.BufferPos += size
	c.BufferSize -= size

	if c.BufferSize == 0 && c.BufferPos != 0 {
		c.BufferPos = 0
	}
	return
}
func (c *_Client) Read() (msg []byte, e error) {
	analyze := c.Analyze
	headerSize := analyze.Header()
	size := -1
	maxSize := len(c.Buffer)
	for {
		// 未知 消息 長度
		if size == -1 {
			// 需要讀取 header
			if c.BufferSize < headerSize {
				if e = c.read(); e != nil {
					break
				}
				continue
			} else {
				// 計算 size
				size, e = analyze.Analyze(c.Buffer[c.BufferPos : c.BufferPos+headerSize])
				if e != nil {
					break
				}
				if size < headerSize {
					e = ErrorMessageSize
					break
				} else if size > maxSize {
					e = ErrorMessageSize
					break
				}
			}
		}

		// 需要 等待 body
		if c.BufferSize < size {
			if e = c.read(); e != nil {
				break
			}
			continue
		}

		// 讀出消息
		msg = c.readMessage(size)
		break
	}
	return
}

func (c *_Client) ReadTimeout(timeout time.Duration) (msg []byte, e error) {
	// 啟動 定時器
	if c.ReadTimer == nil {
		c.ReadTimer = time.NewTimer(timeout)
	} else {
		c.ReadTimer.Reset(timeout)
	}
	// 異步 寫入
	go c.asyncReadF()

	// 等待執行結果
	select {
	case <-c.ReadTimer.C: // 超時
		e = ErrorReadTimeout
	case signal := <-c.ReadSignal:
		// 停止 定時器
		if !c.ReadTimer.Stop() {
			<-c.ReadTimer.C
		}
		msg = signal.Message
		e = signal.Error
	}
	return
}
func (c *_Client) asyncReadF() {
	msg, e := c.Read()
	c.ReadSignal <- rSignal{
		Message: msg,
		Error:   e,
	}
}
func (c *_Client) WaitRead() (msg []byte, e error) {
	signal := <-c.ReadSignal
	msg = signal.Message
	e = signal.Error
	return
}

func (c *_Client) Write(b []byte) (n int, e error) {
	if len(b) == 0 {
		return
	}
	n, e = c.Conn.Write(b)
	return
}
func (c *_Client) asyncWriteF(b []byte) {
	n, e := c.Conn.Write(b)
	c.WriteSignal <- wSignal{
		N:     n,
		Error: e,
	}
}
func (c *_Client) WriteTimeout(b []byte, timeout time.Duration) (n int, e error) {
	if len(b) == 0 {
		return
	}

	// 啟動 定時器
	if c.WriteTimer == nil {
		c.WriteTimer = time.NewTimer(timeout)
	} else {
		c.WriteTimer.Reset(timeout)
	}
	// 異步 寫入
	go c.asyncWriteF(b)

	// 等待執行結果
	select {
	case <-c.WriteTimer.C: // 超時
		e = ErrorWriteTimeout
	case signal := <-c.WriteSignal:
		// 停止 定時器
		if !c.WriteTimer.Stop() {
			<-c.WriteTimer.C
		}
		n = signal.N
		e = signal.Error
	}
	return
}

func (c *_Client) WaitWrite() (n int, e error) {
	signal := <-c.WriteSignal
	n = signal.N
	e = signal.Error
	return
}

func (c *_Client) Close() (e error) {
	e = c.Conn.Close()
	return
}
