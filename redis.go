package cache

import (
	"encoding/json"
	"github.com/go-redis/redis"
	jsoniter "github.com/json-iterator/go"
	"time"
)

// RedisCache 实现了CacheInterface接口，用于kv缓存, 序列化方式暂时只支持json，后期可优化
type RedisCache struct {
	client            *redis.Client
	defaultExpiration time.Duration
}

const (
	defaultMaxIdle        = 5
	defaultMaxActive      = 0
	defaultTimeoutIdle    = 240
	defaultTimeoutConnect = 10000
	defaultTimeoutRead    = 5000
	defaultTimeoutWrite   = 5000
	defaultHost           = "localhost:6379"
	defaultProtocol       = "tcp"
)

// RedisOpts redis连接选项
type RedisOpts struct {
	MaxIdle           int
	MaxActive         int
	Protocol          string
	Host              string
	Password          string
	TimeoutConnect    int
	TimeoutRead       int
	TimeoutWrite      int
	TimeoutIdle       int
	DefaultExpiration time.Duration
}

// padDefaults 填充默认选项
func (r RedisOpts) padDefaults() RedisOpts {
	if r.MaxIdle == 0 {
		r.MaxIdle = defaultMaxIdle
	}
	if r.MaxActive == 0 {
		r.MaxActive = defaultMaxActive
	}
	if r.TimeoutConnect == 0 {
		r.TimeoutConnect = defaultTimeoutConnect
	}
	if r.TimeoutIdle == 0 {
		r.TimeoutIdle = defaultTimeoutIdle
	}
	if r.TimeoutRead == 0 {
		r.TimeoutRead = defaultTimeoutRead
	}
	if r.TimeoutWrite == 0 {
		r.TimeoutWrite = defaultTimeoutWrite
	}
	if r.Host == "" {
		r.Host = defaultHost
	}
	if r.Protocol == "" {
		r.Protocol = defaultProtocol
	}
	return r
}

// Set 设置redis缓存
func (rc *RedisCache) Set(key string, value interface{}, ttl ...time.Duration) error {
	if str, err := jsoniter.MarshalToString(value); err != nil {
		return err
	} else {
		expire := rc.defaultExpiration
		if len(ttl) > 0 {
			expire = ttl[0]
		}
		return rc.client.Set(key, str, expire).Err()
	}
}

// Delete 删除缓存
func (rc *RedisCache) Delete(key string) error {
	return rc.client.Del(key).Err()
}

// Get 获取缓存数据
func (rc *RedisCache) Get(key string, value interface{}) error {
	b, err := rc.client.Get(key).Bytes()
	if err == redis.Nil {
		return ErrNotExist
	}
	if err != nil {
		return err
	}

	return json.Unmarshal(b, value)
}

// NewRedisCache 新建redis缓存
func NewRedisCache(opts RedisOpts) *RedisCache {
	opts = opts.padDefaults()
	toc := time.Millisecond * time.Duration(opts.TimeoutConnect)
	tor := time.Millisecond * time.Duration(opts.TimeoutRead)
	tow := time.Millisecond * time.Duration(opts.TimeoutWrite)
	toi := time.Duration(opts.TimeoutIdle) * time.Second
	opt := &redis.Options{
		Addr:               opts.Host,
		DB:                 0,
		DialTimeout:        toc,
		ReadTimeout:        tor,
		WriteTimeout:       tow,
		PoolSize:           opts.MaxActive,
		PoolTimeout:        30 * time.Second,
		IdleTimeout:        toi,
		Password:           opts.Password,
		IdleCheckFrequency: 500 * time.Millisecond,
	}

	c := redis.NewClient(opt)
	return &RedisCache{client: c, defaultExpiration: opts.DefaultExpiration}
}
