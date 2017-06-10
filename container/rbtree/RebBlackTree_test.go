package rbtree

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("can not create")
	}
	count := 1000
	for i := 0; i < count; i++ {
		if !m.Insert(IKeyInt(i), i) {
			t.Fatal("insert false")
		}
	}
	for i := 0; i < count; i++ {
		if m.Insert(IKeyInt(i), i+10) {
			t.Fatal("insert repeat")
		}
	}
	if m.length != count {
		t.Fatal("insert mis node")
	}
	i := 0
	m.Do(func(k IKey, v IValue) bool {
		if int(k.(IKeyInt)) != i {
			t.Fatal("k error")
		}
		if v.(int) != i+10 {
			t.Fatal("v error", v)
		}
		i++
		return true
	})
	if i != count {
		t.Fatalf("do need %v not %v", count, i)
	}

	fmt.Println("ok")
}
