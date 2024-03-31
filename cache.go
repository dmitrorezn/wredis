package wredis

import (
	"time"

	"github.com/dmitrorezn/golang-lru/v2/expirable"
	"github.com/vmihailenco/go-tinylfu"
)

type Cache interface {
	Get(key string) (any, bool)
	Set(key string, v any, ttl ...time.Duration) bool
	Del(key string)
}

type LRUCache struct {
	cache *expirable.LRU[string, any]
}

func NewLRUCache(size int, onEvict callback[string, any], ttl time.Duration) *LRUCache {
	return &LRUCache{
		cache: expirable.NewLRU[string, any](size, onEvict.Call, ttl),
	}
}

func (lfuc LRUCache) Get(key string) (any, bool) {
	return lfuc.cache.Get(key)
}

func (lfuc LRUCache) Set(key string, v any, ttl ...time.Duration) bool {
	return lfuc.cache.Set(key, v, ttl...)
}

func (lfuc LRUCache) Del(key string) {
	lfuc.cache.Remove(key)
}

var _ Cache = LRUCache{}

type LFUCache struct {
	onEvict callback[string, any]
	ttl     time.Duration
	cache   *tinylfu.SyncT
}

func NewLFUCache(size int, samples int, onEvict callback[string, any], ttl time.Duration) *LFUCache {
	return &LFUCache{
		onEvict: onEvict,
		ttl:     ttl,
		cache:   tinylfu.NewSync(size, samples),
	}
}

func (lfuc LFUCache) Get(key string) (any, bool) {
	return lfuc.cache.Get(key)
}

func (lfuc LFUCache) Set(key string, v any, ttl ...time.Duration) bool {
	lfuc.cache.Set(&tinylfu.Item{
		Key:   key,
		Value: v,
		OnEvict: func() {
			lfuc.onEvict.Call(key, v)
		},
		ExpireAt: time.Now().Add(append(ttl, lfuc.ttl)[0]),
	})

	return false
}

func (lfuc LFUCache) Del(key string) {
	lfuc.cache.Del(key)
}

var _ Cache = LFUCache{}
