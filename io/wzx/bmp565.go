package wzx

var Palette565 = [12]byte{0x0, 0xF8, 0x0, 0x0,
	0xE0, 0x07, 0x0, 0x0,
	0x1F, 0x0, 0x0, 0x0,
}

func newBmp565Header(size uint32) *BmpHeader {
	header := BmpHeader{Size: size, Pos: 14 + 40 + 12}
	header.Flag[0] = 0x42
	header.Flag[1] = 0x4D
	return &header
}
func newBmp565Info(w, h, size uint32) *BmpInfo {
	info := BmpInfo{Size: 40,
		Width:       w,
		Height:      h,
		Planes:      1,
		BitCount:    16,
		Compression: 3,
		SizeImage:   size,
	}
	return &info
}
