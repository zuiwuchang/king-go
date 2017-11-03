package pool

import (
	"net"
	"time"
)

const (
	_StatusIdle   = 0 //空閒
	_StatusGet    = 1 //已經被用戶 使用
	_StatusPing   = 2 //正在執行 ping
	_StatusClose  = 3 //請求關閉socket
	_StatusClosed = 4 //socket已關閉

)

//包裝的 net.Conn
type Conn struct {
	//tcp連接
	net.Conn

	//進入 連接池時間
	lastPut time.Time

	//執行 ping 的 timer
	timer *time.Timer

	//當前狀態
	status int
}

//創建 一個 conn
func newConn(c net.Conn, status int) *Conn {
	now := time.Now()
	return &Conn{
		Conn:    c,
		lastPut: now,
		timer:   nil,
		status:  status,
	}
}

//返回 連接池 模板 Conect 創建的原始 net.Conn
func (c *Conn) Get() net.Conn {
	return c.Conn
}

//關閉連接 釋放 所有資源
func (c *Conn) free(t IPoolTemplate) {
	c.status = _StatusClosed
	t.Close(c.Conn)
}

func (c *Conn) Close() error {
	c.status = _StatusClose
	return nil
}
