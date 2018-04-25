package queue

import (
	"testing"
)

func TestStatic(t *testing.T) {
	iq, e := NewStatic(3)
	if e != nil {
		t.Fatal(e)
	}
	q := iq.(*staticQueue)
	if q.Len() != 0 {
		t.Fatal("bad len")
	}
	if q.Cap() != 3 {
		t.Fatal("bad cap")
	}

	// 0 1 2
	for i := 0; i < 3; i++ {
		e = q.PushBack(i)
		if e != nil {
			t.Fatal(e)
		}
	}
	e = q.PushBack(3)
	if e == nil {
		t.Fatal("bad overflow")
	}
	for i := 0; i < 3; i++ {
		if q.data[i] != i {
			t.Fatal("bad val", q.data[i], i)
		}
	}

	// 1 2
	iv, e := q.PopFront()
	if e != nil {
		t.Fatal(e)
	}
	v := iv.(int)
	if v != 0 {
		t.Fatal("bad val", v, 0)
	}

	// q 1 2 3
	// d 3 1 2
	e = q.PushBack(3)
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 3; i++ {
		if i == 0 {
			if q.data[i] != 3 {
				t.Fatal("bad val", q.data[i], 3)
			}
		} else {
			if q.data[i] != i {
				t.Fatal("bad val", q.data[i], i)
			}
		}
	}

	// q 1 2
	// d x 1 2
	iv, e = q.PopBack()
	if e != nil {
		t.Fatal(e)
	}
	v = iv.(int)
	if v != 3 {
		t.Fatal("bad val", v, 3)
	}

	// q 1
	// d x 1 x
	iv, e = q.PopBack()
	if e != nil {
		t.Fatal(e)
	}
	v = iv.(int)
	if v != 2 {
		t.Fatal("bad val", v, 2)
	}

	// q 2 0 1
	// d 0 1 2
	e = q.PushFront(0)
	if e != nil {
		t.Fatal(e)
	}
	e = q.PushFront(2)
	if e != nil {
		t.Fatal(e)
	}
	for i := 0; i < 3; i++ {
		if q.data[i] != i {
			t.Fatal("bad val", q.data[i], i)
		}
	}

	// [2] 0 1
	iv, e = q.PopFront()
	if e != nil {
		t.Fatal(e)
	}
	v = iv.(int)
	if v != 2 {
		t.Fatal("bad val", v, 2)
	}
	// [0] 1
	iv, e = q.PopFront()
	if e != nil {
		t.Fatal(e)
	}
	v = iv.(int)
	if v != 0 {
		t.Fatal("bad val", v, 0)
	}

	// [1]
	iv, e = q.PopFront()
	if e != nil {
		t.Fatal(e)
	}
	v = iv.(int)
	if v != 1 {
		t.Fatal("bad val", v, 1)
	}

	//
	if q.Len() != 0 {
		t.Fatal("bad len")
	}
	if q.Cap() != 3 {
		t.Fatal("bad cap")
	}
}
