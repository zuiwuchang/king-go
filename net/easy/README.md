# easy
tcp需求的變化和複雜性 讓孤發現 想要為其構造一個 靈活 高效的 框架 是枉費心機 孤想不到一個框架實現方式能很好 適應到各種不同的 tcp應用 

應該只為其提供 簡單的常用工具 讓後依據 需求 組裝這些 工具 為此 孤廢棄了以前 實現的 tcp 庫 重新 提供了 easy 包

easy 包 只包含 幾個 基礎的常用 功能 
1. IListener 接口 就和 net.Listener 一樣 不過多了一個 Closed()bool
2. IAnalyze 接口 讓 包知道如何 拆組 tcp流中的 數據包
3. IClient 類似 net.Coon 不過提供了 超時讀寫 以及配合 IAnalyze 自動解包

# IListener
net.Listener 在closed 後 accept 會返回 錯誤 但 無法 獲取到錯誤是否是因為 Listener 已經關閉 這在需要 動態切換 listen 地址時 是有用的 
故 IListener 在 net.Listener 接口之上 提供了一個 Closed() bool 函數 返回 net.Listener 是否已經 關閉
```go
// IListener 在 net.Listener 之上 提供了 一個 Closed 返回 Listener 是否已經關閉
type IListener interface {
	net.Listener
	// 返回 Listener 是否已經 關閉
	Closed() bool
}
func NewListener(l net.Listener) IListener{
  ...
}
```
func NewListener(l net.Listener) IListener 將一個 net.Listener 轉為 IListener NewListener 會 創建一個 包含 net.Listener 和 closed bool 的結構指針
並且 重寫 Close 函數 當 你調用 Close 時將 closed 設置為 true

# IAnalyze
IAnalyze 定義了 如何 從 tcp流中 解包 IClient 需要 IAnalyze 指導如何 從 tcp流中 獲取數據包
```go
// IAnalyze 定義解包接口
type IAnalyze interface {
	//返回 包頭 長度
	Header() int
	//傳入 包頭 返回 包長 len(header + body)
	Analyze(header []byte) (int, error)
}
// Analyze 返回 默認的 接包器
func Analyze() IAnalyze {
	return defaultAnalyze{}
}
```
Analyze() 函數 返回一個 默認的 解包接口 使用 binary.LittleEndian 以 6字節的header + body 組成
header 由 2字節的flag(1911) + 2字節消息指令 + 2字節包長 組成

```go
// AnalyzeFunc 創建一個 IAnalyze 接口
func AnalyzeFunc(
	headerSize int,
	analyzeFunc func(header []byte,
	) (int, error)) IAnalyze {
	return &analyze{
		headerSize:  headerSize,
		analyzeFunc: analyzeFunc,
	}
}
```
通常 你需要 定義一個 struct 之後 實現 IAnalyze 接口 不過 easy 提供了一個 AnalyzeFunc 輔助函數 讓你可以 傳入 一個 包頭長+解包函數 直接創建一個 IAnalyze 接口

# IClient
IClient 接口 是 easy 包最主要的 功能 其和 net.Conn 接口 高度類似 只是在 Read 時 每次都 固定返回一個 完整的 消息包 且 提供了 超時讀寫

你 基本上 可以 直接 把 IClient 當作 net.Conn 使用 但要注意 下面 幾點
1. Read() ReadTimeout() 如果返回了數據 一定是返回了一個 完整的 消息包
2. Write() WriteTimeout() 和 net.Conn 的 Write() 保持一直 實際上 內部就是直接調用的 net.Conn 的 Write()
3. 當使用 超時 讀寫 ReadTimeout() WriteTimeout() 時 你必須判斷 返回的 錯誤是否是 ErrorReadTimeout ErrorWriteTimeout 
如果是則需要調用 WaitRead() WaitWrite（）以便 讀寫的 goroutine 能夠返回 chan<- 
4. 當出現 ErrorReadTimeout ErrorWriteTimeout 實際的 tcp 錯誤 和 讀寫字節數量 是由 WaitRead() WaitWrite() 返回

```go
// IClient 定義一個 tcp 接口
// ReadXXX 是不保證 goroutine 安全的 你不能在 一個 ReadXXX 返回前 在 其它 goroutine 中 調用 ReadXXX
// WriteXXX 有同樣的 限制
// tcp 的 讀寫流是 分開了 所以 你可以在 一個 goroutine ReadXXX 返回前 在另外一個 goroutine 中 調用 WriteXXX 反之亦然
type IClient interface {
	// 從 tcp 流中 讀取一個 完整的 消息
	// msgReuse 是一個 可選的 消息接收緩存 如果 msgReuse 足夠長 msg 則直接從 msgReuse[:msgSize]中 切片 否則 申請新內存 make([]byte,msgSize)
	Read(msgReuse []byte) (msg []byte, e error)
	// 從 tcp 流中 讀取一個 完整的 消息
	// 注意 如果 ReadTimeout 返回 ErrorReadTimeout 你應該 調用 WaitRead 等待 read 以便 read 的 goroutine 能夠正常退出
	// 此時 你可以 調用 Close 以便 WaitRead 能夠 立刻 返回  但着意味着 此後 將 不能在進行任何 讀寫操作
	//
	// 注意 如果 ReadTimeout 返回 ErrorReadTimeout n 的值爲0
	ReadTimeout(timeout time.Duration, msgReuse []byte) (msg []byte, e error)

	// 通常 只有在 ReadTimeout 失敗後 必須手動 調用 此函數 以便 read 的 goroutine能夠 正確退出
	// 因爲 受限與socket recv 函數 一旦被調用在 函數 成功或失敗 前都 無法使此函數 返回 或 取消操縱
	// 這意味着 即使你調用了 ReadTimeout 如果 超時了 你還是只有 調用 WaitWrite 等待 write結束
	// 不過 WriteTimeout 依然有意義 你可以選擇 直接 在此 調用 Close 後調用 WaitWrite 以便立刻 關閉這個 超時的 socket
	//
	// WaitRead 不會返回 實際寫入從 tcp recv 緩存流中 得到的數據大小 但會返回 錯誤信息
	// 如果 你沒有 調用 Close socket 沒有 錯誤 WaitRead 會在讀取到 完整 消息後 返回 msg,nil 此時 你可以繼續使用此 client
	WaitRead() (msg []byte, e error)

	// 向 tcp 流 寫入 數據流
	// 注意 如果失敗 n 可能是 [0,len(msg)] 意味着可能寫入了 部分數據
	Write(b []byte) (n int, e error)
	// 向 tcp 流 寫入 數據流
	// 注意 如果失敗 n 可能是 [0,len(msg)] 意味着可能寫入了 部分數據
	//
	// 注意 如果 WriteTimeout 返回 ErrorWriteTimeout 你應該 調用 WaitWrite 等待 write 以便 write的 goroutine 能夠正常退出
	// 此時 你可以 調用 Close 以便 WaitWrite 能夠 立刻 返回  但着意味着 此後 將 不能在進行任何 讀寫操作
	//
	// 注意 如果 WriteTimeout 返回 ErrorWriteTimeout n 的值爲0 但 實際寫入到 tcp send 緩存流中 的數據大小 應該由 WaitWrite 返回
	WriteTimeout(b []byte, timeout time.Duration) (n int, e error)

	// 通常 只有在 WriteTimeout 失敗後 必須手動 調用 此函數 以便 write 的 goroutine能夠 正確退出
	// 因爲 受限與socket send 函數 一旦被調用在 函數 成功或失敗 前都 無法使此函數 返回 或 取消操縱
	// 這意味着 即使你調用了 WriteTimeout 如果 超時了 你還是只有 調用 WaitWrite 等待 write結束
	// 不過 WriteTimeout 依然有意義 你可以選擇 直接 在此 調用 Close 後調用 WaitWrite 以便立刻 關閉這個 超時的 socket
	//
	// WaitWrite 會返回 實際寫入到 tcp send 緩存流中 的數據大小 和 錯誤信息
	WaitWrite() (n int, e error)

	// 關閉 socket 連接
	Close() (e error)

	// 返回 遠端 地址
	RemoteAddr() net.Addr
	// 返回 本端 地址
	LocalAddr() net.Addr
}
// NewClient 將一個 net.Conn 包裝爲 IClient
func NewClient(c net.Conn, bufferSize int, analyze IAnalyze) IClient {
  ...
}
```
NewClient 將一個 net.Conn 包裝為 IClient 其參數 分別是 要包裝的 net.Conn 緩衝區 解包接口
IClient 需要自動 解包 為了 提高內存 利用率 會申請一個 緩衝區 大小為 bufferSize 重複使用 這個 緩衝區 read數據 故 消息包 最大不能超過 此值
