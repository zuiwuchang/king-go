package main

import (
	"fmt"
	"github.com/zuiwuchang/king-go/net/easy"
	"log"
	"net"
	"time"
)

func main() {
	count := 20000
	//count = 1
	cs := make([]easy.IClient, count)
	for i := 0; i < count; i++ {
		ok := false
		for !ok {
			c, e := net.Dial("tcp", "127.0.0.1:2920")
			if e != nil {
				log.Println(e)
			}
			ok = true
			cs[i] = easy.NewClient(c, 1024, easy.Analyze())
		}
	}

	var cmd string
	for {
		fmt.Print("#>")
		fmt.Scan(&cmd)
		if cmd == "r" {
			request(cs)
		}
	}
}
func request(cs []easy.IClient) {
	last := time.Now()
	sum := 0
	ch := make(chan int)
	for i := 0; i < len(cs); i++ {
		c := cs[i]
		if c == nil {
			continue
		}
		sum++

		go func(i int) {
			var e error
			c := cs[i]
			defer func() {
				if e == nil {
					ch <- 0
				} else {
					c.Close()
					if e == easy.ErrorReadTimeout {
						c.WaitRead()
					} else if e == easy.ErrorWriteTimeout {
						c.WaitWrite()
					}
					cs[i] = nil
					ch <- 1
				}
			}()
			msg := make([]byte, easy.DefaultHeaderLen+8)
			easy.FormatMessage(msg, 1)
			var n int
			n, e = c.WriteTimeout(msg, time.Second*10)
			if e != nil {
				return
			} else if n != len(msg) {
				e = fmt.Errorf("busy write")
				return
			}

			msg, e = c.ReadTimeout(time.Second*10, nil)
		}(i)
	}
	ok := 0
	err := 0
	for sum != 0 {
		v := <-ch
		if v == 0 {
			ok++
		} else {
			err++
		}
		sum--
	}
	fmt.Printf("request success=%v error=%v %v\n", ok, err, time.Now().Sub(last))
}
