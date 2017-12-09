package echo

import (
	"bytes"
	"errors"
	"net"
	"time"
)

//echo 客戶端 接口
type IClient interface {
	//讀取一個消息
	//timeout 讀取超時(如果超時 自動斷開連接) 爲0 永不超時
	GetMessage(timeout time.Duration) (b []byte, e error)
	net.Conn
}

//echo 客戶端 實現
type client struct {
	net.Conn
	template IClientTemplate
	buffer   *bytes.Buffer
	bufLen   int
	size     int
}

//創建一個 echo 客戶端
func NewClient(addr string, template IClientTemplate) (IClient, error) {
	return NewClient2(addr, DefaultBufferLen, template)
}

//創建一個 echo 客戶端
func NewClient2(addr string, bufLen int, template IClientTemplate) (IClient, error) {
	conn, e := net.Dial("tcp", addr)
	if e != nil {
		return nil, e
	}

	return &client{conn, template, &bytes.Buffer{}, bufLen, -1}, nil
}

func (c *client) GetMessage(timeout time.Duration) ([]byte, error) {
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

		_, e = c.buffer.Write(b[:n])
		if e != nil {
			return nil, e
		}

	}
	return nil, nil
}
func (c *client) getMessage() ([]byte, error) {
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
		_, e = buffer.Read(buf)
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
