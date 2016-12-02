package tcppool

import (
	"net"
	"sync"
	"time"
)

//包裝 net.Conn
type Conn struct {
	Conn net.Conn
	//是否 可用
	IsOk bool

	//最後活動 時間
	lastActive time.Time
	mutex      sync.Mutex
}

//net.Conn Write 如果成功的話 同時 設置 最後活動時間
func (c *Conn) Write(b []byte) (int, error) {
	n, err := c.Conn.Write(b)
	if err == nil {
		c.SetLastActive()
	} else {
		c.IsOk = false
	}
	return n, err
}

//net.Conn Read 如果成功的話 同時 設置 最後活動時間
func (c *Conn) Read(b []byte) (int, error) {
	n, err := c.Conn.Read(b)
	if err == nil {
		c.SetLastActive()
	} else {
		c.IsOk = false
	}

	return n, err
}

//設置 最後活動時間
func (c *Conn) SetLastActive() {
	c.mutex.Lock()
	defer c.mutex.Unlock()
	c.lastActive = time.Now()
}

//返回 最後活動時間
func (c *Conn) GetLastActive() time.Time {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.lastActive
}
