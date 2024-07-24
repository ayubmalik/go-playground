package lru_test

import (
	"fmt"
	"github.com/ayubmalik/cache/lru"
	"testing"
)

func TestLRU(t *testing.T) {

	t.Run("set and get", func(t *testing.T) {
		cache := lru.New(2)

		cache.Set("foo", "bar")

		got := cache.Get("foo")
		if got != "bar" {
			t.Errorf("got %q, want %q", got, "bar")
		}
	})

	t.Run("len with less than max entries", func(t *testing.T) {
		cache := lru.New(2)

		cache.Set("foo1", "bar1")

		got := cache.Len()
		if got != 1 {
			t.Errorf("got %d, want %d", got, 1)
		}
	})

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

func TestSpike(t *testing.T) {
	data := []string{"foo", "bar", "cheese", "tea", "banana"}

	bubble := func(key string, items []string) {

		var i int
		for i = 0; i < len(items); i++ {
			if items[i] == key {
				break
			}
		}
		for j := i; j > 0; j-- {
			data[j-1], data[j] = data[j], data[j-1]
		}
		fmt.Printf("data: %q\n", data)
	}

	bubble("cheese", data)
	bubble("foo", data)
	bubble("banana", data)
}

func TestMoveToFront(t *testing.T) {
	haystack := []string{"a", "b", "c", "d", "e"} // [a b c d e]
	haystack = lru.MoveToFront("d", haystack)     // [c a b d e]
	fmt.Printf("haystackZZ: %q\n", haystack)
}
