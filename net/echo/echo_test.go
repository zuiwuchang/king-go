package echo

import (
	"encoding/binary"
	"errors"
	kio "github.com/zuiwuchang/king-go/io"
	"testing"
)

func TestEcho(t *testing.T) {
	const (
		LAddr      = ":1102"
		TimeOut    = 0
		HeaderFlag = 1102
		Addr       = "127.0.0.1:1102"
	)

	//創建服務器
	s, e := NewServer(LAddr, TimeOut, NewServerTemplate(binary.LittleEndian, HeaderFlag))
	if e != nil {
		t.Fatal(e)
	}

	//運行服務器
	s.Run()

	writeStr := func(c IClient, str string) error {
		order := binary.LittleEndian
		n := len(str) + 4
		b := make([]byte, n)
		order.PutUint16(b, HeaderFlag)
		order.PutUint16(b[2:], uint16(n))
		copy(b[4:], []byte(str))

		return kio.WriteAll(c, b)
	}
	requestStr := func(c IClient, str string) error {
		e := writeStr(c, str)
		if e != nil {
			return e
		}
		b, e := c.GetMessage(0)
		if e != nil {
			return e
		}
		if string(b[4:]) != str {
			return errors.New("rs not equal")
		}
		return nil
	}

	count := 1000
	ch := make(chan int)

	for i := 0; i < count; i++ {
		go func() {
			defer func() {
				ch <- 1
			}()

			//連接服務器
			c, e := NewClient(Addr, NewClientTemplate(binary.LittleEndian, HeaderFlag))
			if e != nil {
				t.Fatal(e)
			}
			defer c.Close()

			e = requestStr(c, "i'm king")
			if e != nil {
				t.Fatal(e)
			}
			e = requestStr(c, "cerberus is an idea")
			if e != nil {
				t.Fatal(e)
			}

		}()
	}
	for count > 0 {
		<-ch
		count--
	}
	s.Close()
	s.Wait()
}
