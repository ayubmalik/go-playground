package lru

import "fmt"

func New(max int) Cache {
	return &LRU{entries: make(map[string]string), max: max}
}

type Cache interface {
	Get(key string) string

	Set(key, val string)

	Len() int
}

type LRU struct {
	max     int
	entries map[string]string
	order   []string
}

func (l *LRU) Get(key string) string {
	return l.entries[key]
}

func (l *LRU) Set(key, val string) {
	if len(l.entries) == l.max {
		delete(l.entries, l.order[0])
		l.order = l.order[1:]
	}
	l.entries[key] = val
	l.order = append(l.order, key)
}

func (l *LRU) Len() int {
	return len(l.entries)
}

func MoveToFront(needle string, haystack []string) []string {
	if len(haystack) != 0 && haystack[0] == needle {
		return haystack
	}
	prev := needle
	for i, elem := range haystack {
		fmt.Println(i, prev, haystack)
		switch {
		case i == 0:
			haystack[0] = needle
			prev = elem
		case elem == needle:
			fmt.Println("cya")
			haystack[i] = prev
			return haystack
		default:
			haystack[i] = prev
			prev = elem
		}
	}

	fmt.Printf("XYZ prev %s, haystack: %q\n", prev, haystack)
	return append(haystack, prev)
}
