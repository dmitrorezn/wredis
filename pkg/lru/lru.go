package lru

import (
	"container/list"
	"sync"
)

type Item[K comparable, V any] struct {
	K K
	V V
}

func newHmap[K comparable, V any](size int) hmap[K, V] {
	return make(hmap[K, V], size)
}

func lockUnlock(locker sync.Locker) func() {
	locker.Lock()

	return locker.Unlock
}

type hmap[K comparable, V any] map[K]*Item[K, V]

func (m hmap[K, V]) Put(key K, val *Item[K, V]) {
	m[key] = val
}
func (m hmap[K, V]) Del(key K) {
	delete(m, key)
}
func (m hmap[K, V]) Get(key K) (*Item[K, V], bool) {
	v, ok := m[key]

	return v, ok
}

func (m hmap[K, V]) Exists(key K) bool {
	_, ok := m[key]

	return ok
}

type LRU[K comparable, V any] struct {
	lru *lru[K, V]
	mu  sync.Mutex
}

func (l *LRU[K, V]) Len() int {
	defer lockUnlock(&l.mu)()

	return l.lru.Len()
}

func (l *LRU[K, V]) Cap() int {
	defer lockUnlock(&l.mu)()

	return l.lru.Cap()
}

func (l *LRU[K, V]) Purge() {
	defer lockUnlock(&l.mu)()

	l.lru.Purge()
}

func (l *LRU[K, V]) Keys() []K {
	defer lockUnlock(&l.mu)()

	return l.lru.Keys()
}

func (l *LRU[K, V]) Values() []V {
	defer lockUnlock(&l.mu)()

	return l.lru.Values()
}

func (l *LRU[K, V]) Contains(key K) bool {
	defer lockUnlock(&l.mu)()

	return l.lru.Contains(key)
}

func (l *LRU[K, V]) Add(key K, val V) (evicted bool) {
	defer lockUnlock(&l.mu)()

	return l.lru.Add(key, val)
}

func (l *LRU[K, V]) GetOldest() (key K, val V, ok bool) {
	defer lockUnlock(&l.mu)()

	return l.lru.GetOldest()
}

func (l *LRU[K, V]) RemoveOldest() (K, bool) {
	defer lockUnlock(&l.mu)()

	return l.lru.RemoveOldest()
}

func (l *LRU[K, V]) Remove(key K) bool {
	defer lockUnlock(&l.mu)()

	return l.lru.Remove(key)
}

func (l *LRU[K, V]) Resize(size int) int {
	defer lockUnlock(&l.mu)()

	return l.lru.Resize(size)
}

func (l *LRU[K, V]) Peek(key K) (V, bool) {
	defer lockUnlock(&l.mu)()

	return l.lru.Peek(key)
}

func (l *LRU[K, V]) Get(key K) (V, bool) {
	defer lockUnlock(&l.mu)()

	return l.lru.Get(key)
}

func NewLRU[K comparable, V any](size int, onEvicted func(K, V)) *LRU[K, V] {
	return &LRU[K, V]{
		lru: newLRU(size, onEvicted),
	}
}

type lru[K comparable, V any] struct {
	size int

	items     hmap[K, *list.Element]
	onEvicted func(K, V)

	list *list.List
}

func newLRU[K comparable, V any](size int, onEvicted func(K, V)) *lru[K, V] {
	return &lru[K, V]{
		size:      size,
		items:     newHmap[K, *list.Element](size),
		list:      list.New(),
		onEvicted: onEvicted,
	}
}

func (l *lru[K, V]) Len() int {
	return l.list.Len()
}

func (l *lru[K, V]) Cap() int {
	return l.size
}

func (l *lru[K, V]) Purge() {
	for e := l.list.Back(); e != nil; e = e.Prev() {
		k := e.Value.(Item[K, V]).K
		l.items.Del(k)
	}
	l.list.Init()
}

func (l *lru[K, V]) Keys() []K {
	keys := make([]K, 0, l.Len())
	for e := l.list.Back(); e != nil; e = e.Prev() {
		keys = append(keys, e.Value.(Item[K, V]).K)
	}

	return keys
}

func (l *lru[K, V]) Values() []V {
	values := make([]V, 0, l.Len())
	for e := l.list.Back(); e != nil; e = e.Prev() {
		v, _ := l.Peek(e.Value.(Item[K, V]).K)
		values = append(values, v)
	}

	return values
}

func (l *lru[K, V]) Contains(key K) bool {
	return l.items.Exists(key)
}

func (l *lru[K, V]) Add(key K, val V) (evicted bool) {
	element, exist := l.items.Get(key)
	if exist {
		l.list.MoveToFront(element.V)
		element.V.Value = Item[K, V]{K: key, V: val}

		return evicted
	}
	el := l.list.PushFront(Item[K, V]{K: key, V: val})
	l.items.Put(key, &Item[K, *list.Element]{
		K: key, V: el,
	})

	if l.list.Len() > l.size {
		_, evicted = l.removeOldest()
	}

	return evicted
}

func (l *lru[K, V]) GetOldest() (key K, val V, ok bool) {
	el := l.list.Back()
	if el == nil {
		return key, val, ok
	}
	l.list.MoveToFront(el)

	item := el.Value.(Item[K, V])

	return item.K, item.V, true
}

func (l *lru[K, V]) RemoveOldest() (K, bool) {
	return l.removeOldest()
}

func (l *lru[K, V]) removeOldest() (key K, ok bool) {
	last := l.list.Back()
	if last == nil {
		return key, false
	}
	item := last.Value.(Item[K, V])

	l.del(item.K, last)

	return item.K, l.onEvict(item)
}

func (l *lru[K, V]) Remove(key K) bool {
	elem, found := l.items.Get(key)
	if !found {
		return false
	}
	l.del(key, elem.V)

	return l.onEvict(elem.V.Value.(Item[K, V]))
}

func (l *lru[K, V]) onEvict(item Item[K, V]) bool {
	if l.onEvicted != nil {
		l.onEvicted(item.K, item.V)
	}

	return true
}

func (l *lru[K, V]) del(key K, el *list.Element) {
	l.list.Remove(el)
	l.items.Del(key)
}

func (l *lru[K, V]) Resize(size int) int {
	diff := l.Len() - size
	if diff < 0 {
		diff = 0
	}
	for i := 0; i < diff; i++ {
		l.removeOldest()
	}
	l.size = size

	return diff
}

func (l *lru[K, V]) Peek(key K) (val V, ok bool) {
	v, ok := l.items.Get(key)
	if !ok {
		return val, ok
	}

	return v.V.Value.(Item[K, V]).V, ok
}

func (l *lru[K, V]) Get(key K) (v V, exist bool) {
	element, exist := l.items.Get(key)
	if exist {
		l.list.MoveToFront(element.V)

		return element.V.Value.(Item[K, V]).V, exist
	}

	return v, exist
}
