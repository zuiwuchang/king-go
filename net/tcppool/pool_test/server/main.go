//pool 測試用的 服務器
package main

import (
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	addr := ":1102"
	s, e := net.Listen("tcp", addr)
	if e != nil {
		log.Fatalln(e)
		return
	}
	defer s.Close()
	log.Println("work at", addr)

	for {
		c, e := s.Accept()
		if e != nil {
			log.Println(e)
			continue
		}
		go read(c)
	}
}
func read(c net.Conn) {
	log.Printf("%v in\n", c.RemoteAddr())
	buf := make([]byte, 1024)
	defer c.Close()

	success := false
	var str string
	for {
		n, e := c.Read(buf)
		if e != nil {
			log.Printf("%v out\n", c.RemoteAddr())
			return
		}
		str += string(buf[:n])
		if !strings.HasSuffix(str, "end") {
			//數據未收完 等待數據
			continue
		}

		//移除標識碼 獲取數據
		cmd := str[0 : len(str)-3]
		if success {
			fmt.Println("echo", cmd)
			str = ""
			if strings.HasPrefix(cmd, "send=") {
				c.Write([]byte("ok"))
			}
		} else if cmd == "king" {
			c.Write([]byte("welcome"))
			success = true
			str = ""
			log.Printf("%v ok\n", c.RemoteAddr())
		} else {
			//斷開 連接
			break
		}

	}
	log.Printf("%v out\n", c.RemoteAddr())
}
