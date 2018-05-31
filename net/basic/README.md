# Discard
此 包已經被 廢棄 請不要再使用

# basic
tcp服務器 通常都需要 做如下繁瑣的 事
1. 接收客戶端 連接
2. 處理客戶端 發送來的消息
3. 一直重複 2 ...
4. 處理客戶斷開
5. 斷開 超時 未活動的 客戶
然 每次寫服務器 都只有 處理客戶 發送來的 消息 不同 其它步驟都是 一樣的 每次重寫 既容易出錯 又繁瑣 basic 包 完成了這最基礎的 工作 讓使用者 只需要關係 如何處理消息 即可  

```Go
package main

import (
	"github.com/zuiwuchang/king-go/net/basic"
	"log"
	"time"
)

const (
	LAddr = ":1102"
)

func main() {
	s, e := basic.NewServer(LAddr, //監聽地址
		time.Minute*10,               //客戶端 超時 時間 如果爲0 永不超時
		basic.NewServerTemplate(nil), //使用 默認的 服務器 模板
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
type IServerTemplate interface {
    //如何 創建 session 連接成功時 自動回調
    //如果 e != nil 將斷開連接
    NewSession(c net.Conn) (session Session, e error)

    //如何 銷毀 session 連接斷開時 自動 回調
    DeleteSession(c net.Conn, session Session)

    //如何 響應 接收到的 數據
    //如果 e != nil 將斷開連接
    Message(c net.Conn, session Session, b []byte) (e error)
}
```
IServerTemplate 接口是 使用者 需要實現的
1. 每當 成功爲 客戶 建立連接 NewSession 會被調用 你可選爲此客戶 創建一個 關聯的 Session (interface{}) 或者直接返回 nil
2. 每當收到 客戶數據 Message 會被 調用 同時 會傳入 關聯的 Conn 和 Session
3. DeleteSession 和 NewSession 一一對應 當close Conn 後 會調用 DeleteSession 以便 你能清除 Session 佔用的資源

# 注意
不要 手動 close Conn 而是 在 IServerTemplate 接口的 NewSession 或 Message 接口中 返回 error 以通知 服務器 斷開這個 Conn
