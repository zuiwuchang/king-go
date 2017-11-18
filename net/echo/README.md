# echo
echo 是一個類似 basic 的 tcp 服務器 框架

相對於 basic echo 會自動 進行 解包操縱 其操縱 基本和 [king-go/net/basic](https://github.com/zuiwuchang/king-go/tree/master/net/basic) 相同

echo 支持的包結構是 有 header + body 組成 且需要從 header 中獲取 包長 

```Go
package main

import (
	"encoding/binary"
	"github.com/zuiwuchang/king-go/net/echo"
	"log"
	"time"
)

const (
	LAddr = ":1102"
)

func main() {
	s, e := echo.NewServer(LAddr, //監聽地址
		time.Minute*10,                                   //客戶端 超時 時間 如果爲0 永不超時
		echo.NewServerTemplate(binary.LittleEndian, 731), //使用 默認的 服務器 模板
	)
	if e != nil {
		log.Fatalln(e)
	}
	//運行服務器
	s.Run()
	log.Println("work at", LAddr)

	//等待服務器停止
	s.Wait()
}
```

# 一般步驟
1. 定義一個 實現了 IServerTemplate 接口的 struct
2. 調用 NewServer 傳入 IServerTemplate 以創建 tcp 服務器
3. 調用 IServer 的 run 運行 服務器
4. 調用 IServer 的 wait 等待服務器 停止

# 一般步驟
1. 定義一個 實現了 IServerTemplate 接口的 struct
2. 調用 NewServer 傳入 IServerTemplate 以創建 tcp 服務器
3. 調用 IServer 的 run 運行 服務器
4. 調用 IServer 的 wait 等待服務器 停止

# NewServer
 func NewServer(laddr string, timeout time.Duration, template IServerTemplate) (IServer, error) 創建一個 tcp 服務器
 * laddr 指定服務器工作 地址
 * timeout 指定了 客戶端  多長時間 沒發送消息 則認爲已經斷線 從而 斷開連接
 * template 接口 定義了如何 主要邏輯 需要 用戶 自己實現 NewServerTemplate 是一個 接口實現示例
# IServer
```Go
type IServer interface {
    //運行服務器
    Run()

    //等待服務器 停止
    Wait()

    //返回服務器是否在運行中
    IsRun() bool

    //關閉服務器 並釋放所有資源
    Close()
}
```

# IServerTemplate
```Go
//echo 服務器 模板 定義如何 解析 消息
type IServerTemplate interface {
	//返回 header 長度 必須大於 0
	GetHeaderSize() int
	//傳入 header 返回 整個消息長
	//如果 e != nil 或 消息長 小於 headerSize 將斷開連接
	GetMessageSize(session Session, header []byte) (n int, e error)

	//如何 創建 session 連接成功時 自動回調
	//如果 e != nil 將斷開連接
	NewSession(c net.Conn) (session Session, e error)

	//如何 銷毀 session 連接斷開時 自動 回調
	DeleteSession(c net.Conn, session Session)

	//如何 響應 消息
	//如果 e != nil 將斷開連接
	Message(c net.Conn, session Session, b []byte) (e error)
}
```
IServerTemplate 接口和 basic.IServerTemplate 基本一樣 現在只是多出了 GetHeaderSize GetMessageSize 兩個函數
 * GetHeaderSize 定義一個消息包 的 header 長度
 * 當 服務器 收一個完整包長時 會回調 GetMessageSize 並 傳入 header 此函數 需要 返回 整個包長 len(header)+lean(body)
 * 現在 Message 不會在 收到數據時被調用 而是每當收到一個 完成消息包 時被回調 每個包 調用一次(如果同時收到多個包 回調多次 如果 包不完整 自動等待收完後回調)
 

# NewServerTemplate
NewServerTemplate 會常見一個 默認的 服務器 模板 其header 爲 2個字節的 標識 加上 2個字節的長度

功能 則是 將收到的 完整 數據包 回發給 客戶端

# IClient IClientTemplate
IClient IClientTemplate 是一個 自動 解包的 客戶端
```Go
//echo 客戶端 接口
type IClient interface {
	//讀取一個消息
	//timeout 讀取超時(如果超時 自動斷開連接) 爲0 永不超時
	GetMessage(timeout time.Duration) (b []byte, e error)
	net.Conn
}

//echo 客戶端 模板 定義如何 解析 消息
type IClientTemplate interface {
	//返回 header 長度 必須大於 0
	GetHeaderSize() int
	//傳入 header 返回 整個消息長
	//如果 e != nil 或 消息長 小於 headerSize 將斷開連接
	GetMessageSize(header []byte) (n int, e error)
}
```
```Go
package main

import (
	"encoding/binary"
	"github.com/zuiwuchang/king-go/net/echo"
	"log"
	"net"
)

const (
	Addr = "127.0.0.1:1102"
	Flag = 731
)

func main() {
	//創建一個 默認 客戶端
	c, e := echo.NewClient(Addr, echo.NewClientTemplate(binary.LittleEndian, Flag))
	if e != nil {
		log.Fatalln(e)
	}
	//結束時 關閉 客戶端
	defer c.Close()

	//發送 數據
	str := "cerberus is an idea"
	e = writeStrings(c, str)
	if e != nil {
		log.Fatalln(e)
	}

	//接受 服務器 返回消息
	b, e := c.GetMessage(0)
	if e != nil {
		log.Fatalln(e)
	}
	if string(b[4:]) != str {
		log.Fatalln("bad rs")
	}
	log.Println("yes")
}
func writeStrings(c net.Conn, str string) error {
	n := 4 + len(str)
	b := make([]byte, n)
	binary.LittleEndian.PutUint16(b, Flag)
	binary.LittleEndian.PutUint16(b[2:], uint16(n))
	copy(b[4:], []byte(str))

	_, e := c.Write(b)
	return e
}
```
