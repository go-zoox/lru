package lru

import (
	"container/list"
	"sync"
)

type AtomicInt int64

func (a AtomicInt) Inc() {
	a += 1
}

type LRU struct {
	mutex    sync.RWMutex
	cache    map[string]*list.Element
	list     *list.List
	Capacity int
	Hits     AtomicInt
	Gets     AtomicInt
}

type entry struct {
	key   string
	value interface{}
}

func New(capacity int) *LRU {
	return &LRU{
		Capacity: capacity,
		cache:    make(map[string]*list.Element),
		list:     list.New(),
	}
}

func (l *LRU) Get(key string) (interface{}, bool) {
	l.mutex.RLock()
	defer l.mutex.RUnlock()
	l.Gets.Inc()

	if elem, ok := l.cache[key]; ok {
		l.Hits.Inc()
		l.list.MoveToFront(elem)
		return elem.Value.(*entry).value, true
	}

	return nil, false
}

func (l *LRU) Set(key string, value interface{}) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.cache == nil {
		l.cache = make(map[string]*list.Element)
		l.list = list.New()
	}

	if el, ok := l.cache[key]; ok {
		l.list.MoveToFront(el)
		el.Value.(*entry).value = value
		return
	}

	ele := l.list.PushFront(&entry{key, value})
	l.cache[key] = ele
	if l.Capacity != 0 && l.list.Len() > l.Capacity {
		l.removeOldest()
	}
}

func (l *LRU) Delete(key string) {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	if l.cache == nil {
		return
	}

	if elem, ok := l.cache[key]; ok {
		l.list.Remove(elem)
		delete(l.cache, key)
	}
}

func (l *LRU) removeOldest() {
	if l.cache == nil {
		return
	}

	elem := l.list.Back()
	if elem != nil {
		l.list.Remove(elem)
		kv := elem.Value.(*entry)
		delete(l.cache, kv.key)
	}
}

func (l *LRU) Len() int {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	return l.list.Len()
}

func (l *LRU) Keys() []string {
	l.mutex.RLock()
	defer l.mutex.RUnlock()

	keys := make([]string, l.list.Len())
	i := 0
	for el := l.list.Front(); el != nil; el = el.Next() {
		keys[i] = el.Value.(*entry).key
		i += 1
	}

	return keys
}

func (l *LRU) Clear() {
	l.mutex.Lock()
	defer l.mutex.Unlock()

	l.cache = make(map[string]*list.Element)
	l.list = list.New()
}
