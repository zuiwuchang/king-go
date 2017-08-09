package echo

import (
	"encoding/binary"
	"errors"
	"net"
)

type Session interface{}

const (
	_HeaderSize = 4
)

//echo 服務器 模板 定義如何 解析 消息
type IServerTemplate interface {
	//返回 header 長度 必須大於 0
	GetHeaderSize() int
	//傳入 header 返回 整個消息長
	//如果 e != nil 或 消息長 小於 headerSize 將斷開連接
	GetMessageSize(header []byte) (n int, e error)

	//如何 創建 session 連接成功時 自動回調
	//如果 e != nil 將斷開連接
	NewSession(c net.Conn) (session Session, e error)

	//如何 銷毀 session 連接斷開時 自動 回調
	DeleteSession(c net.Conn, session Session)

	//如何 響應 消息
	//如果 e != nil 將斷開連接
	Message(c net.Conn, session Session, b []byte) (e error)
}

//返回默認模板實現
func NewServerTemplate(order binary.ByteOrder, flag uint16) IServerTemplate {
	return &serverTemplate{order: order,
		headerFlag: flag,
	}
}

//默認的 模板實現
type serverTemplate struct {
	order      binary.ByteOrder
	headerFlag uint16
}

func (s *serverTemplate) GetHeaderSize() int {
	return _HeaderSize
}
func (s *serverTemplate) GetMessageSize(b []byte) (int, error) {
	if len(b) != _HeaderSize {
		return 0, errors.New("header size not match")
	}
	if s.order.Uint16(b) != s.headerFlag {
		return 0, errors.New("header flag not match")
	}
	n := s.order.Uint16(b[2:])
	if n < _HeaderSize {
		return 0, errors.New("message size not match")
	}
	return int(n), nil
}
func (s *serverTemplate) NewSession(c net.Conn) (session Session, e error) {
	return nil, nil
}
func (s *serverTemplate) DeleteSession(c net.Conn, session Session) {
}
func (s *serverTemplate) Message(c net.Conn, session Session, b []byte) error {
	_, e := c.Write(b)
	return e
}
