package ini

import (
	"fmt"
	"log"
)

// 此注释将会被展示在页面上
// 此函数将被展示在OverView区域
func Example() {
	type Info struct {
		//ini tag 指定 key名稱
		//dec tag 指定 註釋
		Addr string `ini:"12345678" dec:"服務器地址"`
		//默認使用 屬性名 作爲 key
		Key int
		Lv  *int
		//忽略 未導出 屬性
		kk int
	}
	type Role struct {
		Name string
		Lv   int
	}

	//write
	lv := 6
	info := &Info{Addr: ":1102", Key: 10, Lv: &lv}
	role := &Role{Name: "kate", Lv: 8}
	b, e := Marshal(info, role)
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Println(string(b))

	//read
	info1 := &Info{kk: -1}
	role1 := &Role{}
	e = Unmarshal(b, info1, role1)
	if e != nil {
		log.Fatalln(e)
	}
	fmt.Println("info", info1.Addr, info1.Key, *info1.Lv, info1.kk)
	fmt.Println("role", role1.Name, role1.Lv)
	//Output:
	//[Info]
	//;服務器地址
	//12345678=:1102
	//Key=10
	//Lv=6
	//[Role]
	//Name=kate
	//Lv=8
	//
	//info :1102 10 6 -1
	//role kate 8
}
