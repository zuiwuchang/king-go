package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	addr := "127.0.0.1:1102"

	c, e := net.Dial("tcp", addr)
	if e != nil {
		log.Fatalln(e)
	}
	defer c.Close()
	log.Println("connect at", addr)

	go func(c net.Conn) {
		var cmd string
		for {
			fmt.Print(">")
			fmt.Scan(&cmd)
			cmd += "end"
			c.Write([]byte(cmd))
		}
	}(c)

	b := make([]byte, 1024)
	for {
		n, err := c.Read(b)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("read", string(b[:n]))
	}

}
