package localcache

import (
	"sync"
)

// LocalCache 线程安全的本地缓存
type LocalCache[K comparable, V any] struct {
	mu    sync.RWMutex
	cache map[K]V
}

// NewLocalCache 创建一个新的本地缓存
func NewLocalCache[K comparable, V any]() *LocalCache[K, V] {
	return &LocalCache[K, V]{
		cache: make(map[K]V),
	}
}

// Get 从缓存中获取一个值
func (c *LocalCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	value, ok = c.cache[key]
	return
}

// Set 在缓存中设置一个值
func (c *LocalCache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = value
}

// Delete 从缓存中删除一个值
func (c *LocalCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}
