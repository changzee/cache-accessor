// cache 是一个针对web应用设计的缓存系统
package cache

import (
	"errors"
	"time"
)

// ErrNotExist 缓存项不存在
var ErrNotExist = errors.New("item is not exist")

// CacheInterface 为缓存接口协议
type CacheInterface interface {
	// Get 获取key值对应的缓存, 不存在则返回nil
	Get(key string, value interface{}) error

	// Set 设置key值对应的缓存
	Set(key string, value interface{}, ttl ...time.Duration) error

	// Delete 删除key值对应的缓存
	Delete(key string) error
}
