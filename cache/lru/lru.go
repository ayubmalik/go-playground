package lru

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
	val := l.entries[key]
	l.moveToFront(key)
	return val
}

func (l *LRU) Set(key, val string) {
	if len(l.entries) == l.max { // TODO: bug
		delete(l.entries, l.order[l.max])
		l.order = l.order[1:]
	}
	l.entries[key] = val
	l.order = append(l.order, key)
}

func (l *LRU) Len() int {
	return len(l.entries)
}

func (l *LRU) moveToFront(key string) {
	if l.order[0] == key {
		return
	}

	prev := key
	for i, elem := range l.order {
		switch {
		case i == 0:
			l.order[0] = key
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
