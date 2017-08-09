package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"king-go/net/echo"
	"log"
	"time"
)

const (
	Addr       = "127.0.0.1:1102"
	TimeOut    = 0
	HeaderFlag = 1102
)

func main() {
	last := time.Now()

	count := 10000
	ch := make(chan int)

	for i := 0; i < count; i++ {
		go test_one(ch)
	}
	for count > 0 {
		<-ch
		count--
	}
	fmt.Println(time.Now().Sub(last))
}
func test_one(ch chan int) {
	defer func() {
		ch <- 1
	}()

	//連接服務器
	c, e := echo.NewClient(Addr, echo.NewClientTemplate(binary.LittleEndian, HeaderFlag))
	if e != nil {
		log.Fatalln(e)
	}
	defer c.Close()
	e = requestStr(c, "i'm king")
	if e != nil {
		log.Fatalln(e)
	}
	e = requestStr(c, "cerberus is an idea")
	if e != nil {
		log.Fatalln(e)
	}
}
func requestStr(c echo.IClient, str string) error {
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
func writeStr(c echo.IClient, str string) error {
	order := binary.LittleEndian
	n := len(str) + 4
	b := make([]byte, n)
	order.PutUint16(b, HeaderFlag)
	order.PutUint16(b[2:], uint16(n))
	copy(b[4:], []byte(str))

	_, e := c.Write(b)
	return e
}
