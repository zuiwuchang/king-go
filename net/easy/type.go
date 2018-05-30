// Package easy 簡化 tcp 編程
// tcp 是 十分靈活的流協議 通常這樣的靈活 帶來了一些編碼上的 繁瑣 例如 從流中 拆包組包 ...
// easy 包 使用一些 約定來 約束tcp 限制了一些靈活性 以換來 編程上的 一些簡化
// 同時 完成了 一些 tcp 編程中 常用的功能 例如 帶超時的讀寫
package easy

import (
	"net"
	"time"
)

// 定義了一個 默認的 tcp 包協議
const (
	// DefaultHeaderLen 默認 包頭長度
	DefaultHeaderLen = 6

	// DefaultFlag 默認 包 標記
	DefaultFlag = 1911

	// DefaultFlagPos 默認的 包 標記 偏移
	DefaultFlagPos = 0
	// DefaultCommandPos 默認 包 命名偏移
	DefaultCommandPos = 2
	// DefaultLenPos 默認 包 長度 偏移
	DefaultLenPos = 4

	// DefaultRecvBuffer 默認 recv 緩衝區 大小
	DefaultRecvBuffer = 1024 * 16
)

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

// IClient 定義一個 tcp 接口
// ReadXXX 是不保證 goroutine 安全的 你不能在 一個 ReadXXX 返回前 在 其它 goroutine 中 調用 ReadXXX
// WriteXXX 有同樣的 限制
// tcp 的 讀寫流是 分開了 所以 你可以在 一個 goroutine ReadXXX 返回前 在另外一個 goroutine 中 調用 WriteXXX 反之亦然
type IClient interface {
	// 從 tcp 流中 讀取一個 完整的 消息
	Read() (msg []byte, e error)
	// 從 tcp 流中 讀取一個 完整的 消息
	// 注意 如果 ReadTimeout 返回 ErrorReadTimeout 你應該 調用 WaitRead 等待 read 以便 read 的 goroutine 能夠正常退出
	// 此時 你可以 調用 Close 以便 WaitRead 能夠 立刻 返回  但着意味着 此後 將 不能在進行任何 讀寫操作
	//
	// 注意 如果 ReadTimeout 返回 ErrorReadTimeout n 的值爲0
	ReadTimeout(timeout time.Duration) (msg []byte, e error)

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

// IListener 在 net.Listener 之上 提供了 一個 Closed 返回 Listener 是否已經關閉
type IListener interface {
	net.Listener
	// 返回 Listener 是否已經 關閉
	Closed() bool
}
