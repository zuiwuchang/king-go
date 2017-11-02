package pool

import (
	"net"
	"time"
)

//包裝的 net.Conn
type Conn struct {
	//連接是否正常
	ok bool
	//tcp連接
	c net.Conn
	//最後工作時間
	lastWork time.Time
	//進入 連接池時間
	lastPut time.Time

	//執行 ping 的 timer
	timer *time.Timer
}

//創建 一個 conn
func newConn(c net.Conn) *Conn {
	now := time.Now()
	return &Conn{
		ok:       true,
		c:        c,
		lastWork: now,
		lastPut:  now,
	}
}

//返回 連接池 模板 Conect 創建的原始 net.Conn
func (c *Conn) Get() net.Conn {
	return c.c
}

//關閉連接 釋放 所有資源
func (c *Conn) free(t IPoolTemplate) {
	t.Close(c.c)
}

//實現 net.Conn接口
func (c *Conn) Read(b []byte) (n int, err error) {
	n, e := c.c.Read(b)
	if e != nil {
		c.lastWork = time.Now()
	}
	return n, e
}

func (c *Conn) Write(b []byte) (n int, err error) {
	n, e := c.c.Write(b)
	if e != nil {
		c.lastWork = time.Now()
	}
	return n, e
}

func (c *Conn) Close() error {
	c.ok = false
	return nil
}
func (c *Conn) LocalAddr() net.Addr {
	return c.c.LocalAddr()
}
func (c *Conn) RemoteAddr() net.Addr {
	return c.c.RemoteAddr()
}
func (c *Conn) SetDeadline(t time.Time) error {
	return c.SetDeadline(t)
}
func (c *Conn) SetReadDeadline(t time.Time) error {
	return c.SetReadDeadline(t)
}
func (c *Conn) SetWriteDeadline(t time.Time) error {
	return c.SetWriteDeadline(t)
}
