package set

import (
	"testing"
)

func TestSetInt64(t *testing.T) {
	Count := 10
	set := NewInt64()
	for i := 0; i < Count; i++ {
		set.Insert(int64(i))
	}
	if set.Len() != Count {
		t.Fatal("bad len")
	}

	for i := 0; i < Count; i++ {
		if !set.Ok(int64(i)) {
			t.Fatal("bad Ok", i)
		}
	}
	for i := 0; i < Count; i++ {
		if i%2 == 0 {
			set.Remove(int64(i))
		}
	}
	for i := 0; i < Count; i++ {
		if i%2 == 0 {
			if set.Ok(int64(i)) {
				t.Fatal("bad Ok", i)
			}
		} else {
			if !set.Ok(int64(i)) {
				t.Fatal("bad Ok", i)
			}
		}
	}
}
