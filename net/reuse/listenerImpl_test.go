package reuse

import (
	"bufio"
	"net"
	"testing"
	"time"
)

const (
	testLAddr = "127.0.0.1:1731"
)

func TestListenerImpl(t *testing.T) {
	//lister
	l, e := newListenerImpl("tcp", testLAddr)
	if e != nil {
		t.Fatal(e)
	}
	ch := make(chan (int))
	go func() {
		for i := 0; i < 2; i++ {
			c, e := l.Accept()
			if e != nil {
				t.Fatal(e)
				return
			}

			go func() {
				defer c.Close()
				r := bufio.NewReader(c)
				for {
					b, _, e := r.ReadLine()
					if e != nil {
						return
					}

					if n, e := c.Write([]byte(string(b) + "\n")); e != nil || n != len(b) {
						return
					}
				}

			}()
		}
		ch <- 0
	}()

	//client
	c0, e := net.Dial("tcp", testLAddr)
	if e != nil {
		t.Fatal(e)
	}
	defer c0.Close()

	c1, e := net.Dial("tcp", testLAddr)
	if e != nil {
		t.Fatal(e)
	}
	defer c1.Close()

	//write
	s0 := "this is c0"
	s1 := "this is c1"
	n, e := c0.Write([]byte(s0 + "\n"))
	if e != nil {
		t.Fatal(e)
	} else if n != len(s0)+1 {
		t.Fatal("bad write")
	}
	n, e = c1.Write([]byte(s1 + "\n"))
	if e != nil {
		t.Fatal(e)
	} else if n != len(s1)+1 {
		t.Fatal("bad write")
	}
	//read
	r0 := bufio.NewReader(c0)
	r1 := bufio.NewReader(c1)
	if b, _, e := r1.ReadLine(); e != nil {
		t.Fatal(e)
	} else if string(b) != s1 {
		t.Fatal("bad read s1 :", string(b))
	}
	if b, _, e := r0.ReadLine(); e != nil {
		t.Fatal(e)
	} else if string(b) != s0 {
		t.Fatal("bad read s0 :", string(b))
	}

	//wait exit
	time.AfterFunc(time.Second, func() {
		ch <- 1
	})
	rs := <-ch
	if rs != 0 {
		t.Fatal("wait rs timeout")
	}
	e = l.Close()
	if e != nil {
		t.Fatal(e)
	}
}
