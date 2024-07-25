package lru

import "container/list"

func New2(max int) LRU2 {
	return LRU2{list.New(), max}
}

type LRU2 struct {
	entries *list.List
	max     int
}

type entry struct {
	k string
	v string
}

func (l LRU2) Set(key, val string) {
	if l.entries.Len() == l.max {
		l.entries.Remove(l.entries.Back())
	}
	l.entries.PushFront(&entry{key, val})
}

func (l LRU2) Get(key string) string {
	for elem := l.entries.Front(); elem != nil; elem = elem.Next() {
		e := elem.Value.(*entry)
		if e.k == key {
			l.entries.MoveToFront(elem)
			return e.v
		}
	}
	return ""
}

func (l LRU2) Len() int {
	return l.entries.Len()
}
