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
	//如果 返回 0 將不會執行ping
	PingInterval() time.Duration

	//返回 連接池 最少連接數
	//初始化時 會創建 min個 連接 (如果 min 大於0的話)
	MinConn() int

	//多久 未活動的 連接 會被 釋放
	Timeout() time.Duration

	//連接池 執行釋放 操作 的最小 執行週期
	MinFreeInterval() time.Duration
	//連接池 多久執行一次自動 縮容
	FreeInterval() time.Duration
}
