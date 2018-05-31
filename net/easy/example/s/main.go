package main

import (
	"fmt"
	"github.com/zuiwuchang/king-go/net/easy"
	"log"
	"net"
	"runtime"
	"time"
)

func main() {
	l, e := net.Listen("tcp", "127.0.0.1:2920")
	if e != nil {
		log.Fatalln(e)
	}
	go func() {
		for {
			c, e := l.Accept()
			if e != nil {
				continue
			}
			go read(easy.NewClient(c, 1024, easy.Analyze()))
		}
	}()
	var cmd string
	for {
		fmt.Print("#>")
		fmt.Scan(&cmd)
		if cmd == "info" {
			fmt.Println("NumGoroutine", runtime.NumGoroutine())
		}
	}
}
func read(c easy.IClient) {
	timeout := time.Minute
	var msg []byte
	var e error
	var n int
	for {
		msg, e = c.ReadTimeout(timeout, nil)
		if e != nil {
			break
		}
		n, e = c.Write(msg)
		if e != nil || n != len(msg) {
			break
		}
	}
	c.Close()
	if e == easy.ErrorReadTimeout {
		c.WaitRead()
	} else if e == easy.ErrorWriteTimeout {
		c.WaitWrite()
	}
}
