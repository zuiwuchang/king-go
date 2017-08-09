package echo

import (
	"encoding/binary"
	"errors"
)

//echo 客戶端 模板 定義如何 解析 消息
type IClientTemplate interface {
	//返回 header 長度 必須大於 0
	GetHeaderSize() int
	//傳入 header 返回 整個消息長
	//如果 e != nil 或 消息長 小於 headerSize 將斷開連接
	GetMessageSize(header []byte) (n int, e error)
}

//返回默認模板實現
func NewClientTemplate(order binary.ByteOrder, flag uint16) IClientTemplate {
	return &clientTemplate{order: order,
		headerFlag: flag,
	}
}

//默認的 模板實現
type clientTemplate struct {
	order      binary.ByteOrder
	headerFlag uint16
}

func (c *clientTemplate) GetHeaderSize() int {
	return _HeaderSize
}
func (c *clientTemplate) GetMessageSize(b []byte) (int, error) {
	if len(b) != _HeaderSize {
		return 0, errors.New("header size not match")
	}
	if c.order.Uint16(b) != c.headerFlag {
		return 0, errors.New("header flag not match")
	}
	n := c.order.Uint16(b[2:])
	if n < _HeaderSize {
		return 0, errors.New("message size not match")
	}
	return int(n), nil
}
