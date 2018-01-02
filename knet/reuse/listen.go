package reuse

import (
	"net"
	"sync"
)

//服務器 定義
type Server struct {
	Listener net.Listener
	//
	AccpeN int
	Description
}

func (s *Server) format() {
	s.Description.Format()

	if s.AccpeN < 1 {
		if kLog.Warn != nil {
			kLog.Warn.Printf("AccpeN(%v) < 1 use default(%v)",
				s.AccpeN,
				DefaultAcceptN,
			)
		}

		s.AccpeN = DefaultAcceptN
	}
	if s.Listener == nil {
		if kLog.Fault != nil {
			kLog.Fault.Println("Listener can't be empty")
		}
		panic("Listener can't be empty")
	}
}

//可復用的 Listener
type listenImpl struct {
	Mutex  sync.Mutex
	Server *Server

	//服務器 退出 通知
	chExit   chan interface{}
	flagExit bool

	//通知 成功 建立連接
	chAccept chan net.Conn
	//已經使用的 accept 緩存
	acceptN int
}

//將 一個 net.Listener 包裝爲 可復用的
func NewListen(srv *Server) net.Listener {
	if srv == nil {
		if kLog.Fault != nil {
			kLog.Fault.Println("Server can't be nil")
		}
		panic("Server can't be nil")
	}

	srv.format()
	s := &Server{}
	*s = *srv

	l := &listenImpl{
		Server: s,

		chExit: make(chan interface{}),

		chAccept: make(chan net.Conn),
		acceptN:  s.AccpeN,
	}
	//啟動 accept goroutine
	go l.acceptRoutine()

	//返回 復用 Listener
	return l
}

func (l *listenImpl) Addr() net.Addr {
	return l.Server.Listener.Addr()
}
func (l *listenImpl) Accept() (net.Conn, error) {
	c := <-l.chAccept
	if c == nil {
		if kLog.Warn != nil {
			kLog.Warn.Println(ErrorListenerClosed)
		}
		return nil, ErrorListenerClosed
	} else {
		//中間 緩存 空閒容量
		l.Mutex.Lock()
		l.acceptN++
		l.Mutex.Unlock()
	}
	return c, nil
}
func (l *listenImpl) Close() (e error) {
	l.Mutex.Lock()
	if l.flagExit { //已經退出
		l.Mutex.Unlock()
		e = ErrorListenerClosed

		if kLog.Warn != nil {
			kLog.Warn.Println(e)
		}
		return e
	}
	e = l.Server.Listener.Close()

	close(l.chExit)

	l.flagExit = true
	l.Mutex.Unlock()
	return nil
}
func (l *listenImpl) acceptRoutine() {
	if kLog.Trace != nil {
		kLog.Trace.Println("start goroutine acceptRoutine")
		defer kLog.Trace.Println("stop goroutine acceptRoutine")
	}
	nl := l.Server.Listener
	var c net.Conn
	var e error
	for !l.flagExit {
		//接收 連接
		c, e = nl.Accept()
		if e != nil {
			if kLog.Error != nil {
				kLog.Error.Println(e)
			}
			l.Mutex.Lock()
			if l.flagExit { //是否需要退出
				l.Mutex.Unlock()
				break
			}
			l.Mutex.Unlock()
			continue
		}

		//創建 復用 Conn
		l.newConn(c)
	}
}
func (l *listenImpl) newConn(c net.Conn) {
	l.Mutex.Lock()

	if l.acceptN < 1 {
		l.Mutex.Unlock()
		if kLog.Warn != nil {
			kLog.Warn.Println("accept buffer is full , reject socket connect .", c.RemoteAddr())
		}
		c.Close()
		return
	}
	newWrapperConn(l, c)
	l.Mutex.Unlock()

	/*//減少 緩存 空閒容量
	l.acceptN--

	//創建 復用 Conn
	nc := newWrapperConn(l, c)
	//通知 Accept 返回
	l.chAccept <- nc*/
}
