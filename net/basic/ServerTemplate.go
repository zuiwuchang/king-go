package basic

import (
	"encoding/binary"
	"net"
)

type Session interface{}

//服務器 模板 定義如何 創建 session 處理接受到的數據
type IServerTemplate interface {
	//如何 創建 session 連接成功時 自動回調
	//如果 e != nil 將斷開連接
	NewSession(c net.Conn) (session Session, e error)

	//如何 銷毀 session 連接斷開時 自動 回調
	DeleteSession(c net.Conn, session Session)

	//如何 響應 接收到的 數據
	//如果 e != nil 將斷開連接
	Message(c net.Conn, session Session, b []byte) (e error)
}

//返回默認模板實現
func NewServerTemplate(order binary.ByteOrder) IServerTemplate {
	return &serverTemplate{order: order}
}

//默認的 模板實現
type serverTemplate struct {
	order binary.ByteOrder
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
