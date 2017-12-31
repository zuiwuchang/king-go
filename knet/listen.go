package knet

import (
	"net"
)

//服務器 定義
type Server struct {
	Listener net.Listener
	Description
}

func (s *Server) format() {
	if s.Listener == nil {
		if kLog.Fault != nil {
			kLog.Fault.Println("Listener can't be empty")
		}
		panic("Listener can't be empty")
	}

	s.Description.Format()
}

type listenImpl struct {
	Server *Server
}

//將 一個 net.Listener 包裝爲 ILister
func NewListen(srv *Server) ILister {
	if srv == nil {
		if kLog.Fault != nil {
			kLog.Fault.Println("Server can't be nil")
		}
		panic("Server can't be nil")
	}

	srv.format()
	s := &Server{}
	*s = *srv

	return &listenImpl{
		Server: s,
	}
}
func (l *listenImpl) Addr() net.Addr {
	return l.Server.Listener.Addr()
}
func (l *listenImpl) Accept() (IConn, error) {
	c, e := l.Server.Listener.Accept()
	if e != nil {
		if kLog.Warn != nil {
			kLog.Warn.Println(e)
		}
		return nil, e
	}
	return NewConn(
			c,
			&(l.Server.Description),
		),
		nil

}
func (l *listenImpl) Close() (e error) {
	e = l.Server.Listener.Close()
	if e != nil && kLog.Warn != nil {
		kLog.Warn.Println(e)
	}
	return e
}
