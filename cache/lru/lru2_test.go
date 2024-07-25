package lru_test

import (
	"testing"

	"github.com/ayubmalik/cache/lru"
)

func TestLRU2(t *testing.T) {
	t.Run("set and get", func(t *testing.T) {
		cache := lru.New2(2)

		cache.Set("foo", "bar")

		got := cache.Get("foo")
		if got != "bar" {
			t.Errorf("got %q, want %q", got, "bar")
		}
	})

	t.Run("len with less than max entries", func(t *testing.T) {
		cache := lru.New2(2)

		cache.Set("foo1", "bar1")

		got := cache.Len()
		if got != 1 {
			t.Errorf("got %d, want %d", got, 1)
		}
	})

	t.Run("len does not exceed max", func(t *testing.T) {
		cache := lru.New2(2)

		cache.Set("foo1", "bar1")
		cache.Set("foo2", "bar2")
		cache.Set("foo3", "bar3")

		gotLen := cache.Len()
		if gotLen != 2 {
			t.Errorf("got len %d, want %d", gotLen, 2)
		}

		got := cache.Get("foo3")
		if got != "bar3" {
			t.Errorf("got %q, want %q", got, "bar3")
		}

		got2 := cache.Get("foo1")
		if got2 != "" {
			t.Errorf("got1 %q, want empty string", got2)
		}
	})

	t.Run("evict least recently used no Get", func(t *testing.T) {
		cache := lru.New2(4)

		cache.Set("foo1", "bar1")
		cache.Set("foo2", "bar2")
		cache.Set("foo3", "bar3")
		cache.Set("foo4", "bar4")
		cache.Set("foo5", "bar5")

		if cache.Get("foo5") != "bar5" {
			t.Errorf("got %q, want %q", cache.Get("foo5"), "bar5")
		}

		if cache.Get("foo1") != "" {
			t.Errorf("got %q, want %q", cache.Get("foo3"), "bar3")
		}
	})

	t.Run("evict least recently used with Get", func(t *testing.T) {
		cache := lru.New2(4)

		cache.Set("foo1", "bar1")
		cache.Set("foo2", "bar2")
		cache.Set("foo3", "bar3")
		cache.Set("foo4", "bar4")

		cache.Get("foo1")
		cache.Get("foo2")
		cache.Get("foo3")

		cache.Set("foo5", "bar5")

		got := cache.Get("foo5")
		want := "bar5"
		if got != "bar5" {
			t.Errorf("got %s, want %s", got, want)
		}

		got = cache.Get("foo1")
		want = "bar1"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		got = cache.Get("foo2")
		want = "bar2"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		got = cache.Get("foo3")
		want = "bar3"
		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}

		got = cache.Get("foo4")
		want = ""
		if got != want {
			t.Errorf("no evicted - got %q, want %q", got, want)
		}
	})
}
