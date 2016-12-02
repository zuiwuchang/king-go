package crypto

//提供了對 byte 操作 的封裝
type Byte struct {
	Data byte
}

//
func (b *Byte) getBitFlag(bit int) byte {
	switch bit {
	case 0:
		return 0x01
	case 1:
		return 0x02
	case 2:
		return 0x04
	case 3:
		return 0x08
	case 4:
		return 0x10
	case 5:
		return 0x20
	case 6:
		return 0x40
	case 7:
		return 0x80
	}
	return 0
}

//返回 指定 bit 是否被 設置
func (b *Byte) IsBitSet(bit int) bool {
	flag := b.getBitFlag(bit)
	return b.Data&flag != 0
}

//設置 指定位 爲 1 or 0
func (b *Byte) SetBit(bit int, set bool) {
	flag := b.getBitFlag(bit)
	if set {
		b.Data |= flag
	} else {
		b.Data &= (^flag)
	}
}

//交換 指定的 兩個位的 值
func (b *Byte) SwapBit(l, r int) {
	lSet := b.IsBitSet(l)
	rSet := b.IsBitSet(r)

	b.SetBit(l, rSet)
	b.SetBit(r, lSet)
}

//異或
func (b *Byte) Xor(r byte) {
	b.Data ^= r
}

//或
func (b *Byte) Or(r byte) {
	b.Data |= r
}

//與
func (b *Byte) And(r byte) {
	b.Data &= r
}

//取反
func (b *Byte) Not() {
	b.Data = ^b.Data
}

//左移 n 位
func (b *Byte) Shl(n uint) {
	b.Data <<= n
}
func (b *Byte) ShlLoop(n uint) {
	n %= 8
	if n == 0 {
		return
	}

	h := b.Data << n
	l := b.Data >> (8 - n)
	b.Data = h | l
}

//右移 n 位
func (b *Byte) Shr(n uint) {
	b.Data >>= n
}
func (b *Byte) ShrLoop(n uint) {
	n %= 8
	if n == 0 {
		return
	}

	l := b.Data >> n
	h := b.Data << (8 - n)
	b.Data = h | l
}
