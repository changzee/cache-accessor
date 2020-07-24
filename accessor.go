package cache

import (
	"golang.org/x/sync/singleflight"
	"reflect"
)

// CacheAccessor 用于封装缓存访问逻辑
type CacheAccessor struct {
	L1Cache CacheInterface // 一级缓存
	L2Cache CacheInterface // 二级缓存

	group *singleflight.Group // 重复抑制，防止缓存击穿
}

// Get 获取key值对应的缓存数据，顺序由L1 -》 L2
func (accessor *CacheAccessor) Get(key string, value interface{}) error {
	// 如果有配置一级缓存先从一级缓存取数据
	if accessor.L1Cache != nil {
		err := accessor.L1Cache.Get(key, value)
		if err != ErrNotExist { // err不为ErrNotExist说明err为nil或发生其他错误，直接返回
			return err
		}
	}

	// 如果有配置二级缓存则从二级缓存取数据
	if accessor.L2Cache != nil {
		err := accessor.L2Cache.Get(key, value)
		if err != ErrNotExist { // err不为ErrNotExist说明err为nil或发生其他错误，直接返回
			return err
		}
	}
	return ErrNotExist
}

// Set 设置缓存数据, 顺序由L2-》L1，中断则向上抛出错误，防止数据跳动
func (accessor *CacheAccessor) Set(key string, value interface{}) error {
	if accessor.L2Cache != nil { // 写入二级缓存层
		if err := accessor.L2Cache.Set(key, value); err != nil {
			return err
		}
	}
	if accessor.L1Cache != nil { // 写入一级缓存层
		if err := accessor.L1Cache.Set(key, value); err != nil {
			return err
		}
	}
	return nil
}

// Delete 删除缓存数据
func (accessor *CacheAccessor) Delete(key string) error {
	if accessor.L2Cache != nil { // 删除二级缓存数据
		if err := accessor.L2Cache.Delete(key); err != nil {
			return err
		}
	}
	if accessor.L1Cache != nil { // 删除一级缓存数据
		if err := accessor.L1Cache.Delete(key); err != nil {
			return err
		}
	}
	return nil
}

// LazyGet 用于惰性加载数据
func (accessor *CacheAccessor) LazyGet(key string, res interface{}, slower func() (interface{}, error)) error {
	// 先从缓存中拿数据
	if err := accessor.Get(key, res); err != ErrNotExist {
		return err
	}

	// 重复抑制
	group := accessor.group
	slowerRet, err, _ := group.Do(key, slower)
	if err == nil {
		reflect.ValueOf(res).Elem().Set(reflect.ValueOf(slowerRet))
		err = accessor.Set(key, slowerRet)
		if err != nil {
			return err
		}
	}
	return err
}

// NewCacheAccessor 新建一个缓存访问器
func NewCacheAccessor(levelOneCache CacheInterface, levelTwoCache CacheInterface) *CacheAccessor {
	return &CacheAccessor{
		L1Cache: levelOneCache,
		L2Cache: levelTwoCache,
		group:   &singleflight.Group{},
	}
}
