package echo

import (
	"bytes"
	"encoding/binary"
	"net"
	"time"
)

type DialFunc func(network, address string) (net.Conn, error)

type IDialer interface {
	Dial(network, address string) (net.Conn, error)
}

//echo 客戶端 接口
type IClient interface {
	//讀取一個消息
	//timeout 讀取超時(如果超時 自動斷開連接) 爲0 永不超時
	GetMessage(timeout time.Duration) (b []byte, e error)
	net.Conn
}

//客戶端 定義
type Client struct {
	Network string
	Addr    string

	//如何 連接 服務器
	Dialer IDialer

	//客戶端 模板
	Template IClientTemplate

	//recv 緩衝區 大小
	RecvBuffer int
}

//初始化 無效 選項為 默認值
func (c *Client) format() {
	if c.Network == "" {
		c.Network = "tcp"
	}
	if c.Dialer == nil {
		c.Dialer = g_netDialer
	}
	if c.Template == nil {
		c.Template = NewClientTemplate(binary.LittleEndian, 731)
	}
	if c.RecvBuffer < MinBufferLen {
		c.RecvBuffer = DefaultBufferLen
	}
}

//創建一個 客戶端
func NewEchoClient(client *Client) (IClient, error) {
	client.format()
	conn, e := client.Dialer.Dial(client.Network, client.Addr)
	if e != nil {
		return nil, e
	}
	return &clientImpl{conn, client.Template, &bytes.Buffer{}, client.RecvBuffer, -1}, nil
}

//定義 默認的 IDialer
var g_netDialer netDialer

type netDialer struct {
}

func (netDialer) Dial(network, address string) (net.Conn, error) {
	return net.Dial(network, address)
}
