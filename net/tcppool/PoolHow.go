package tcppool

import (
	"net"
	"time"
)

const (
	//連接池最小 連接量
	POOL_MIN_CONN = 10
	//連接池最大 激活量
	POOL_MAX_CONN = 0x00FFFFFF
)

//爲連接池 指定 必要的 行爲
type PoolHow interface {
	//如何 連接 服務器
	Conect() (net.Conn, error)

	//如何 斷開 服務器
	Close(net.Conn) error

	//如何 ping 服務器 以防止 連接 因長時間未使用 而被服務器 關閉
	Ping(*Conn) error

	//返回 自動 ping 間隔
	//如果 返回 0 ping
	PingDuration() time.Duration
	//超過 此時間的 連接 將被 ping
	PingIntervalDuration() time.Duration

	//返回 初始連接服務器時 需要創建的 連接量
	//當 小於 MinConn() 時 使用 MinConn() 的 值
	InitConn() int

	//返回 最小連接 量
	//當 值在 [POOL_MIN_CONN,POOL_MAX_CONN] 之外時 使用 POOL_MIN_CONN
	MinConn() int

	//返回 最大連接 量
	//當 大於 POOL_MAX_CONN 時 使用 POOL_MAX_CONN
	MaxConn() int

	//每次 擴容時 最大擴容量
	//小於 POOL_MIN_CONN 則不限制
	MaxAddStep() int

	//返回是否 需要 擴容 如果 需要 同時 返回 擴容量
	//true,0 需要 擴容 擴容量 自動計算
	ResizeMore(use, free int) (bool, int)
}
