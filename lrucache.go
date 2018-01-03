package lrucache

import (
	"sync"
	"time"
)

type item struct {
	sync.RWMutex
	data    []byte
	ttl     time.Duration
	expires *time.Time
}

func (i *item) touch() {
	i.Lock()
	defer i.Unlock()
	expiration := time.Now().Add(i.ttl)
	i.expires = &expiration
}

func (i *item) expired() bool {
	var value bool
	i.RLock()
	defer i.RUnlock()
	if i.expires == nil {
		value = true
	} else {
		value = i.expires.Before(time.Now())
	}
	return value
}

type coverageCache struct {
	mutex sync.RWMutex
	items map[string]*item
}

func (cache *coverageCache) set(key string, data []byte, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	item := &item{
		data: data,
		ttl:  ttl,
	}
	item.touch()
	cache.items[key] = item
}

func (cache *coverageCache) get(key string) (data []byte, found bool) {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	item, exists := cache.items[key]
	if !exists || item.expired() {
		data = nil
		found = false
	} else {
		data = item.data
		found = true
	}
	return
}

func (cache *coverageCache) count() int {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	count := len(cache.items)
	return count
}

func (cache *coverageCache) cleanup() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

func (cache *coverageCache) clean(duration time.Duration) {
	cleanTicker := time.NewTicker(duration)
	defer cleanTicker.Stop()

	for {
		select {
		case <-cleanTicker.C:
			cache.cleanup()
		}
	}
}

// Cache 默认全局缓存
var cache *coverageCache

// Get 从默认缓存中读取
func Get(key string) ([]byte, bool) {
	return cache.get(key)
}

// Set 设置缓存到默认缓存中
func Set(key string, value []byte, ttl time.Duration) {
	cache.set(key, value, ttl)
}

func init() {
	cache = &coverageCache{
		items: map[string]*item{},
	}
	go cache.clean(time.Duration(20 * time.Minute))
}
