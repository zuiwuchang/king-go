package wzx

import (
	"encoding/binary"
)

//位圖 img
type Bmp struct {
	Width  uint16
	Height uint16
	Sign   uint32

	Data []byte
}

func (i *Bmp) GetSize() (uint16, uint16) {
	return i.Width, i.Height
}
func (i *Bmp) GetSign() uint32 {
	return i.Sign
}
func (i *Bmp) GetData() ([]byte, error) {
	return i.Data, nil
}

//銷毀 資源
func (i *Bmp) Destory() {
	i.Data = nil
}

//位圖頭信息
type BmpHeader struct {
	//bm 標記 0x42 0x4D
	Flag [2]byte

	//檔案大小
	Size uint32

	//保留位 爲0
	Reserved1 uint16

	//保留位 爲0
	Reserved2 uint16

	//位圖數據 相對os.SEEK_SET 偏移
	Pos uint32
}

func (h *BmpHeader) ToBinary() []byte {
	b := make([]byte, 14, 14)
	copy(b, h.Flag[:])

	binary.LittleEndian.PutUint32(b[2:2+4], h.Size)
	binary.LittleEndian.PutUint32(b[10:], h.Pos)
	return b
}

//位圖詳情
type BmpInfo struct {
	//BmpInfo 結構大小 故 爲40
	Size uint32

	//圖像寬度
	Width uint32
	//圖像高度
	Height uint32

	//平面數 固定爲1
	Planes uint16

	//bit 數 2 4 8 16 24 32
	BitCount uint16

	//是否 壓縮
	//0	RI_RGB		不壓縮
	//1	RI_RLE8	8bit 編碼時 使用 rle算法
	//2	RI_RLE4	4bit 編碼時 使用 rle算法
	//3	RI_BITFIElDS	16/32 bit
	//4	RI_JPEG	包含 jpeg 圖像
	//5	RI_PNG		包含 png 圖像
	Compression uint32
	//圖像數據大小 biCompression爲RI_RGB時 此值可爲0
	SizeImage uint32

	//水平分辨率
	XPelsPerMeter uint32
	//垂直分辨率
	YPelsPerMeter uint32

	//調色板中顏色 使用數 爲0 全部使用
	ClrUsed uint32
	//重要顏色 索引數 爲0 全部重要
	ClrImportant uint32
}

func (i *BmpInfo) ToBinary() []byte {
	b := make([]byte, 40, 40)

	pos := 0
	binary.LittleEndian.PutUint32(b[pos:], i.Size)
	pos += 4
	binary.LittleEndian.PutUint32(b[pos:], i.Width)
	pos += 4
	binary.LittleEndian.PutUint32(b[pos:], i.Height)
	pos += 4

	binary.LittleEndian.PutUint16(b[pos:], i.Planes)
	pos += 2
	binary.LittleEndian.PutUint16(b[pos:], i.BitCount)
	pos += 2
	binary.LittleEndian.PutUint32(b[pos:], i.Compression)
	pos += 4

	binary.LittleEndian.PutUint32(b[pos:], i.SizeImage)
	pos += 4
	binary.LittleEndian.PutUint32(b[pos:], i.XPelsPerMeter)
	pos += 4
	binary.LittleEndian.PutUint32(b[pos:], i.YPelsPerMeter)
	pos += 4
	binary.LittleEndian.PutUint32(b[pos:], i.ClrUsed)
	pos += 4
	binary.LittleEndian.PutUint32(b[pos:], i.ClrImportant)

	return b
}
