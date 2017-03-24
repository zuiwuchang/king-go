//解析 mir 使用的 圖像資源 wix wil/wzx wzl
package wzx

import (
	"fmt"
)

const (
	//bitmap 16 565
	COLOR_DEPTH_16_565 = 65536 //wil 標記
	DEPTH_16_565       = 0x105 //wzl 標記
	//bitmap 256
	COLOR_DEPTH_256 = 256   //wil 標記
	DEPTH_256       = 0x103 //wzl 標記

	//wix wil
	TYPE_WIX = 1
	//wzx wzl
	TYPE_WZX = 2
	//自定義 擴展格式 可用值
	TYPE_USER = 1000
)

//檔案 解析器定義
type Parser interface {
	NewSource(t int, xPath, lPath string) (Source, error)
}
type Img interface {
	GetSize() (uint16, uint16)
	GetSign() uint32
	GetData() ([]byte, error)

	//銷毀 資源
	Destory()
}

//資源定義
type Source interface {
	//返回 資源類型
	GetType() int

	//返回 圖像數量
	GetSize() int
	//返回 圖像偏移
	GetPos(i int) (int64, error)

	//返回 圖像 數據
	GetImg(i int) (Img, error)

	//銷毀 圖像 緩存數據
	ClearData(i int)

	//銷毀 資源
	Destory()
}

//全局 解析器
var g_Parser map[int]Parser

func init() {
	g_Parser = make(map[int]Parser)
	RegisterParser(TYPE_WIX, &wixParser{})
	RegisterParser(TYPE_WZX, &wzxParser{})
}

//註冊 資源 解析器
func RegisterParser(t int, parser Parser) {
	g_Parser[t] = parser
}

//創建一個 資源 檔案
func NewSource(t int, xPath, lPath string) (Source, error) {
	if parser, ok := g_Parser[t]; ok {
		return parser.NewSource(t, xPath, lPath)
	}
	return nil, fmt.Errorf("not found source type %v", t)
}
