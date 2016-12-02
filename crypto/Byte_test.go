package crypto

import (
	"testing"
)

func TestByte(t *testing.T) {
	b := Byte{Data: 0x80}
	for i := 0; i < 7; i++ {
		if b.IsBitSet(i) {
			t.Fatal("IsBitSet true", i)
		}
	}
	if !b.IsBitSet(7) {
		t.Fatal("IsBitSet false")
	}

	b.SwapBit(0, 7)
	for i := 1; i < 8; i++ {
		if b.IsBitSet(i) {
			t.Fatal("IsBitSet true", i)
		}
	}
	if !b.IsBitSet(0) {
		t.Fatal("IsBitSet false")
	}
}
