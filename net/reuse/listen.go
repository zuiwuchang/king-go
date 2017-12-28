package reuse

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
)

const (
	HeaderSize = 2 + 2 //  size + cmd

	DefaultAcceptN = 5
	DefaultReadN   = 5
	DefaultWriteN  = 5

	DefaultRecvBuffer = 1024 * 8
)
const (
	//連接成功
	commandAccept = 1
	//read 到數據
	commandRead = 2
	//關閉 復用的 socket
	commandClose = 3
)

//服務器 定義
type Server struct {
	//緩存 多少個 Accept 消息
	AcceptN int
	//每個 復用的tcp 緩存 多少個 Read 消息
	ReadN int
	//每個 復用的tcp 緩存 多少個 Write 消息
	WriteN int

	//recv 緩衝區 大小
	RecvBuffer int
}

//將 無效選項 設置爲 默認值
func (s *Server) Format() {
	if s.AcceptN < 1 {
		s.AcceptN = DefaultAcceptN
	}
	if s.ReadN < 1 {
		s.ReadN = DefaultReadN
	}
	if s.WriteN < 1 {
		s.ReadN = DefaultWriteN
	}

	if s.RecvBuffer < 1024 {
		s.RecvBuffer = DefaultRecvBuffer
	}
}

var ErrorListenerClosed error = errors.New("listener already closed")

type listenerImpl struct {
	Server   *Server
	Listener net.Listener

	Mutex sync.Mutex

	chExit     chan (interface{})
	ChanAccept chan (IConn)

	closed bool
}

//創建 一個 可復用的 默認 Listener
func NewListener(l net.Listener) net.Listener {
	return NewListener2(l, &Server{})
}

//創建 一個 可復用的 定製化 Listener
func NewListener2(l net.Listener, s *Server) net.Listener {
	s.Format()
	rl := &listenerImpl{
		Server:   s,
		Listener: l,

		chExit:     make(chan (interface{})),
		ChanAccept: make(chan (IConn), s.AcceptN),
	}
	go rl.accept()
	return rl
}
func (l *listenerImpl) Addr() net.Addr {
	return l.Listener.Addr()
}
func (l *listenerImpl) Accept() (c net.Conn, e error) {
	c = <-l.ChanAccept
	if c == nil {
		e = ErrorListenerClosed
	}
	return
}
func (l *listenerImpl) Close() (e error) {
	l.Mutex.Lock()
	if l.closed {
		l.Mutex.Unlock()
		e = ErrorListenerClosed
		return
	}
	l.Listener.Close()

	close(l.chExit)
	close(l.ChanAccept)

	l.closed = true
	l.Mutex.Unlock()
	return
}
func (l *listenerImpl) accept() {
	if LogTrace != nil {
		LogTrace.Println("goroutine start [l.accept]")
	}
	for {
		c, e := l.Listener.Accept()
		if e != nil {
			fmt.Println(e)
			break
		}
		go l.read(newWrapperConn(l.Server, c))
	}
	if LogTrace != nil {
		LogTrace.Println("goroutine stop [l.accept]")
	}
}
func (l *listenerImpl) read(w *wrapperConnServer) {
	if LogTrace != nil {
		LogTrace.Println("goroutine start [l.read]")
	}

	go func() {
		if LogTrace != nil {
			LogTrace.Println("goroutine start [l.read chExit]")
		}

		select {
		case <-l.chExit:
		case <-w.chExit:
		}
		w.Close()

		if LogTrace != nil {
			LogTrace.Println("goroutine stop [l.read chExit]")
		}
	}()

	var buffer bytes.Buffer
	b := make([]byte, l.Server.RecvBuffer)
	var e error
	var n, size int
	var buf []byte
	for {
		n, e = w.conn.Read(b)
		if e != nil {
			if LogWarn != nil {
				LogWarn.Println(e)
			}
			break
		}
		e = wrapperWriteBytes(&buffer, b[:n])
		if e != nil {
			if LogWarn != nil {
				LogWarn.Println(e)
			}
			break
		}

		for {
			//讀取 header
			if size == 0 {
				if buffer.Len() < 2 {
					//等待 header
					break
				}
				buf = buffer.Bytes()
				size = int(binary.LittleEndian.Uint16(buf))
				if size < HeaderSize {
					//錯誤的 消息
					if LogWarn != nil {
						LogWarn.Println("bad msg len")
					}
					goto Exit
				}
			}

			//讀取body
			if buffer.Len() < size {
				//等待 body
				break
			}
			buf = make([]byte, size)
			e = wrapperReadBytes(&buffer, buf)
			if e != nil {
				if LogWarn != nil {
					LogWarn.Println(e)
				}
				//讀取錯誤
				goto Exit
			}
			e = l.dealMsg(w, buf)
			if e != nil {
				//消息處理錯誤
				goto Exit
			}

			//重置 消息 解析狀態
			size = 0
		}
	}
Exit:
	w.Close()

	if LogTrace != nil {
		LogTrace.Println("goroutine stop [l.read]")
	}
}
func (l *listenerImpl) dealMsg(w *wrapperConnServer, b []byte) (e error) {
	defer func() {
		//捕獲異常
		if err := recover(); err != nil {
			e = fmt.Errorf("%v", err)
		}
	}()

	cmd := binary.LittleEndian.Uint16(b[2:])
	if LogInfo != nil {
		LogInfo.Println("dealMsg ", cmd)
	}
	switch cmd {
	case commandAccept:
		//執行accept
		l.ChanAccept <- w.Accept(l.Server)
	case commandRead:
		//執行 read
		if len(b) <= HeaderSize+8 {
			if LogError != nil {
				LogError.Println("get commandRead but not data")
			}
			//沒有 需要寫入的 數據
			return
		}
		id := binary.LittleEndian.Uint64(b[4:])
		if id == 0 {
			if LogError != nil {
				LogError.Println("get commandRead but id == 0")
			}
			return
		}
		w.WriteRead(id, b[HeaderSize+8:])
	case commandClose:
		//執行 close
		if len(b) < HeaderSize+8 {
			if LogError != nil {
				LogError.Println("get commandClose but not id")
			}
			//沒有 需要關閉 IConn
			return
		}
		id := binary.LittleEndian.Uint64(b[4:])
		if id == 0 {
			if LogError != nil {
				LogError.Println("get commandClose but id == 0")
			}
			return
		}
		w.Lock()
		conn, ok := w.Conns[id]
		w.Unlock()
		if ok {
			conn.Lock()
			conn.closeRemote = true
			conn.unsafeClose()
			conn.Unlock()
		} else {
			if LogError != nil {
				LogError.Println("get commandClose but", ErrorReuseConnClosed)
			}
			return
		}
	default:
		if LogError != nil {
			LogError.Println("cmd not found", cmd)
		}
	}
	return
}
func wrapperReadBytes(r io.Reader, b []byte) (e error) {
	var n, pos int
	for pos != len(b) {
		n, e = r.Read(b[pos:])
		if e != nil {
			return
		}
		pos += n
	}
	return
}
func wrapperWriteBytes(w io.Writer, b []byte) (e error) {
	var n, pos int
	for pos != len(b) {
		n, e = w.Write(b[pos:])
		if e != nil {
			return e
		}
		pos += n
	}
	return
}
