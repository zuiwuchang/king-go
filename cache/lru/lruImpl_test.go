package lru

import (
	"testing"
)

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
	cache.Set("cat0", "me0")
	cache.Set("cat1", "me1")
	cache.Set("cat2", "me2")
	if cache.Ok("kate") ||
		cache.Get("cat0") != "me0" ||
		cache.Get("cat1") != "me1" ||
		cache.Get("cat2") != "me2" {

		t.Fatal("bad front")
	}

}
