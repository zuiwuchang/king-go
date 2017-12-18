package reuse

import (
	"net"
	"sync"
)

type listenerImpl struct {
	Mutex    sync.Mutex
	Listener net.Listener

	//標記 服務 是否 運行中
	flag bool

	//執行 accept
	chAccpet chan (interface{})
	//返回 accept
	chRsAccept chan (rsConn)
	//退出 main
	chExit chan (interface{})
}
type rsConn struct {
	Error error
	Conn  net.Conn
}

func newListenerImpl(network, address string) (*listenerImpl, error) {
	l, e := net.Listen(network, address)
	if e != nil {
		return nil, e
	}

	impl := &listenerImpl{
		Listener: l,

		flag:       true,
		chAccpet:   make(chan (interface{}), 1),
		chRsAccept: make(chan (rsConn), 1),
		chExit:     make(chan (interface{}), 1),
	}
	go impl.doMain()
	return impl, nil
}
func (l *listenerImpl) doMain() {
	var c net.Conn
	var e error
	for {
		select {
		case <-l.chExit:
			return
		case <-l.chAccpet:
			if c, e = l.Listener.Accept(); e == nil {
				go l.doRead(c)
			} else {
				l.chRsAccept <- rsConn{
					Error: e,
				}
			}

		}
	}
}
func (l *listenerImpl) doRead(c net.Conn) {
	b := make([]byte, 1024)
	var e error
	for {
		if _, e = c.Read(b); e != nil {
			break
		}

	}
	c.Close()
	/*l.chRsAccept <- rsConn{
		Conn: c,
	}*/
}
func (l *listenerImpl) Accept() (net.Conn, error) {
	l.Mutex.Lock()
	if !l.flag {
		l.Mutex.Unlock()
		return nil, ErrorListenerClosed
	}

	//發送 accept 請求
	l.chAccpet <- nil
	l.Mutex.Unlock()

	//獲取 請求結果
	rs := <-l.chRsAccept
	return rs.Conn, rs.Error
}

func (l *listenerImpl) Close() error {
	l.Mutex.Lock()
	if !l.flag {
		l.Mutex.Unlock()
		return nil
	}
	//關閉 服務器
	l.flag = false
	l.Mutex.Unlock()
	e := l.Listener.Close()

	//通知 main gorutine 退出
	l.chExit <- nil
	return e
}

func (l *listenerImpl) Addr() net.Addr {
	return l.Listener.Addr()
}
