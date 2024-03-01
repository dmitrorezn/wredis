package lru

import (
	"container/list"
	"fmt"
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

type hmap[K comparable, V any] map[K]Item[K, V]

func (m hmap[K, V]) Put(key K, val Item[K, V]) {
	m[key] = val
}
func (m hmap[K, V]) Del(key K) {
	delete(m, key)
}
func (m hmap[K, V]) Get(key K) (Item[K, V], bool) {
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

	values    hmap[K, V]
	elements  hmap[K, *list.Element]
	onEvicted func(K, V)

	list *list.List
}

func newLRU[K comparable, V any](size int, onEvicted func(K, V)) *lru[K, V] {
	return &lru[K, V]{
		size:      size,
		values:    newHmap[K, V](size),
		elements:  newHmap[K, *list.Element](size),
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
		k := e.Value.(K)
		l.elements.Del(k)
		l.values.Del(k)
	}
	l.list.Init()
}

func (l *lru[K, V]) Keys() []K {
	keys := make([]K, 0, l.Len())
	for e := l.list.Back(); e != nil; e = e.Prev() {
		keys = append(keys, e.Value.(K))
	}

	return keys
}

func (l *lru[K, V]) Values() []V {
	values := make([]V, 0, l.Len())
	for e := l.list.Back(); e != nil; e = e.Prev() {
		v, _ := l.Peek(e.Value.(K))
		values = append(values, v)
	}

	return values
}

func (l *lru[K, V]) Contains(key K) bool {
	return l.values.Exists(key)
}

func (l *lru[K, V]) Add(key K, val V) (evicted bool) {
	element, exist := l.elements.Get(key)
	if exist {
		l.list.MoveToFront(element.V)
	} else {
		el := l.list.PushFront(key)
		l.elements.Put(key, Item[K, *list.Element]{
			K: key, V: el,
		})
	}
	l.values.Put(key, Item[K, V]{
		K: key, V: val,
	})
	if l.list.Len() > l.size {
		l.removeOldest()
		evicted = true
	}

	return evicted
}

func (l *lru[K, V]) GetOldest() (key K, val V, ok bool) {
	el := l.list.Back()
	if el == nil {
		return
	}
	l.list.MoveToFront(el)

	key = el.Value.(K)
	val, ok = l.Peek(key)

	return
}

func (l *lru[K, V]) RemoveOldest() (K, bool) {
	return l.removeOldest()
}

func (l *lru[K, V]) removeOldest() (key K, ok bool) {
	last := l.list.Back()
	if last == nil {
		return key, false
	}
	key = last.Value.(K)
	fmt.Println("removeOldest ", key)

	defer l.del(key, last)

	return key, l.onEvict(key)
}

func (l *lru[K, V]) Remove(key K) bool {
	fmt.Println("l", l.elements)
	elem, found := l.elements.Get(key)
	if !found {
		return false
	}
	defer l.del(key, elem.V)

	return l.onEvict(key)
}

func (l *lru[K, V]) onEvict(key K) bool {
	val, found := l.values.Get(key)
	if !found {
		return false
	}

	if l.onEvicted != nil {
		l.onEvicted(key, val.V)
	}

	return true
}

func (l *lru[K, V]) del(key K, el *list.Element) {
	l.list.Remove(el)
	l.elements.Del(key)
	l.values.Del(key)
}

func (l *lru[K, V]) Resize(size int) int {
	diff := l.Len() - size
	if diff < 0 {
		return 0
	}
	for i := 0; i < diff; i++ {
		l.removeOldest()
	}
	l.size = size

	return diff
}

func (l *lru[K, V]) Peek(key K) (V, bool) {
	v, ok := l.values.Get(key)

	return v.V, ok
}

func (l *lru[K, V]) Get(key K) (V, bool) {
	element, exist := l.elements.Get(key)
	if exist {
		l.list.MoveToFront(element.V)
	}
	item, found := l.values.Get(key)

	return item.V, found
}
