package lru

import (
	"testing"
)

func TestResize(t *testing.T) {
	cache := newLRUImpl(3)
	cache.Set(0, 100)
	cache.Set(1, 101)
	cache.Set(2, 102)
	cache.Set(3, 103)

	if cache.Resize(0.67) != 2 ||
		cache.Ok(0) || cache.Ok(1) ||
		!cache.Ok(2) || !cache.Ok(3) {
		t.Fatal("bad Resize to 2 len")
	}
	cache.debugPrint(t)
	if cache.Resize(0.67) != 2 ||
		cache.Ok(0) || cache.Ok(1) ||
		!cache.Ok(2) || !cache.Ok(3) {
		t.Fatal("bad Resize to 2 len")
	}
	cache.debugPrint(t)

	if cache.Resize(0.34) != 1 ||
		cache.Ok(0) || cache.Ok(1) ||
		cache.Ok(2) || !cache.Ok(3) {
		t.Fatal("bad Resize to 1 len")
	}
	cache.debugPrint(t)
	if cache.Resize(0.34) != 1 ||
		cache.Ok(0) || cache.Ok(1) ||
		cache.Ok(2) || !cache.Ok(3) {
		t.Fatal("bad Resize to 1 len")
	}
	cache.debugPrint(t)

	if cache.Resize(0.1) != 0 ||
		cache.Ok(0) || cache.Ok(1) ||
		cache.Ok(2) || cache.Ok(3) {
		t.Fatal("bad Resize to clear")
	}
	cache.debugPrint(t)
	if cache.Resize(0.1) != 0 ||
		cache.Ok(0) || cache.Ok(1) ||
		cache.Ok(2) || cache.Ok(3) {
		t.Fatal("bad Resize to clear")
	}
	cache.debugPrint(t)
}
func TestLRU2(t *testing.T) {
	cache := newLRUImpl(3)
	cache.Set(0, 100)
	cache.debugPrint(t)
	cache.Set(1, 101)
	cache.debugPrint(t)
	cache.Get(0)
	cache.debugPrint(t)
	cache.Set(2, 102)
	cache.debugPrint(t)
	cache.Set(3, 103)
	cache.debugPrint(t)
}
func TestLRU(t *testing.T) {
	cache := newLRUImpl(3)
	cache.Set("kate", "my love")

	if cache.Len() != 1 {
		t.Fatal("bad len")
	}
	if !cache.Ok("kate") || cache.Ok("Leo") {
		t.Fatal("bad ok")
	}
	if "my love" != cache.Get("kate").(string) {
		t.Fatal("bad get")
	}
	if cache.Get("Leo") != nil {
		t.Fatal("bad get")
	}

	cache.Set("kate", "my love")
	if cache.Len() != 1 {
		t.Fatal("bad update")
	}
	cache.Set("cat0", "me0")
	cache.Set("cat1", "me1")
	cache.Set("cat2", "me2")
	if cache.Ok("kate") ||
		cache.Get("cat0") != "me0" ||
		cache.Get("cat1") != "me1" ||
		cache.Get("cat2") != "me2" {

		t.Fatal("bad front")
	}

	cache.Set("kate", "my love")
	cache.Set("cat0", "me0")
	cache.Get("kate")
	cache.Set("cat1", "me1")
	cache.Set("cat2", "me2")
	if cache.Get("kate") != "my love" ||
		cache.Ok("cat0") ||
		cache.Get("cat1") != "me1" ||
		cache.Get("cat2") != "me2" {

		t.Fatal("bad front")
	}
}
