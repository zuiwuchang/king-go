package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	os.MkdirAll("king-go/go-xorm/params", 0774)
	f, err := os.Create("king-go/go-xorm/params/params.go")
	n := 20

	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	write(f)
	writeWhere(f, n)

	log.Println("success")
}
func write(f *os.File) {
	f.WriteString("/*	靈活的爲 xorm 執行的sql 中的 where 創建 ? 參數 	*/\n")
	f.WriteString("package params\n\n")
	f.WriteString(`import (
	"bytes"
	"github.com/go-xorm/xorm"
)

type Params struct {
	buf    bytes.Buffer
	params []interface{}
}

//創建 一個 參數集
//cap	? 參數 數量的 參考值
func NewParams(cap int) Params {
	return Params{params: make([]interface{}, 0, cap)}
}

//增加 where 條件
func (p *Params) WriteWhere(where string) {
	p.buf.WriteString(where)
}

//增加 ?
func (p *Params) WriteParam(param interface{}) {
	p.params = append(p.params, param)
}`)
	f.WriteString("\n\n")
}
func writeWhere(f *os.File, n int) {
	f.WriteString("//執行 Where\n")
	f.WriteString("func (p *Params) Where(session *xorm.Session) {\n")
	f.WriteString("	switch len(p.params) {\n")
	for i := 0; i <= n; i++ {
		f.WriteString(fmt.Sprintf("	case %d:\n", i))
		f.WriteString("		session.Where(p.buf.String()")
		for j := 0; j < i; j++ {
			f.WriteString(fmt.Sprintf(", p.params[%d]", j))
		}
		f.WriteString(")\n")
	}
	f.WriteString("	}\n")
	f.WriteString("}\n")
}
