package main

import (
	"fmt"
	"king-go/net/tcppool"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

type myPoolHow struct {
}

func (m myPoolHow) Conect() (net.Conn, error) {
	c, err := net.Dial("tcp", "127.0.0.1:1102")
	if err != nil {
		return nil, err
	}
	_, err = c.Write([]byte("kingend"))
	if err != nil {
		return nil, err
	}

	b := make([]byte, 1024)
	n, err := c.Read(b)
	if err != nil {
		return nil, err
	}
	str := string(b[0:n])
	if str != "welcome" {
		return nil, err
	}

	_, err = c.Write([]byte("pool oneend"))
	if err != nil {
		return nil, err
	}
	return c, nil
}
func (m myPoolHow) Close(c net.Conn) error {
	return c.Close()
}
func (m myPoolHow) InitConn() int {
	return tcppool.POOL_MIN_CONN
}
func (m myPoolHow) MinConn() int {
	return tcppool.POOL_MIN_CONN
}
func (m myPoolHow) MaxConn() int {
	return tcppool.POOL_MAX_CONN
}
func (m myPoolHow) MaxAddStep() int {
	return 256
}
func (m myPoolHow) ResizeMore(use, free int) (bool, int) {
	sum := use + free
	if sum < 20 {
		return free < 3, 0
	} else if sum < 50 {
		return free < 10, 0
	}
	return free < use/5, 0
}
func (m myPoolHow) Ping(c *tcppool.Conn) error {
	_, err := c.Write([]byte("pingend"))
	return err
}
func (m myPoolHow) PingDuration() time.Duration {
	return tcppool.POOL_MIN_PING_DURATION
}
func (m myPoolHow) PingIntervalDuration() time.Duration {
	return tcppool.POOL_MIN_PING_INTERVAL_DURATION
}
func main() {
	//log.SetFlags(log.LstdFlags | log.Llongfile)

	pool := tcppool.NewPool(myPoolHow{})
	log.Println("pool create success")

	cs := make(map[*tcppool.Conn]bool)
	for {
		fmt.Print("$>")

		var cmd string
		fmt.Scan(&cmd)
		if cmd == "q" {
			break
		} else if cmd == "sum" {
			fmt.Println("sum =", pool.GetSum())
		} else if cmd == "free" {
			fmt.Println("free =", pool.GetFree())
		} else if cmd == "use" {
			fmt.Println("use =", pool.GetUse())
		} else if cmd == "new" {
			c, err := pool.Malloc()
			if err != nil {
				fmt.Println("error ", err)
			}
			cs[c] = true
		} else if cmd == "delete" {
			for c, _ := range cs {
				pool.Free(c)
				delete(cs, c)
				break
			}
		} else if strings.HasPrefix(cmd, "new=") {
			cmd = cmd[len("new="):]
			n, _ := strconv.ParseInt(cmd, 10, 32)
			for n > 0 {
				c, err := pool.Malloc()
				if err != nil {
					fmt.Println("error ", err)
					break
				}
				cs[c] = true
				n--
			}
		} else if strings.HasPrefix(cmd, "delete=") {
			cmd = cmd[len("delete="):]
			n, _ := strconv.ParseInt(cmd, 10, 32)
			for c, _ := range cs {
				if n < 1 {
					break
				}
				pool.Free(c)
				delete(cs, c)
				n--
			}
		} else if strings.HasPrefix(cmd, "send=") {
			b := make([]byte, 1024, 1024)
			for c, _ := range cs {
				if _, err := c.Write([]byte(cmd + "end")); err != nil {
					pool.Free(c)
					delete(cs, c)
					log.Println(err)
				}
				if n, err := c.Read(b); err != nil {
					pool.Free(c)
					delete(cs, c)
					log.Println(err)
				} else if string(b[:n]) != "ok" {
					c.IsOk = false
					pool.Free(c)
					delete(cs, c)
					log.Println("Read error")
				}
			}
			fmt.Println("send ok")
		}
	}
}
