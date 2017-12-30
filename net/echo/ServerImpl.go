package echo

import (
	kio "github.com/zuiwuchang/king-go/io"
	"io"
	"net"
	"time"
)

//echo 服務器 實現
type serverImpl struct {
	Server *Server

	//服務器 是否運行中
	run bool

	//控制指令
	cmd chan int

	signalWait chan int

	clients   map[net.Conn]int
	signalIn  chan net.Conn
	signalOut chan net.Conn
}

func newServerImpl(srv *Server) *serverImpl {
	srv.format()

	s := &serverImpl{
		Server: srv,

		cmd:        make(chan int),
		signalWait: make(chan int),

		clients:   make(map[net.Conn]int),
		signalIn:  make(chan net.Conn),
		signalOut: make(chan net.Conn),
	}

	go s.runController()
	return s
}

//運行 控制 goroutine
func (s *serverImpl) runController() {
	defer s.Server.Listener.Close()

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
func (s *serverImpl) runCommand(cmd int) bool {
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
		s.Server.Listener.Close()
		s.Wait()
		for c, _ := range s.clients {
			c.Close()
		}
		s.clients = nil
		return true
	}
	return false
}
func (s *serverImpl) accept() {
	l := s.Server.Listener
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
func (s *serverImpl) read(c net.Conn) {
	//創建session
	template := s.Server.Template
	session, e := template.NewSession(c)
	if e != nil {
		c.Close()
		return
	}
	s.signalIn <- c

	var timer *time.Timer
	//銷毀 session
	defer func() {
		if timer != nil {
			timer.Stop()
		}
		template.DeleteSession(c, session)
		c.Close()
		s.signalOut <- c
	}()

	headerSize := template.GetHeaderSize()
	timeout := s.Server.Timeout

	buffer := make([]byte, s.Server.RecvBuffer)
	header := buffer[:headerSize]
	var msg []byte
	var size int
	for {
		//啓動 超時 關閉
		if timeout != 0 {
			timer = time.AfterFunc(timeout, func() {
				c.Close()
			})
		}

		//讀取 header
		e = kio.ReadAll(io.LimitReader(c, int64(headerSize)), header)
		if e != nil {
			return
		}
		//獲取 消息長度
		size, e = template.GetMessageSize(session, header)
		if e != nil || size < headerSize {
			//錯誤的 消息
			return
		}
		//緩衝區不夠
		if len(buffer) < size {
			msg = make([]byte, size)
			copy(msg, buffer[:headerSize])
		} else {
			msg = buffer[:size]
		}
		//讀取 body
		e = kio.ReadAll(c, msg[headerSize:])
		if e != nil {
			//讀取錯誤
			return
		}

		//通知 處理消息
		e = template.Message(c, session, msg)
		if e != nil {
			//斷開連接
			return
		}

		if timer != nil {
			timer.Stop()
		}
	}
}
func (s *serverImpl) Run() {
	s.cmd <- cmdRun
}
func (s *serverImpl) Close() {
	s.cmd <- cmdClose
}
func (s *serverImpl) IsRun() bool {
	return s.run
}
func (s *serverImpl) Wait() {
	<-s.signalWait
}
