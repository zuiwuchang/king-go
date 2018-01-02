package reuse

import (
	kio "github.com/zuiwuchang/king-go/io"
	"net"
)

//包裝的 net.Conn
type wrapperConn struct {
	net.Conn

	Description *Description
	l           *listenImpl

	chExit chan interface{}
}

//返回 一個 包裝的 net.Conn
func newWrapperConn(l *listenImpl, c net.Conn) {
	wc := &wrapperConn{
		c,

		&(l.Server.Description),
		l,

		make(chan interface{}),
	}

	go wc.read()
}
func (w *wrapperConn) read() {
	if kLog.Trace != nil {
		kLog.Trace.Println("start goroutine read", w.RemoteAddr())
	}
	defer func() {
		if kLog.Trace != nil {
			kLog.Trace.Println("stop goroutine read", w.RemoteAddr())
		}
		w.Conn.Close()
		close(w.chExit)
	}()

	//var n int
	var e error
	b := make([]byte, w.Description.RecvBuffer)
	for {
		//讀取 命令指令 頭
		e = kio.ReadAll(w, b[1:])
		if e != nil {
			if kLog.Warn != nil {
				kLog.Warn.Println(e)
			}
			break
		}

		switch b[0] {
		case NetAccept:
		case NetAcceptOk:
		case NetAcceptReject:

		case NetForward:
		case NetAck:

		case NetClose:

		default:
			if kLog.Error != nil {
				kLog.Error.Println("unknow control command", b[0])
			}
			return
		}
	}
}
