package ttlcache

/*
 *
 * Created by 0x5010 on 2018/01/03.
 * pf
 * https://github.com/0x5010/pf
 *
 * Copyright 2018 0x5010.
 * Licensed under the MIT license.
 *
 */

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

// CoverageCache 缓存
type CoverageCache struct {
	mutex sync.RWMutex
	items map[string]*item
}

// Set 缓存数据
func (cache *CoverageCache) Set(key string, data []byte, ttl time.Duration) {
	cache.mutex.Lock()
	defer cache.mutex.Unlock()
	item := &item{
		data: data,
		ttl:  ttl,
	}
	item.touch()
	cache.items[key] = item
}

// Get 获取缓存
func (cache *CoverageCache) Get(key string) (data []byte, found bool) {
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

// Count 获取缓存个数
func (cache *CoverageCache) Count() int {
	cache.mutex.RLock()
	defer cache.mutex.RUnlock()
	count := len(cache.items)
	return count
}

func (cache *CoverageCache) cleanup() {
	cache.mutex.Lock()
	for key, item := range cache.items {
		if item.expired() {
			delete(cache.items, key)
		}
	}
	cache.mutex.Unlock()
}

func (cache *CoverageCache) clean(duration time.Duration) {
	cleanTicker := time.NewTicker(duration)
	defer cleanTicker.Stop()

	for {
		select {
		case <-cleanTicker.C:
			cache.cleanup()
		}
	}
}

// New 新建缓存
func New(duration time.Duration) *CoverageCache {
	cache = &CoverageCache{
		items: map[string]*item{},
	}
	go cache.clean(duration)
	return cache
}

// Cache 默认全局缓存
var cache *CoverageCache

// Get 从默认缓存中读取
func Get(key string) ([]byte, bool) {
	return cache.Get(key)
}

// Set 设置缓存到默认缓存中
func Set(key string, value []byte, ttl time.Duration) {
	cache.Set(key, value, ttl)
}

// Count 当前缓存个数
func Count() int {
	return cache.Count()
}

func init() {
	cache = New(time.Duration(20 * time.Minute))
}
