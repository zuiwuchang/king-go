package ini

import (
	"testing"
)

const (
	testOkStr = `[Node]
I8=-1
I16=-2
I32=-3
I64=-4
I=-5
U8=1
U16=2
U32=3
U64=4
U=5
Float32=1.1
Float64=1.2
Bool=true
;str test
str=ok
`
)

func TestEncoder(t *testing.T) {
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

	node := Node{-1, -2, -3, -4, -5, //int
		1, 2, 3, 4, 5, //uint
		1.1, 1.2, //float
		true, //bool
		"ok", //str
		nil,
	}
	b, e := Marshal(node)
	if e != nil {
		t.Fatal(e)
	}
	s := string(b)

	b1, e := Marshal(&node)
	if e != nil {
		t.Fatal(e)
	}
	s1 := string(b1)

	if s != s1 {
		t.Fatal("s != s1")
	}

	if s != testOkStr {
		t.Fatal("s != testOkStr\n")
	}

}
func TestEncoderPtr(t *testing.T) {
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

	var i8 int8 = -1
	var i16 int16 = -2
	var i32 int32 = -3
	var i64 int64 = -4
	var i int = -5

	var u8 uint8 = 1
	var u16 uint16 = 2
	var u32 uint32 = 3
	var u64 uint64 = 4
	var u uint = 5

	var f32 float32 = 1.1
	var f64 float64 = 1.2

	ok := true
	str := "ok"

	node := Node{&i8, &i16, &i32, &i64, &i, //int
		&u8, &u16, &u32, &u64, &u, //uint
		&f32, &f64, //float
		&ok,  //bool
		&str, //str
	}

	b, e := Marshal(node)
	if e != nil {
		t.Fatal(e)
	}
	s := string(b)

	b1, e := Marshal(&node)
	if e != nil {
		t.Fatal(e)
	}
	s1 := string(b1)

	if s != s1 {
		t.Fatal("s != s1")
	}

	if s != testOkStr {
		t.Fatal("s != testOkStr\n")
	}
}
