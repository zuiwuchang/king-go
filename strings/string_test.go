package strings

import (
	"testing"
)

func TestSubString(t *testing.T) {
	var v string
	str := "草泥馬b的測試"
	if Left(str, 0) != "" || Left(str, -1) != "" {
		t.Fatal("left 0")
	}
	v = Left(str, 1)
	if v != "草" {
		t.Fatal("left 1", v)
	}
	v = Left(str, 6)
	if v != "草泥馬b的測" {
		t.Fatal("left 6", v)
	}
	v = Left(str, 7)
	if v != str {
		t.Fatal("left 7", v)
	}
	v = Left(str, 10)
	if v != str {
		t.Fatal("left 10", v)
	}

	if Right(str, 0) != "" || Right(str, -1) != "" {
		t.Fatal("right 0")
	}
	v = Right(str, 1)
	if v != "試" {
		t.Fatal("right 1", v)
	}
	v = Right(str, 6)
	if v != "泥馬b的測試" {
		t.Fatal("right 6", v)
	}
	v = Right(str, 7)
	if v != str {
		t.Fatal("right 7", v)
	}
	v = Right(str, 10)
	if v != str {
		t.Fatal("right 7", v)
	}

	v = Sub(str)
	if v != str {
		t.Fatal("Sub", v)
	}

	v = Sub(str, 1)
	if v != "泥馬b的測試" {
		t.Fatal("Sub 1", v)
	}
	v = Sub(str, 100)
	if v != "" {
		t.Fatal("Sub 100", v)
	}
	v = Sub(str, 1, 2)
	if v != "泥馬" {
		t.Fatal("Sub 1 2", v)
	}
}
