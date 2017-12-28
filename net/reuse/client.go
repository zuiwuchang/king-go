package reuse

import (
	"net"
	"sync"
)

type DialFunc func(network, address string) (net.Conn, error)

type IDialer interface {
	Dial(network, address string) (net.Conn, error)
}

//客戶端 定義
type Client struct {
	//如何 連接 服務器
	Dial DialFunc

	//每個 復用的tcp 緩存 多少個 Read 消息
	ReadN int
	//每個 復用的tcp 緩存 多少個 Write 消息
	WriteN int

	//recv 緩衝區 大小
	RecvBuffer int
}

//將 無效選項 設置爲 默認值
func (c *Client) Format() {
	if c.Dial == nil {
		c.Dial = net.Dial
	}

	if c.ReadN < 1 {
		c.ReadN = DefaultReadN
	}
	if c.WriteN < 1 {
		c.WriteN = DefaultWriteN
	}

	if c.RecvBuffer < 1024 {
		c.RecvBuffer = DefaultRecvBuffer
	}
}

type dialerImpl struct {
	Client *Client

	sync.Mutex

	//已創建的 連接
	Clients map[string]*wrapperConnClient
}

//創建一個 可復用的 默認 IDialer
func NewDialer() IDialer {
	return NewDialer2(&Client{})
}

//創建一個 可復用的 定製 IDialer
func NewDialer2(c *Client) IDialer {
	c.Format()
	impl := &dialerImpl{
		Client: c,

		Clients: make(map[string]*wrapperConnClient),
	}
	return impl
}

//創建一個 復用 Conn
func (d *dialerImpl) Dial(network, address string) (c net.Conn, e error) {
	d.Lock()
	defer d.Unlock()

	key := network + address
	client, ok := d.Clients[key]
	if !ok {
		//連接 不存在 創建新連接
		var conn net.Conn
		conn, e = d.Client.Dial(network, address)
		if e != nil {
			if LogWarn != nil {
				LogWarn.Println(e)
			}
			return
		}
		client = newWrapperConnClient(d.Client, conn)
		//監聽 socket 關閉
		go func() {
			if LogTrace != nil {
				LogTrace.Println("goroutine start [Dial wait]")
			}

			<-client.chExit
			d.Lock()
			delete(d.Clients, key)
			d.Unlock()

			if LogTrace != nil {
				LogTrace.Println("goroutine stop [Dial wait]")
			}
		}()

		//保存 client
		d.Clients[key] = client
	}
	//創建 復用 Conn
	c, e = client.Dial(d.Client)
	if e != nil {
		if LogWarn != nil {
			LogWarn.Println(e)
		}
	}
	return
}
