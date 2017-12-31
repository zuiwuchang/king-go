//建立在tcp之上的一些工具
package knet

import (
	"net"
)

const (
	DefaultHeaderLen = 6

	DefaultFlag = 1911

	DefaultFlagPos    = 0
	DefaultCommandPos = 2
	DefaultLenPos     = 4

	DefaultRecvBuffer = 1024 * 32
)

type Description struct {
	Analyze    IAnalyze
	RecvBuffer int

	headerSize int
}

func (d *Description) Format() {
	if d.Analyze == nil {
		if kLog.Info != nil {
			kLog.Info.Println("use default Analyze")
		}
		d.Analyze = defaultAnalyze{}
	}
	d.headerSize = d.Analyze.Header()
	if d.headerSize < 1 {
		if kLog.Fault != nil {
			kLog.Fault.Println("Header must larger than 0")
		}
		panic("Header must larger than 0")
	}

	if d.RecvBuffer < d.headerSize {
		if kLog.Warn != nil {
			kLog.Warn.Printf("RecvBuffer(%v) < Header(%v) use default(%v)",
				d.RecvBuffer,
				d.headerSize,
				DefaultRecvBuffer,
			)
		}

		d.RecvBuffer = DefaultRecvBuffer
	}
}

//一個 socket 連接
type IConn interface {
	net.Conn

	//讀取一條完整的 消息
	//
	//如果 返回 任何 錯誤 則 socket 已經不可靠 應該關閉之
	ReadMessage() ([]byte, error)

	//讀取一條消息的 header
	//
	//多次 調用 ReadHeader 不會將 message 移除緩衝區 你只會得到同一個消息的 header
	//
	//你必須 調用 ReadMessage/WriteTo/WriteSkipTo 才能將 message 移除緩衝區 之後繼續讀取 後續的 緩衝區
	ReadHeader() ([]byte, error)
}

//一個 監聽狀態的 socket 服務器
type ILister interface {
	//返回 loster 監聽 地址
	Addr() net.Addr
	//接受一個 客戶端的 連接請求
	Accept() (IConn, error)
	//關閉 socket 此後 無法在繼續 執行 Accept 操作
	Close() error
}

//一個 可以連接 服務器的 撥號器
type IDialer interface {
	//向服務器 請求創建一個 連接
	Dial(network, address string) (IConn, error)
}

//tcp 解包器
type IAnalyze interface {
	//返回 包頭 長度
	Header() int
	//傳入 包頭 返回 包長 len(header + body)
	Analyze(header []byte) (int, error)
}

//創建一個 默認的 IAnalyze 接口
//header 爲 uint16 的 binary.LittleEndian 包長數據
func NewAnalyze() IAnalyze {
	return defaultAnalyze{}
}

//創建一個 IAnalyze 接口
func AnalyzeFunc(
	headerSize int,
	analyzeFunc func(header []byte,
	) (int, error)) IAnalyze {
	return &analyze{
		headerSize:  headerSize,
		analyzeFunc: analyzeFunc,
	}
}
