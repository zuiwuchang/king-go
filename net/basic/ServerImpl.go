package basic

import (
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

func newBasicServer(srv *Server) IServer {
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

	//銷毀 session
	defer func() {
		template.DeleteSession(c, session)
		c.Close()
		s.signalOut <- c
	}()

	b := make([]byte, s.Server.RecvBuffer)
	var timer *time.Timer
	timeout := s.Server.Timeout
	var n int
	for {
		if timeout != 0 {
			timer = time.AfterFunc(timeout, func() {
				c.Close()
			})
		}
		n, e = c.Read(b)
		if timer != nil {
			timer.Stop()
		}
		if e != nil {
			return
		}

		//通知 處理消息
		e = template.Message(c, session, b[:n])
		if e != nil {
			//斷開連接
			return
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
