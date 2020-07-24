package cache

import (
	"errors"
	"github.com/patrickmn/go-cache"
	"reflect"
	"time"
)

// MemoryCache 是一个基于内存的缓存，一般来说缓存的数据时效要短一些
type MemoryCache struct {
	expire time.Duration
	cache  *cache.Cache
}

// NewMemoryCache 用于新建一个内存缓存
func NewMemoryCache(defaultExpiration, cleanupInterval time.Duration) *MemoryCache {
	return &MemoryCache{
		expire: defaultExpiration,
		cache:  cache.New(defaultExpiration, cleanupInterval),
	}
}

// Get 用于获取缓存数据
func (mc *MemoryCache) Get(key string, value interface{}) error {
	cacheValue, ok := mc.cache.Get(key)
	if !ok {
		return ErrNotExist
	}

	if reflect.TypeOf(value).Elem().Kind() == reflect.TypeOf(cacheValue).Kind() {
		reflect.ValueOf(value).Elem().Set(reflect.ValueOf(cacheValue))
		return nil
	} else {
		return errors.New("error type assigned")
	}
}

// Set 设置key值对应的缓存
func (mc *MemoryCache) Set(key string, value interface{}, ttl ...time.Duration) error {
	expire := mc.expire
	if len(ttl) > 0 {
		expire = ttl[0]
	}

	mc.cache.Set(key, value, expire)
	return nil
}

// Delete 删除key值对应的缓存
func (mc *MemoryCache) Delete(key string) error {
	mc.cache.Delete(key)
	return nil
}
