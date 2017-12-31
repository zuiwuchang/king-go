package knet

import (
	"net"
)

//客戶端 定義
type Client struct {
	DialFunc func(network, address string) (net.Conn, error)
	Description
}

func (c *Client) format() {
	if c.DialFunc == nil {
		if kLog.Warn != nil {
			kLog.Warn.Println("DialFunc empty use default(net.Dial)")
		}
		c.DialFunc = net.Dial
	}

	c.Description.format()
}

type dialerImpl struct {
	Client *Client

	DialFunc func(network, address string) (net.Conn, error)
}

//將 一個 net.Dialer 包裝爲 IDialer
func NewDialer(client *Client) IDialer {
	c := &Client{}
	if client == nil {
		c.format()
	} else {
		client.format()
		*c = *client
	}

	return &dialerImpl{
		Client: c,
	}
}
func (d *dialerImpl) Dial(network, address string) (IConn, error) {
	c, e := d.Client.DialFunc(network, address)
	if e != nil {
		return nil, e
	}
	return newConnImpl(
			c,
			&(d.Client.Description),
		),
		nil
}
