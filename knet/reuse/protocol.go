package reuse

const (
	//默認 recv 窗口大小 64 k
	DefaultRecvWindow = 1024 * 64
	//默認 recv 緩衝區 大小 32k
	DefaultRecvBuffer = 1024 * 16

	//Accept 緩衝大小
	DefaultAcceptN = 10
)
const (
	// c -> s 客戶端 請求建立 連接
	// id(uint32) recv 窗口大小 (uint32)
	NetAccept = iota + 1
	// s -> c 成功建立 連接
	// id(uint32)
	NetAcceptOk
	// s -> c 服務器 拒絕建立連接
	NetAcceptReject

	// s <-> c 轉發對端 發來的數據到 上層
	// id(uint32) len(uint32) data([]byte)
	NetForward
	// s <-> c 向遠端 通知 已處理的 數據
	// id(uint32) len(uint32)
	//
	// 每當處理完 recv 窗口大小 * 0.7 的數據 就向對端發送 ack 確認
	NetAck

	// s <-> c 通知 對端 關閉連接
	// id(uint32)
	NetClose
)
