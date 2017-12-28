package echo

import (
	"bytes"
	"errors"
	kio "github.com/zuiwuchang/king-go/io"
	"net"
	"time"
)

//echo 客戶端 實現
type clientImpl struct {
	net.Conn
	template IClientTemplate
	buffer   *bytes.Buffer
	bufLen   int
	size     int
}

func (c *clientImpl) GetMessage(timeout time.Duration) ([]byte, error) {
	b := make([]byte, c.bufLen)
	var timer *time.Timer
	for {
		if c.buffer.Len() > 0 {
			if b, e := c.getMessage(); e != nil {
				return nil, e
			} else if b != nil {
				return b, nil
			}
		}

		if timeout != 0 {
			timer = time.AfterFunc(timeout, func() {
				c.Close()
			})
		}

		n, e := c.Read(b)
		if timer != nil {
			timer.Stop()
		}

		if e != nil {
			return nil, e
		}

		e = kio.WriteAll(c.buffer, b[:n])
		if e != nil {
			return nil, e
		}

	}
	return nil, nil
}
func (c *clientImpl) getMessage() ([]byte, error) {
	template := c.template
	headerSize := template.GetHeaderSize()
	buffer := c.buffer
	var e error
	for {
		//讀取 header
		if c.size == -1 {
			if buffer.Len() < headerSize {
				//等待 header
				break
			}
			buf := buffer.Bytes()
			c.size, e = template.GetMessageSize(buf[:headerSize])
			if e != nil || c.size < headerSize {
				//錯誤的 消息
				return nil, errors.New("message size not match")
			}

		}

		//讀取body
		if buffer.Len() < c.size {
			//等待 body
			break
		}
		buf := make([]byte, c.size)
		e = kio.ReadAll(buffer, buf)
		if e != nil {
			//讀取錯誤
			return nil, e
		}

		//返回消息
		c.size = -1
		return buf, nil
	}
	return nil, nil
}
