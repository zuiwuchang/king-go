package echo

import (
	"net"
	"time"
)

//創建一個 echo 服務器
//timeout 客戶端 未活動 斷開時間 如果爲0 不主動斷開
//
//此 api 已經被廢棄 不建議繼續使用 請使用 NewEchoServer
func NewServer(laddr string,
	timeout time.Duration,
	template IServerTemplate,
) (IServer, error) {
	return NewServer2(laddr,
		DefaultBufferLen,
		timeout,
		template,
	)
}

//創建一個 basic 服務器
//timeout 客戶端 未活動 斷開時間 如果爲0 不主動斷開
//
//此 api 已經被廢棄 不建議繼續使用 請使用 NewEchoServer
func NewServer2(laddr string,
	bufLen int,
	timeout time.Duration,
	template IServerTemplate,
) (IServer, error) {
	l, e := net.Listen("tcp", laddr)
	if e != nil {
		return nil, e
	}

	s := newServerImpl(
		&Server{
			Listener:   l,
			Timeout:    timeout,
			RecvBuffer: bufLen,
			Template:   template,
		},
	)

	return s, nil
}
