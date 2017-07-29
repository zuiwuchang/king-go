package ini

import (
	"testing"
)

func TestDecoder(t *testing.T) {
	type Node struct {
		I8  int8
		I16 int16
		I32 int32
		I64 int64
		I   int

		U8  uint8
		U16 uint16
		U32 uint32
		U64 uint64
		U   uint

		Float32 float32
		Float64 float64

		Bool bool
		Str  string `ini:"str" dec:"str test"`
		List []int
	}
	node := &Node{}
	e := Unmarshal([]byte(testOkStr), node)
	if e != nil {
		t.Fatal(e)
	}

	if node.I8 != -1 {
		t.Fatal("I8")
	}
	if node.I16 != -2 {
		t.Fatal("I16")
	}
	if node.I32 != -3 {
		t.Fatal("I32")
	}
	if node.I64 != -4 {
		t.Fatal("I64")
	}
	if node.I != -5 {
		t.Fatal("I")
	}

	if node.U8 != 1 {
		t.Fatal("U8")
	}
	if node.U16 != 2 {
		t.Fatal("U16")
	}
	if node.U32 != 3 {
		t.Fatal("U32")
	}
	if node.U64 != 4 {
		t.Fatal("U64")
	}
	if node.U != 5 {
		t.Fatal("U")
	}

	if node.Float32 != 1.1 {
		t.Fatal("Float32")
	}
	if node.Float64 != 1.2 {
		t.Fatal("Float64")
	}

	if !node.Bool {
		t.Fatal("Bool")
	}

	if node.Str != "ok" {
		t.Fatal("Str")
	}
}
func TestDecoderPtr(t *testing.T) {
	type Node struct {
		I8  *int8
		I16 *int16
		I32 *int32
		I64 *int64
		I   *int

		U8  *uint8
		U16 *uint16
		U32 *uint32
		U64 *uint64
		U   *uint

		Float32 *float32
		Float64 *float64

		Bool *bool
		Str  *string `ini:"str" dec:"str test"`
	}
	node := &Node{}
	e := Unmarshal([]byte(testOkStr), node)
	if e != nil {
		t.Fatal(e)
	}

	if *node.I8 != -1 {
		t.Fatal("I8")
	}
	if *node.I16 != -2 {
		t.Fatal("I16")
	}
	if *node.I32 != -3 {
		t.Fatal("I32")
	}
	if *node.I64 != -4 {
		t.Fatal("I64")
	}
	if *node.I != -5 {
		t.Fatal("I")
	}

	if *node.U8 != 1 {
		t.Fatal("U8")
	}
	if *node.U16 != 2 {
		t.Fatal("U16")
	}
	if *node.U32 != 3 {
		t.Fatal("U32")
	}
	if *node.U64 != 4 {
		t.Fatal("U64")
	}
	if *node.U != 5 {
		t.Fatal("U")
	}

	if *node.Float32 != 1.1 {
		t.Fatal("Float32")
	}
	if *node.Float64 != 1.2 {
		t.Fatal("Float64")
	}

	if !*node.Bool {
		t.Fatal("Bool")
	}

	if *node.Str != "ok" {
		t.Fatal("Str")
	}
}
