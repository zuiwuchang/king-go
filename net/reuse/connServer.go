package reuse

import (
	"bytes"
	"encoding/binary"
	"net"
	"sync"
)

type wrapperConnServer struct {
	//當前的 可用id
	ID uint64

	conn net.Conn
	sync.Mutex
	closed bool

	//通知 所有 復用的 Conn 退出
	chExit chan (interface{})

	//發包 chan
	chWrite chan ([]byte)

	//通知 連接建立 成功
	chRequestAccept  chan (uint64)
	chResponseAccept chan (bool)

	//復用的 Conn
	Conns map[uint64](*connImpl)
}

func newWrapperConn(srv *Server, c net.Conn) *wrapperConnServer {
	wrapper := &wrapperConnServer{
		conn: c,

		chExit:  make(chan (interface{})),
		chWrite: make(chan ([]byte), srv.WriteN),

		chRequestAccept:  make(chan (uint64)),
		chResponseAccept: make(chan (bool)),

		Conns: make(map[uint64](*connImpl)),
	}

	//啟動 發包 goroutine
	go wrapper.writeRoutine()

	return wrapper
}
func (w *wrapperConnServer) Conn() net.Conn {
	return w.conn
}
func (w *wrapperConnServer) CloseId(id uint64, closeRemote bool) {
	if id == 0 {
		return
	}
	w.Lock()
	delete(w.Conns, id)
	w.Unlock()
	if !closeRemote {
		//通知 遠端 復用socket 應該被關閉
		b := make([]byte, HeaderSize+8)
		binary.LittleEndian.PutUint16(b, HeaderSize+8)
		binary.LittleEndian.PutUint16(b[2:], commandClose)
		binary.LittleEndian.PutUint64(b[4:], id)
		w.Write(b)
	}
}
func (w *wrapperConnServer) writeRoutine() {
	//捕獲異常 closed ch<-
	defer func() {
		e := recover()
		if e != nil && LogWarn != nil {
			LogWarn.Println(e)
		}

		if LogTrace != nil {
			LogTrace.Println("goroutine stop [w.Write]")
		}
	}()

	if LogTrace != nil {
		LogTrace.Println("goroutine start [w.Write]")
	}
	var e error
	var b []byte
	var id uint64
	bufferAccept := make([]byte, HeaderSize+8)
	binary.LittleEndian.PutUint16(bufferAccept, HeaderSize+8)
	binary.LittleEndian.PutUint16(bufferAccept[2:], commandAccept)
	for {
		select {
		case <-w.chExit:
			return
		case id = <-w.chRequestAccept:
			//通知 客戶端 連接 成功
			binary.LittleEndian.PutUint64(bufferAccept[HeaderSize:], id)
			e = wrapperWriteBytes(w.conn, bufferAccept)
			if e != nil {
				//socket 錯誤 關閉 連接
				w.Close()

				if LogWarn != nil {
					LogWarn.Println(e)
				}
				return
			}
			w.chResponseAccept <- true
		case b = <-w.chWrite:
			if e = wrapperWriteBytes(w.conn, b); e != nil {
				if LogWarn != nil {
					LogWarn.Println(e)
				}

				//socket 錯誤 關閉 所有連接
				w.Close()
				return
			}
		}
	}
}
func (w *wrapperConnServer) Close() (e error) {
	w.Lock()

	if w.closed {
		w.Unlock()
		e = ErrorConnClosed
		return
	}
	e = w.conn.Close()
	if e != nil && LogError != nil {
		LogError.Println(e)
	}
	close(w.chExit)
	close(w.chWrite)

	close(w.chRequestAccept)
	close(w.chResponseAccept)

	w.closed = true

	w.Unlock()
	return
}

//創建一個 和客戶端連接的 IConn
func (w *wrapperConnServer) Accept(srv *Server) (c IConn) {
	w.Lock()
	if w.ID == 0 {
		w.ID++
	}
	id := w.ID
	w.ID++
	w.Unlock()

	//通知 客戶 id
	defer func() {
		//捕獲異常 closed w.chRequestAccept
		if e := recover(); e != nil {
			if LogWarn != nil {
				LogWarn.Println(e)
			}
		}
	}()
	w.chRequestAccept <- id
	ok := <-w.chResponseAccept
	if !ok { //創建 連接失敗
		return
	}

	if LogInfo != nil {
		LogInfo.Println("one in", id, w.conn.RemoteAddr())
	}
	conn := &connImpl{
		id:      id,
		wrapper: w,

		chExit: make(chan interface{}),
		chRead: make(chan *bytes.Buffer, srv.ReadN),
	}
	w.Lock()
	w.Conns[id] = conn
	w.Unlock()

	go func() {
		if LogTrace != nil {
			LogTrace.Println("goroutine start [w.Accept chExit]")
		}
		//監聽 socket 狀態
		select {
		case <-w.chExit:
			//通知 復用節點 socket 已經斷開
			conn.Lock()
			conn.unsafeClose()
			conn.Unlock()
		case <-conn.chExit:
		}

		if LogTrace != nil {
			LogTrace.Println("goroutine stop [w.Accept chExit]")
		}
	}()
	c = conn
	return
}
func (w *wrapperConnServer) WriteRead(id uint64, b []byte) {
	//捕獲異常 closed impl.chRead
	defer recover()

	w.Lock()
	impl, ok := w.Conns[id]
	w.Unlock()
	if !ok { //id不存在 忽略
		return
	}

	dist := make([]byte, len(b))
	copy(dist, b)
	impl.chRead <- bytes.NewBuffer(dist)
	return
}
func (w *wrapperConnServer) Write(b []byte) (e error) {
	//捕獲異常 closed chWrite
	defer func() {
		if err := recover(); err != nil {
			if LogWarn != nil {
				LogWarn.Println(err)
			}
			e = ErrorConnClosed
		}
	}()

	w.chWrite <- b
	return nil
}
