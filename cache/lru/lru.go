package lru

func New(max int) Cache {
	return &LRU{entries: make(map[string]string, max), max: max}
}

type Cache interface {
	Get(key string) string

	Set(key, val string)

	Len() int
}

type LRU struct {
	entries map[string]string
	order   []string
	max     int
}

func (l *LRU) Get(key string) string {
	val := l.entries[key]
	l.move(key)
	return val
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

func (l *LRU) move(key string) {
	last := len(l.order) - 1
	if l.order[last] == key {
		return
	}

	prev := key
	for i := last; i >= 0; i-- {
		elem := l.order[i]
		switch {
		case i == last:
			l.order[last] = key
			prev = elem
		case elem == key:
			l.order[i] = prev
			return
		default:
			l.order[i] = prev
			prev = elem
		}
	}
	// should not be here
}
