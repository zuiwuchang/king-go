package lru

import (
	"testing"
	"time"
)

func TestLRU(t *testing.T) {
	cache := newLRUImpl(time.Millisecond*100, 3)
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

	time.Sleep(time.Millisecond * 200)
	if cache.Ok("kate") || cache.Get("kate") != nil {
		t.Fatal("bad timeout")
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

	cache.Set("kate", "my love")
	go func() {
		for i := 0; i < 10; i++ {
			time.Sleep(time.Millisecond * 50)
			if cache.Get("kate") != "my love" {
				t.Fatal("bad get", i)
			}
		}
	}()

	time.Sleep(time.Millisecond * 500)
	if cache.Get("kate") != "my love" {
		t.Fatal("bad timeout reset")
	}
}
