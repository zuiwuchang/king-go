//tcp 服務器 客戶端
package basic

import (
	"net"
	"time"
)

const (
	MinBufferLen     = 256
	DefaultBufferLen = 1024 * 8
)
const (
	cmdRun   = 1
	cmdClose = 2
)

//echo 服務器 接口
type IServer interface {
	//運行服務器
	Run()

	//等待服務器 停止
	Wait()

	//返回服務器是否在運行中
	IsRun() bool

	//關閉服務器 並釋放所有資源
	Close()
}

//服務器 定義
type Server struct {
	//監聽的 socket
	Listener net.Listener
	//服務器 模板
	Template IServerTemplate
	//客戶端 未活動 斷開時間 如果爲0 不主動斷開
	Timeout time.Duration
	//recv 緩衝區 大小
	RecvBuffer int
}

//初始化 無效 選項為 默認值
func (s *Server) format() {
	if s.Template == nil {
		s.Template = NewServerTemplate()
	}
	if s.RecvBuffer < MinBufferLen {
		s.RecvBuffer = DefaultBufferLen
	}
}

//創建一個 服務器
func NewBasicServer(srv *Server) IServer {
	return newBasicServer(srv)
}
