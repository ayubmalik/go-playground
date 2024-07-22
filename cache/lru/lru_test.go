package lru_test

import (
	"github.com/ayubmalik/cache/lru"
	"testing"
)

func TestLRU(t *testing.T) {

	//t.Run("set and get", func(t *testing.T) {
	//	cache := lru.New(2)
	//
	//	cache.Set("foo", "bar")
	//
	//	got := cache.Get("foo")
	//	if got != "bar" {
	//		t.Errorf("got %q, want %q", got, "bar")
	//	}
	//})
	//
	//t.Run("len with less than max entries", func(t *testing.T) {
	//	cache := lru.New(2)
	//
	//	cache.Set("foo1", "bar1")
	//
	//	got := cache.Len()
	//	if got != 1 {
	//		t.Errorf("got %d, want %d", got, 1)
	//	}
	//})

	t.Run("len does not exceed capacity", func(t *testing.T) {
		cache := lru.New(2)

		cache.Set("foo1", "bar1")
		cache.Set("foo2", "bar2")
		cache.Set("foo3", "bar3")

		got := cache.Get("foo3")
		if got != "bar3" {
			t.Errorf("got %q, want %q", got, "bar3")
		}

		gotLen := cache.Len()
		if gotLen != 2 {
			t.Errorf("got len %d, want %d", gotLen, 2)
		}

		got2 := cache.Get("foo1")
		if got2 != "" {
			t.Errorf("got1 %q, want empty string", got2)
		}
	})

}

//func TestSet(t *testing.T) {
//	lru := New(1)
//
//	lru.Set("foo", "bar")
//
//	got := lru.Get("foo")
//	if got != "bar" {
//		t.Errorf("got %s wanted 0", got)
//	}
//
//}
//
//func TestCapacity(t *testing.T) {
//	lru := New(2)
//
//	lru.Set("foo1", "bar1")
//	lru.Set("foo2", "bar1")
//	lru.Set("foo3", "bar1")
//
//	l := lru.Len()
//	if l != 2 {
//		t.Errorf("got %d want 2", l)
//	}
//
//}
