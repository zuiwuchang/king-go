package reuse

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"sync"
)

var bufferRequestNew []byte

func init() {
	bufferRequestNew = make([]byte, HeaderSize)
	binary.LittleEndian.PutUint16(bufferRequestNew, HeaderSize)
	binary.LittleEndian.PutUint16(bufferRequestNew[2:], commandAccept)
}

type wrapperConnClient struct {
	conn net.Conn
	sync.Mutex
	closed bool

	//通知 所有 復用的 Conn 退出
	chExit chan (interface{})

	//發包 chan
	chWrite chan ([]byte)

	//請求創建 新的 Conn
	chRequestAccept  chan (bool)
	chResponseAccept chan (uint64)

	//復用的 Conn
	Conns map[uint64](*connImpl)
}

func newWrapperConnClient(client *Client, c net.Conn) *wrapperConnClient {
	wrapper := &wrapperConnClient{
		conn: c,

		chExit:  make(chan (interface{})),
		chWrite: make(chan ([]byte), client.WriteN),

		chRequestAccept:  make(chan (bool)),
		chResponseAccept: make(chan (uint64)),

		Conns: make(map[uint64](*connImpl)),
	}

	//啟動 發包 goroutine
	go wrapper.writeRoutine()
	//啟動 收包 goroutine
	go wrapper.readRoutine(client)

	return wrapper
}
func (w *wrapperConnClient) Conn() net.Conn {
	return w.conn
}
func (w *wrapperConnClient) CloseId(id uint64, closeRemote bool) {
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
func (w *wrapperConnClient) readRoutine(client *Client) {
	if LogTrace != nil {
		LogTrace.Println("goroutine start [c.read]")
	}

	var buffer bytes.Buffer
	b := make([]byte, client.RecvBuffer)
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
			e = w.dealMsg(buf)
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
		LogTrace.Println("goroutine stop [c.read]")
	}
}
func (w *wrapperConnClient) dealMsg(b []byte) (e error) {
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
		b = b[HeaderSize:]
		if len(b) < 8 {
			if LogError != nil {
				LogError.Println("recv commandAccept but not set id")
			}
			return
		}
		id := binary.LittleEndian.Uint64(b)
		if id == 0 {
			LogError.Println("recv commandAccept but id = 0")
		}
		w.chResponseAccept <- id
	case commandRead:
		//執行 read
		if len(b) <= HeaderSize+8 {
			if LogError != nil {
				LogError.Println("get commandRead bug not data")
			}
			//沒有 需要寫入的 數據
			return
		}
		id := binary.LittleEndian.Uint64(b[4:])
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
func (w *wrapperConnClient) writeRoutine() {
	if LogTrace != nil {
		LogTrace.Println("goroutine start [c w.Write]")
	}
	var e error
	var b []byte
	for {
		select {
		case <-w.chExit:
			goto END
		case <-w.chRequestAccept:
			if e = wrapperWriteBytes(w.conn, bufferRequestNew); e != nil {
				if LogWarn != nil {
					LogWarn.Println(e)
				}

				//socket 錯誤 關閉 所有連接
				w.Close()
				goto END
			}
		case b = <-w.chWrite:
			if e = wrapperWriteBytes(w.conn, b); e != nil {
				if LogWarn != nil {
					LogWarn.Println(e)
				}

				//socket 錯誤 關閉 所有連接
				w.Close()
				goto END
			}
		}
	}
END:
	if LogTrace != nil {
		LogTrace.Println("goroutine stop [c w.Write]")
	}
}
func (w *wrapperConnClient) Close() (e error) {
	w.Lock()
	if w.closed {
		w.Unlock()
		e = ErrorConnClosed
		return
	}

	e = w.conn.Close()
	close(w.chExit)
	close(w.chWrite)

	close(w.chRequestAccept)
	close(w.chResponseAccept)

	w.closed = true
	w.Unlock()
	return
}

//創建一個 和 服務器 連接的 net.Conn
func (w *wrapperConnClient) Dial(srv *Client) (c IConn, e error) {
	//捕獲異常 closed impl.chRequestNew
	defer func() {
		if err := recover(); err != nil {
			e = ErrorConnClosed
		}
	}()

	w.chRequestAccept <- true
	id := <-w.chResponseAccept
	if id == 0 {
		e = ErrorConnClosed
		return
	}
	//創建 Conn
	conn := &connImpl{
		id:      id,
		wrapper: w,

		chExit: make(chan interface{}),
		chRead: make(chan (*bytes.Buffer), srv.ReadN),
	}
	w.Lock()
	w.Conns[id] = conn
	w.Unlock()

	if LogInfo != nil {
		LogInfo.Println("New Client", id, w.conn.LocalAddr(), w.conn.RemoteAddr())
	}

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
func (w *wrapperConnClient) WriteRead(id uint64, b []byte) {
	//捕獲異常 closed impl.chRead
	defer recover()

	w.Lock()
	impl, ok := w.Conns[id]
	w.Unlock()
	if !ok { //id不存在 忽略
		if LogWarn != nil {
			LogWarn.Println("recv but id not found", id)
		}
		return
	}

	dist := make([]byte, len(b))
	copy(dist, b)
	impl.chRead <- bytes.NewBuffer(dist)
	return
}
func (w *wrapperConnClient) Write(b []byte) (e error) {
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
