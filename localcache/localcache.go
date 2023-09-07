package localcache

import (
	"sync"
	"time"
)

type CachedValue[V any] struct {
	Expired        time.Time
	WithoutExpired bool
	Value          V
}

// LocalCache 线程安全的本地缓存
type LocalCache[K comparable, V any] struct {
	mu    sync.RWMutex
	cache map[K]CachedValue[V]
}

// NewLocalCache 创建一个新的本地缓存
func NewLocalCache[K comparable, V any]() *LocalCache[K, V] {
	return &LocalCache[K, V]{
		cache: make(map[K]CachedValue[V]),
	}
}

// Get 从缓存中获取一个值
func (c *LocalCache[K, V]) Get(key K) (value V, ok bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	v, ok := c.cache[key]
	if !ok {
		return value, false
	}
	if v.WithoutExpired {
		return v.Value, true
	} else if v.Expired.Before(time.Now()) {
		delete(c.cache, key)
		return value, false
	}
	return v.Value, true
}

// Set 在缓存中设置一个值
func (c *LocalCache[K, V]) Set(key K, value V, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if expiration > 0 {
		c.cache[key] = CachedValue[V]{Expired: time.Now().Add(expiration), WithoutExpired: false, Value: value}
	} else {
		c.cache[key] = CachedValue[V]{Expired: time.Time{}, WithoutExpired: true, Value: value}
	}
	return
}

// SetWithExpiredTime 在缓存中设置一个值，过期时间为指定的过期时间点
func (c *LocalCache[K, V]) SetWithExpiredTime(key K, value V, expiredTime time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = CachedValue[V]{Expired: expiredTime, WithoutExpired: false, Value: value}
	return
}

// Delete 从缓存中删除一个值
func (c *LocalCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}
