package echo

import (
	"bytes"
	"net"
)

//創建一個 echo 客戶端
//
//此 api 已經被廢棄 不建議繼續使用 請使用 NewEchoClient
func NewClient(addr string, template IClientTemplate) (IClient, error) {
	return NewClient2(addr, DefaultBufferLen, template)
}

//創建一個 echo 客戶端
//
//此 api 已經被廢棄 不建議繼續使用 請使用 NewEchoClient
func NewClient2(addr string, bufLen int, template IClientTemplate) (IClient, error) {
	conn, e := net.Dial("tcp", addr)
	if e != nil {
		return nil, e
	}

	return &clientImpl{conn, template, &bytes.Buffer{}, bufLen, -1}, nil
}
