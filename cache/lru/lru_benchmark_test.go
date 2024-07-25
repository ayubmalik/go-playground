package lru_test

import (
	"fmt"
	"testing"

	"github.com/ayubmalik/cache/lru"
)

func BenchmarkLRU(b *testing.B) {
	capacity := 1_000_000

	b.Run("LRU1", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			cache := lru.New(capacity)
			for i := 0; i < capacity; i++ {
				key := fmt.Sprintf("key%0d", i)
				val := fmt.Sprintf("val%0d", i)
				cache.Set(key, val)
			}
		}
	})

	b.Run("LRU2", func(b *testing.B) {
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			cache := lru.New2(capacity)
			for i := 0; i < capacity; i++ {
				key := fmt.Sprintf("key%0d", i)
				val := fmt.Sprintf("val%0d", i)
				cache.Set(key, val)
			}
		}
	})
}
