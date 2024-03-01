package lru

type ILru[K comparable, V any] interface {
	Len() int
	Cap() int
	Purge()
	Keys() []K
	Values() []V
	Contains(key K) bool
	Add(key K, val V) (evicted bool)
	GetOldest() (key K, val V, ok bool)
	RemoveOldest() (K, bool)
	Remove(key K) bool
	Resize(size int) int
	Peek(key K) (V, bool)
	Get(key K) (V, bool)
}
