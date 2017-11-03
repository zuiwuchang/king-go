package pool

import (
	"net"
	"time"
)

//連接池 模板 定義了 連接池 工作細節
type IPoolTemplate interface {
	//連接池 將使用此 函數 連接 服務器
	Conect() (net.Conn, error)
	//連接池 將使用此 函數 斷開 連接
	Close(net.Conn) error
	//連接池 將使用此 函數 ping 服務器
	Ping(*Conn) error

	//連接池 中的連接 多久自動執行一次 ping
	//返回 0 將不會執行ping
	//返回值 必須 恆定
	PingInterval() time.Duration

	//多久 未活動的 連接 會被 釋放
	//返回 0 永不超時
	//返回值 必須 恆定
	Timeout() time.Duration

	//連接池 執行縮擴容 操作 的最小 執行週期
	MinInterval() time.Duration
	//連接池 多久執行一次自動 縮擴容
	//返回 0 不執行
	//返回值 必須 恆定
	Interval() time.Duration

	//是否需要 縮擴容
	//返回 0 不需要擴容
	//返回 >0 擴容 |n|
	//返回 <0 縮容 |n|
	Resize(use /*已被get使用的 連接數*/, free /*未被使用的 連接數*/ int) int
}

//創建一個 默認的 連接池 模板
func NewPoolTemplate(addr string, //服務器地址
	timeout time.Duration, //未活動超時時間
	min int, //最少 空閒 連接數
	max int, //最多 空閒 連接數
) *PoolTemplate {
	if min < 1 {
		min = 1
	}
	if max < 1 {
		max = 10
	}
	if min < max {
		max, min = min, max
	}
	return &PoolTemplate{
		addr:    addr,
		timeout: timeout,
		min:     min,
		max:     max,
	}
}

//默認的 連接池 模板
type PoolTemplate struct {
	//服務器 地址
	addr string
	//未活動超時時間
	timeout time.Duration
	//最少 空閒 連接數
	min int
	//最多 空閒 連接數
	max int
}

func (p *PoolTemplate) Conect() (net.Conn, error) {
	return net.Dial("tcp", p.addr)
}
func (p *PoolTemplate) Close(c net.Conn) error {
	return c.Close()
}
func (p *PoolTemplate) Ping(*Conn) error {
	return nil
}
func (p *PoolTemplate) PingInterval() time.Duration {
	return 0
}
func (p *PoolTemplate) Timeout() time.Duration {
	return p.timeout
}
func (p *PoolTemplate) MinInterval() time.Duration {
	return time.Second * 5
}
func (p *PoolTemplate) Interval() time.Duration {
	return time.Minute * 5
}
func (p *PoolTemplate) Resize(use, free int) int {
	if free < p.min {
		return free - (p.min+p.max)/2
	} else if free > p.max {
		return free - p.max
	}
	return 0
}
