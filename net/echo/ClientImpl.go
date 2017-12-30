package echo

import (
	kio "github.com/zuiwuchang/king-go/io"
	"io"
	"net"
	"time"
)

//echo 客戶端 實現
type clientImpl struct {
	net.Conn
	template IClientTemplate
	buffer   []byte
}

func (c *clientImpl) GetMessage(timeout time.Duration) (b []byte, e error) {
	//啓動 定時器
	var timer *time.Timer
	if timeout != 0 {
		timer = time.AfterFunc(timeout, func() {
			c.Close()
		})
	}

	template := c.template
	headerSize := template.GetHeaderSize()

	buffer := c.buffer
	header := buffer[:headerSize]
	var msg []byte
	var size int
	//讀取 header
	e = kio.ReadAll(io.LimitReader(c, int64(headerSize)), header)
	if e != nil {
		return
	}
	//獲取 消息長度
	size, e = template.GetMessageSize(header)
	if e != nil || size < headerSize {
		//錯誤的 消息
		return
	}

	//緩衝區不夠
	if len(buffer) < size {
		msg = make([]byte, size)
		copy(msg, buffer[:headerSize])
	} else {
		msg = buffer[:size]
	}
	//讀取 body
	e = kio.ReadAll(c, msg[headerSize:])
	if e != nil {
		//讀取錯誤
		return
	}

	//關閉 定時器
	if timer != nil {
		timer.Stop()
	}
	b = msg

	return
}
