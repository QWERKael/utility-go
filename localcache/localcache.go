package localcache

import (
	"sync"
	"time"
)

type CachedValue[V any] struct {
	Expired     time.Time
	WithExpired bool
	Value       V
}

func NewCachedValue[V any](value V, expiration time.Duration) CachedValue[V] {
	if expiration > 0 {
		return CachedValue[V]{Expired: time.Now().Add(expiration), WithExpired: true, Value: value}
	} else {
		return CachedValue[V]{Expired: time.Time{}, WithExpired: false, Value: value}
	}
}

func NewCachedValueWithExpiredTime[V any](value V, expiredTime time.Time) CachedValue[V] {
	return CachedValue[V]{Expired: expiredTime, WithExpired: true, Value: value}
}
func NewCachedValueWithoutExpiredTime[V any](value V) CachedValue[V] {
	return CachedValue[V]{Expired: time.Time{}, WithExpired: false, Value: value}
}

func (c *CachedValue[V]) IsExpired() bool {
	if c.WithExpired && c.Expired.Before(time.Now()) {
		return true
	}
	return false
}

func (c *CachedValue[V]) Get() (value V) {
	if c.IsExpired() {
		return value
	}
	return c.Value
}

// Expire 强制 CachedValue 过期
func (c *CachedValue[V]) Expire() {
	c.WithExpired = true
	c.Expired = time.Unix(0, 0)
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
	if v.IsExpired() {
		delete(c.cache, key)
		return value, false
	}
	return v.Value, true
}

// Set 在缓存中设置一个值
func (c *LocalCache[K, V]) Set(key K, value V, expiration time.Duration) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = NewCachedValue[V](value, expiration)
	return
}

// SetWithExpiredTime 在缓存中设置一个值，过期时间为指定的过期时间点
func (c *LocalCache[K, V]) SetWithExpiredTime(key K, value V, expiredTime time.Time) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = NewCachedValueWithExpiredTime[V](value, expiredTime)
	return
}

// SetWithoutExpiredTime 在缓存中设置一个值，不设置过期时间
func (c *LocalCache[K, V]) SetWithoutExpiredTime(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = NewCachedValueWithoutExpiredTime[V](value)
	return
}

// Delete 从缓存中删除一个值
func (c *LocalCache[K, V]) Delete(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

// Keys 列出缓存中素有的key
func (c *LocalCache[K, V]) Keys() []K {
	c.mu.RLock()
	defer c.mu.RUnlock()
	keys := make([]K, 0)
	for key, value := range c.cache {
		if value.IsExpired() {
			// 如果有过期时间，并且已过期，则删除key
			delete(c.cache, key)
		} else {
			// 否则，记录key
			keys = append(keys, key)
		}
	}
	return keys
}

// Values 列出缓存中素有的value
func (c *LocalCache[K, V]) Values() []V {
	c.mu.RLock()
	defer c.mu.RUnlock()
	values := make([]V, 0)
	for key, cachedValue := range c.cache {
		if cachedValue.IsExpired() {
			// 如果有过期时间，并且已过期，则删除key
			delete(c.cache, key)
		} else {
			// 否则，记录key
			values = append(values, cachedValue.Value)
		}
	}
	return values
}
