package knet

import (
	kio "github.com/zuiwuchang/king-go/io"
	"io"
	"net"
)

type connImpl struct {
	net.Conn
	Description *Description
	//recv 緩衝區
	Buffer []byte
	//緩存的 header
	Header []byte
	//緩存 消息 長度
	Size int
}

//將一個 net.Conn 包裝爲 IConn
func NewConn(c net.Conn, desc *Description) IConn {
	return &connImpl{
		c,
		desc,
		make([]byte, desc.RecvBuffer),
		nil,
		0,
	}
}

//讀取一條完整的 消息
//
//如果 返回 任何 錯誤 則 socket 已經不可靠 應該關閉之
func (c *connImpl) ReadMessage() (b []byte, e error) {
	size := c.Size
	if size == 0 {
		//讀取 Header
		var header []byte
		header, e = c.ReadHeader()
		if e != nil {
			return
		}
		size, e = c.Description.Analyze.Analyze(header)
		if e != nil {

			if kLog.Error != nil {
				kLog.Error.Println(e, size)
			}
			return
		} else if size < c.Description.headerSize {
			e = ErrorMessageSize

			if kLog.Error != nil {
				kLog.Error.Println(e)
			}
			return
		}
	}
	//讀取 body
	if size == c.Description.headerSize {
		b = c.Header
	} else {
		var msg []byte
		if len(c.Buffer) < size {
			//緩衝區不夠 創建新緩衝區
			msg = make([]byte, size)
			copy(msg, c.Header)
		} else {
			msg = c.Buffer[:size]
		}
		//讀取 body
		e = kio.ReadAll(
			io.LimitReader(
				c.Conn,
				int64(size-c.Description.headerSize),
			),
			msg[c.Description.headerSize:],
		)

		if e != nil {

			if kLog.Warn != nil {
				kLog.Warn.Println(e)
			}
			return
		}
		b = msg
	}

	//重置 消息讀取 狀態
	c.Header = nil
	c.Size = 0
	return
}

//讀取一條消息的 header
//
//多次 調用 ReadHeader 不會將 message 移除緩衝區 你只會得到同一個消息的 header
//
//你必須 調用 ReadMessage 才能將 message 移除緩衝區 之後繼續讀取 後續的 緩衝區
func (c *connImpl) ReadHeader() (b []byte, e error) {
	if c.Header != nil {
		b = c.Header
		return
	}

	buf := c.Buffer[:c.Description.headerSize]
	e = kio.ReadAll(io.LimitReader(c.Conn, int64(c.Description.headerSize)), buf)
	if e != nil {
		if e != io.EOF {
			if kLog.Warn != nil {
				kLog.Warn.Println(e)
			}
		}
		return
	}
	c.Header = buf
	b = buf
	return
}
