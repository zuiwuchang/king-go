package wzx

import (
	"encoding/binary"
	"fmt"
	"os"
)

type wilHeader struct {
	//圖像數量
	Size uint32
	//圖像格式
	ColorDepth uint32
}

type wixSource struct {
	//圖像偏移
	imgsPos []uint32

	//圖像緩存
	imgsCache []Img

	//數據檔案
	dataFile *os.File

	//wil 頭信息
	header *wilHeader
}

func (w *wixSource) GetType() int {
	return TYPE_WIX
}
func (w *wixSource) GetSize() int {
	return len(w.imgsPos)
}
func (w *wixSource) GetPos(i int) (int64, error) {
	if i < 0 || i >= len(w.imgsPos) {
		return 0, fmt.Errorf("index error (%v)", i)
	}
	return int64(w.imgsPos[i]), nil
}

func (w *wixSource) GetImg(i int) (Img, error) {
	if i < 0 || i >= len(w.imgsPos) {
		return nil, fmt.Errorf("index error (%v)", i)
	}
	if w.imgsCache[i] != nil {
		return w.imgsCache[i], nil
	}
	if w.header.ColorDepth == COLOR_DEPTH_256 {
		img, e := w.createImg256(i)
		if e != nil {
			return nil, e
		}
		w.imgsCache[i] = img
		return img, nil
	} else if w.header.ColorDepth == COLOR_DEPTH_16_565 {
		img, e := w.createImg565(i)
		if e != nil {
			return nil, e
		}
		w.imgsCache[i] = img
		return img, nil
	}
	return nil, fmt.Errorf("not found depth (%v)", w.header.ColorDepth)
}
func (w *wixSource) createImg256(i int) (Img, error) {
	pos, e := w.GetPos(i)
	if e != nil {
		return nil, e
	}
	f := w.dataFile
	if _, e = f.Seek(pos, os.SEEK_SET); e != nil {
		return nil, e
	}

	//w
	b := make([]byte, 2, 2)
	if _, e = f.Read(b); e != nil {
		return nil, e
	}
	width := binary.LittleEndian.Uint16(b)
	//h
	b = make([]byte, 2, 2)
	if _, e = f.Read(b); e != nil {
		return nil, e
	}
	height := binary.LittleEndian.Uint16(b)

	//s
	b = make([]byte, 4, 4)
	if _, e = f.Read(b); e != nil {
		return nil, e
	}
	sign := binary.LittleEndian.Uint32(b)

	//b
	offset := 14 + 40 + 1024
	size := int(width)*int(height) + offset
	b = make([]byte, size, size)
	if _, e = f.Read(b[offset:]); e != nil {
		return nil, e
	}
	//b palette
	copy(b[54:], Palette[:])
	//b 14
	header := newBmp256Header(uint32(size))
	copy(b, header.ToBinary())
	//b 40
	info := newBmp256Info(uint32(width), uint32(height))
	copy(b[14:], info.ToBinary())

	img := &Bmp{Width: width,
		Height: height,
		Sign:   sign,
		Data:   b,
	}
	return img, nil
}
func (w *wixSource) createImg565(i int) (Img, error) {
	pos, e := w.GetPos(i)
	if e != nil {
		return nil, e
	}
	f := w.dataFile
	if _, e = f.Seek(pos, os.SEEK_SET); e != nil {
		return nil, e
	}

	//w
	b := make([]byte, 2, 2)
	if _, e = f.Read(b); e != nil {
		return nil, e
	}
	width := binary.LittleEndian.Uint16(b)
	//h
	b = make([]byte, 2, 2)
	if _, e = f.Read(b); e != nil {
		return nil, e
	}
	height := binary.LittleEndian.Uint16(b)

	//s
	b = make([]byte, 4, 4)
	if _, e = f.Read(b); e != nil {
		return nil, e
	}
	sign := binary.LittleEndian.Uint32(b)

	//b
	offset := 14 + 40 + 12
	size := int(width)*int(height)*2 + offset
	b = make([]byte, size, size)
	if _, e = f.Read(b[offset:]); e != nil {
		return nil, e
	}

	//b 14
	header := newBmp256Header(uint32(size))
	copy(b, header.ToBinary())
	//b 40
	info := newBmp256Info(uint32(width), uint32(height))
	copy(b[14:], info.ToBinary())
	//b12
	copy(b[14+40:], Palette565[:])

	img := &Bmp{Width: width,
		Height: height,
		Sign:   sign,
		Data:   b,
	}
	return img, nil
}
func (w *wixSource) ClearData(i int) {
	if i < 0 || i >= len(w.imgsPos) {
		return
	}
	if w.imgsCache[i] == nil {
		return
	}

	w.imgsCache[i].Destory()
	w.imgsCache[i] = nil
}
func (w *wixSource) Destory() {
	for _, img := range w.imgsCache {
		if img != nil {
			img.Destory()
		}
	}

	if w.dataFile != nil {
		w.dataFile.Close()
	}

}

type wixParser struct {
}

func (w *wixParser) NewSource(t int, xPath, lPath string) (Source, error) {
	//讀取索引 檔案
	f, e := os.Open(xPath)
	if e != nil {
		return nil, e
	}
	defer f.Close()
	info, e := f.Stat()
	if e != nil {
		return nil, e
	}
	size := info.Size()
	if size < 48 {
		return nil, fmt.Errorf("it's not a wix file (%v)", xPath)
	}

	b := make([]byte, 4)
	if _, e = f.ReadAt(b, 44); e != nil {
		return nil, e
	}
	n := binary.LittleEndian.Uint32(b)

	if size < int64(48+n*4) {
		return nil, fmt.Errorf("it's not a wix file (%v)", xPath)
	}
	imgsPos := make([]uint32, n, n)
	if _, e = f.Seek(48, os.SEEK_SET); e != nil {
		return nil, e
	}
	imgsCache := make([]Img, n, n)

	for i := uint32(0); i < n; i++ {
		b := make([]byte, 4)
		if _, e = f.Read(b); e != nil {
			return nil, e
		}
		imgsPos[i] = binary.LittleEndian.Uint32(b)
	}

	//讀取 數據檔案
	dataFile, header, e := newWilHeader(lPath)
	if e != nil {
		return nil, e
	}

	s := &wixSource{imgsPos: imgsPos,
		imgsCache: imgsCache,
		dataFile:  dataFile,
		header:    header,
	}
	return s, nil
}

//加載並驗證 wil header
func newWilHeader(path string) (*os.File, *wilHeader, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, nil, e
	}
	info, e := f.Stat()
	if e != nil {
		f.Close()
		return nil, nil, e
	}
	if info.Size() < 44+2*4 {
		f.Close()
		return nil, nil, fmt.Errorf("it's not a wil file (%v)", path)
	}
	header := &wilHeader{}
	b := make([]byte, 4, 4)
	_, e = f.ReadAt(b, 44)
	if e != nil {
		f.Close()
		return nil, nil, e
	}
	header.Size = binary.LittleEndian.Uint32(b)

	b = make([]byte, 4, 4)
	_, e = f.ReadAt(b, 48)
	if e != nil {
		f.Close()
		return nil, nil, e
	}
	header.ColorDepth = binary.LittleEndian.Uint32(b)

	if header.ColorDepth != COLOR_DEPTH_16_565 &&
		header.ColorDepth != COLOR_DEPTH_256 {
		f.Close()
		return nil, nil, fmt.Errorf("wil color depth not found (%v) (%v)", path, header.ColorDepth)
	}
	return f, header, nil
}
