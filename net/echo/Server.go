//tcp 實現的 echo 服務器 客戶端
package echo

import (
	"bytes"
	"net"
	"time"
)

const (
	MinBufferLen     = 256
	DefaultBufferLen = 1024 * 4
)
const (
	cmdRun   = 1
	cmdClose = 2
)

//echo 服務器 接口
type IServer interface {
	//運行服務器
	Run()

	//等待服務器 停止
	Wait()

	//返回服務器是否在運行中
	IsRun() bool

	//關閉服務器 並釋放所有資源
	Close()
}

//echo 服務器 實現
type server struct {
	template IServerTemplate

	listener net.Listener
	timeout  time.Duration

	//recv 緩衝區 大小
	bufLen int

	//服務器 是否運行中
	run bool

	//控制指令
	cmd chan int

	signalWait chan int

	clients   map[net.Conn]int
	signalIn  chan net.Conn
	signalOut chan net.Conn
}

//創建一個 echo 服務器
//timeout 客戶端 未活動 斷開時間 如果爲0 不主動斷開
func NewServer(laddr string,
	timeout time.Duration,
	template IServerTemplate,
) (IServer, error) {
	return NewServer2(laddr,
		DefaultBufferLen,
		timeout,
		template,
	)
}

//創建一個 basic 服務器
//timeout 客戶端 未活動 斷開時間 如果爲0 不主動斷開
func NewServer2(laddr string,
	bufLen int,
	timeout time.Duration,
	template IServerTemplate,
) (IServer, error) {
	if bufLen < MinBufferLen {
		bufLen = MinBufferLen
	}
	l, e := net.Listen("tcp", laddr)
	if e != nil {
		return nil, e
	}

	s := &server{
		listener: l,

		bufLen: bufLen,

		timeout:    timeout,
		cmd:        make(chan int),
		signalWait: make(chan int),
		template:   template,

		clients:   make(map[net.Conn]int),
		signalIn:  make(chan net.Conn),
		signalOut: make(chan net.Conn),
	}

	go s.runController()
	return s, nil
}
func (s *server) runController() {
	defer func() {
		if s.listener != nil {
			s.listener.Close()
		}
	}()
	for {
		select {
		case cmd := <-s.cmd:
			if s.runCommand(cmd) {
				return
			}
		case c := <-s.signalIn:
			s.clients[c] = 1
		case c := <-s.signalOut:
			delete(s.clients, c)
		}
	}
}

//返回是否 結束
func (s *server) runCommand(cmd int) bool {
	switch cmd {
	case cmdRun:
		if s.run { //已經運行 直接返回
			return false
		}

		s.run = true
		go s.accept()
	case cmdClose:
		if !s.run { //未運行 直接返回
			return true
		}

		s.run = false
		s.listener.Close()
		s.Wait()
		for c, _ := range s.clients {
			c.Close()
		}
		s.clients = nil
		return true
	}
	return false
}
func (s *server) accept() {
	l := s.listener
	for {
		c, e := l.Accept()
		if e != nil {
			if !s.run {
				break
			}
			continue
		}
		go s.read(c)
	}
	close(s.signalWait)
}
func (s *server) read(c net.Conn) {
	//創建session
	template := s.template
	session, e := template.NewSession(c)
	if e != nil {
		c.Close()
		return
	}
	s.signalIn <- c

	//銷毀 session
	defer func() {
		template.DeleteSession(c, session)
		c.Close()
		s.signalOut <- c
	}()

	var buffer bytes.Buffer
	b := make([]byte, s.bufLen)
	size := -1
	headerSize := template.GetHeaderSize()
	var timer *time.Timer
	for {
		if s.timeout != 0 {
			timer = time.AfterFunc(s.timeout, func() {
				c.Close()
			})
		}
		n, e := c.Read(b)
		if timer != nil {
			timer.Stop()
		}
		if e != nil {
			return
		}
		_, e = buffer.Write(b[:n])
		if e != nil {
			return
		}

		for {
			//讀取 header
			if size == -1 {
				if buffer.Len() < headerSize {
					//等待 header
					break
				}
				buf := buffer.Bytes()
				size, e = template.GetMessageSize(session, buf[:headerSize])
				if e != nil || size < headerSize {
					//錯誤的 消息
					return
				}
			}

			//讀取body
			if buffer.Len() < size {
				//等待 body
				break
			}
			buf := make([]byte, size)
			_, e = buffer.Read(buf)
			if e != nil {
				//讀取錯誤
				return
			}
			//通知 處理消息
			e = template.Message(c, session, buf)
			if e != nil {
				//斷開連接
				return
			}

			//重置 消息 解析狀態
			size = -1
		}
	}
}
func (s *server) Run() {
	s.cmd <- cmdRun
}
func (s *server) Close() {
	s.cmd <- cmdClose
}
func (s *server) IsRun() bool {
	return s.run
}
func (s *server) Wait() {
	<-s.signalWait
}
