package rbtree

import (
	"math/rand"
	"testing"
	"time"
)

func TestBase(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("can not create")
	}
	//insert
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
	ele := m.Min()
	if int(ele.Key().(IKeyInt)) != 0 {
		t.Fatal("min error")
	}
	ele = m.Max()
	if int(ele.Key().(IKeyInt)) != count-1 {
		t.Fatal("max error")
	}

	i := 0
	m.Do(func(ele IElement) bool {
		k := ele.Key()
		v := ele.Value()
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
	i = 0
	m.DoTree(m.root, func(ele IElement) bool {
		k := ele.Key()
		v := ele.Value()
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

	i = 0
	m.DoReverse(func(ele IElement) bool {
		k := ele.Key()
		v := ele.Value()

		if int(k.(IKeyInt)) != (count - i - 1) {
			t.Fatal("k error")
		}
		if v.(int) != (count-i-1)+10 {
			t.Fatal("v error", v)
		}
		i++
		return true
	})
	if i != count {
		t.Fatalf("do need %v not %v", count, i)
	}
	i = 0
	m.DoTreeReverse(m.root, func(ele IElement) bool {
		k := ele.Key()
		v := ele.Value()

		if int(k.(IKeyInt)) != (count - i - 1) {
			t.Fatal("k error")
		}
		if v.(int) != (count-i-1)+10 {
			t.Fatal("v error", v)
		}
		i++
		return true
	})
	if i != count {
		t.Fatalf("do need %v not %v", count, i)
	}
	//is read blak tree
	inorderWalk(t, m.root)

	//find

	//reset
	m.Reset()
	if m.length != 0 {
		t.Fatal("Reset error")
	}
	ele = m.Min()
	if ele != nil {
		t.Fatal("min error")
	}
	ele = m.Max()
	if ele != nil {
		t.Fatal("max error")
	}
}

//驗證 紅黑樹 性質
func inorderWalk(t *testing.T, x *_Element) {
	if x != nil {
		inorderWalk(t, x.L)
		if x != _ElementNil {
			if x.Red { //驗證性質
				if x.P.Red || x.L.Red || x.P.Red {
					t.Fatal("bad rbt red")
				}
			}
			if x.L == nil || x.R == nil {
				t.Fatal("bad rbt lr")
			}
		}
		inorderWalk(t, x.R)
	}
}
func TestErase(t *testing.T) {
	m := New()
	count := 1000
	for i := 0; i < count; i++ {
		m.Insert(IKeyInt(i), i+10)
	}
	//get
	for i := 0; i < count; i++ {
		key := +IKeyInt(i)
		ele := m.Get(key)
		v := m.GetValue(key)
		if ele.Key() != key {
			t.Fatal("bad get")
		}
		if v.(int) != i+10 || v != ele.Value() {
			t.Fatal("bad val")
		}
	}

	m.Erase(IKeyInt(count))
	if m.length != count {
		t.Fatal("erase bad key")
	}
	key := IKeyInt(10)
	m.Erase(key)
	if m.length != count-1 {
		t.Fatal("erase bad key")
	}
	m.Erase(key)
	if m.length != count-1 {
		t.Fatal("erase bad key")
	}
	if m.Get(key) != nil || m.GetValue(key) != nil {
		t.Fatal("get bad key")
	}
	length := m.length
	for i := 0; i < count/2; i++ {
		m.Erase(IKeyInt(i))
		if i != 10 {
			length--
		}
		if m.length != length {
			t.Fatal("for bad erase", m.length, length)
		}
	}
	inorderWalk(t, m.root)
}
func TestRand(t *testing.T) {
	rand.Seed(time.Now().Unix())
	for i := 0; i < 500; i++ {
		rand.Seed(int64(i))

		m := New()
		count := (rand.Int()%5 + 1) + 1024
		length := 0
		for i := 0; i < count; i++ {
			if rand.Int()%2 == 0 {
				//insert
				k := IKeyInt(rand.Int() % 1024 * 100)
				sum := m.Len()
				ele := m.Get(k)
				m.Insert(k, i)
				if ele == nil {
					length++
					if sum+1 != m.Len() {
						t.Fatal("bad insert")
					}
				} else {
					if sum != m.Len() {
						t.Fatal("bad insert")
					}
				}
			} else {
				//delete
				k := IKeyInt(rand.Int() % 1024 * 100)
				sum := m.Len()
				ele := m.Get(k)
				m.Erase(k)
				if ele == nil {
					if sum != m.Len() {
						t.Fatal("bad erase")
					}
				} else {
					length--
					if sum-1 != m.Len() {
						t.Fatal("bad erase")
					}
				}
			}
		}
		if length != m.Len() {
			t.Fatal("bad rand")
		}
		inorderWalk(t, m.root)
	}
}
