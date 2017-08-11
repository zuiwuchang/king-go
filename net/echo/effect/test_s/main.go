package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"king-go/net/echo"
	"log"
	"net"
	"os"
)

const (
	LAddr      = ":1102"
	TimeOut    = 0
	HeaderFlag = 1102
)

type Server struct {
	impl echo.IServerTemplate
}

func (s *Server) GetHeaderSize() int {
	return s.impl.GetHeaderSize()
}
func (s *Server) GetMessageSize(session echo.Session, header []byte) (n int, e error) {
	return s.impl.GetMessageSize(session, header)
}
func (s *Server) NewSession(c net.Conn) (session echo.Session, e error) {
	//fmt.Println("one in :", c.RemoteAddr())
	return nil, e
}
func (s *Server) DeleteSession(c net.Conn, session echo.Session) {
	//fmt.Println("one out :", c.RemoteAddr())
}
func (s *Server) Message(c net.Conn, session echo.Session, b []byte) error {
	return s.impl.Message(c, session, b)
}
func main() {
	//創建服務器 模板
	server := &Server{
		impl: echo.NewServerTemplate(binary.LittleEndian, HeaderFlag),
	}

	//創建服務器
	s, e := echo.NewServer(LAddr, TimeOut, server)
	if e != nil {
		log.Fatalln(e)
	}
	log.Println("work at", LAddr)

	//運行服務器
	s.Run()

	go getInput(s)

	//等待服務器 停止
	s.Wait()
}
func getInput(s echo.IServer) {
	r := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("\n$>")
		line, _, e := r.ReadLine()
		if e != nil {
			log.Fatalln(e)
		}
		cmd := string(line)

		if cmd == "e" {
			s.Close()
			break
		}
	}
}
