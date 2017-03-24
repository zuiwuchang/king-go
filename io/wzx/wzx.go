package wzx

import (
	"bytes"
	"compress/zlib"
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
)

type wzxSource struct {
	wixSource
}

func (w *wzxSource) GetType() int {
	return TYPE_WIX
}
func (w *wzxSource) GetImg(i int) (Img, error) {
	if i < 0 || i >= len(w.imgsPos) {
		return nil, fmt.Errorf("index error (%v)", i)
	}
	if w.imgsCache[i] == nil {
		return w.createImg(i)
	}
	return w.imgsCache[i], nil
}
func (w *wzxSource) createImg(i int) (Img, error) {
	pos, e := w.GetPos(i)
	if e != nil {
		return nil, e
	}
	f := w.dataFile
	if _, e = f.Seek(pos, os.SEEK_SET); e != nil {
		return nil, e
	}

	//depth
	b := make([]byte, 2, 2)
	if _, e = f.Read(b); e != nil {
		return nil, e
	}
	depth := binary.LittleEndian.Uint16(b)
	if depth == DEPTH_256 {
		if img, e := w.createImg256(pos); e != nil {
			return nil, e
		} else {
			w.imgsCache[i] = img
			return img, nil
		}
	} else if depth == DEPTH_16_565 {
		if img, e := w.createImg565(pos); e != nil {
			return nil, e
		} else {
			w.imgsCache[i] = img
			return img, nil
		}
	}
	return nil, fmt.Errorf("wzl not found depth (%v)", depth)

}
func (w *wzxSource) createImg256(pos int64) (Img, error) {
	f := w.dataFile
	if _, e := f.Seek(pos+4, os.SEEK_SET); e != nil {
		return nil, e
	}

	//w
	b := make([]byte, 2, 2)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	width := binary.LittleEndian.Uint16(b)
	//h
	b = make([]byte, 2, 2)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	height := binary.LittleEndian.Uint16(b)

	//s
	b = make([]byte, 4, 4)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	sign := binary.LittleEndian.Uint32(b)

	//size
	b = make([]byte, 4, 4)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	size := binary.LittleEndian.Uint32(b)

	//zlib
	b = make([]byte, size, size)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	var buf bytes.Buffer
	if _, e := buf.Write(b); e != nil {
		return nil, e
	}
	r, e := zlib.NewReader(&buf)
	if e != nil {
		return nil, e
	}
	defer r.Close()
	data, e := ioutil.ReadAll(r)

	if e != nil {
		return nil, e
	}

	//b
	offset := 14 + 40 + 1024
	size = uint32(len(data) + offset)
	b = make([]byte, size, size)
	copy(b[offset:], data)

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
func (w *wzxSource) createImg565(pos int64) (Img, error) {
	f := w.dataFile
	if _, e := f.Seek(pos+4, os.SEEK_SET); e != nil {
		return nil, e
	}

	//w
	b := make([]byte, 2, 2)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	width := binary.LittleEndian.Uint16(b)
	//h
	b = make([]byte, 2, 2)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	height := binary.LittleEndian.Uint16(b)

	//s
	b = make([]byte, 4, 4)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	sign := binary.LittleEndian.Uint32(b)

	//size
	b = make([]byte, 4, 4)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	size := binary.LittleEndian.Uint32(b)

	//zlib
	b = make([]byte, size, size)
	if _, e := f.Read(b); e != nil {
		return nil, e
	}
	var buf bytes.Buffer
	if _, e := buf.Write(b); e != nil {
		return nil, e
	}
	r, e := zlib.NewReader(&buf)
	if e != nil {
		return nil, e
	}
	defer r.Close()
	data, e := ioutil.ReadAll(r)
	if e != nil {
		return nil, e
	}

	//b
	offset := 14 + 40 + 12
	dataSize := len(data)
	size = uint32(dataSize + offset)
	b = make([]byte, size, size)
	copy(b[offset:], data)

	//b 14
	header := newBmp565Header(uint32(size))
	copy(b, header.ToBinary())
	//b 40
	info := newBmp565Info(uint32(width), uint32(height), uint32(dataSize))
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

type wzxParser struct {
}

func (w *wzxParser) NewSource(t int, xPath, lPath string) (Source, error) {
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
	dataFile, header, e := newWzlHeader(lPath)
	if e != nil {
		return nil, e
	}

	s := &wzxSource{wixSource{imgsPos: imgsPos,
		imgsCache: imgsCache,
		dataFile:  dataFile,
		header:    header,
	}}
	return s, nil
}

//加載並驗證 wzl header
func newWzlHeader(path string) (*os.File, *wilHeader, error) {
	f, e := os.Open(path)
	if e != nil {
		return nil, nil, e
	}
	info, e := f.Stat()
	if e != nil {
		f.Close()
		return nil, nil, e
	}
	if info.Size() < 44+4 {
		f.Close()
		return nil, nil, fmt.Errorf("it's not a wzl file (%v)", path)
	}
	header := &wilHeader{}
	b := make([]byte, 4, 4)
	_, e = f.ReadAt(b, 44)
	if e != nil {
		f.Close()
		return nil, nil, e
	}
	header.Size = binary.LittleEndian.Uint32(b)

	return f, header, nil
}
