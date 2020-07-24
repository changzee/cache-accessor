![Golang](https://golang.org/lib/godoc/images/footer-gopher.jpg)

# cache-accessor

`cache-accessor` 为golang下多级缓存访问器实现

喜欢的话请给个小星星吧~

## Keywords
1. 多级缓存
2. 基于`singleflight`防止缓存击穿
3. Redis/Memory

## 版本要求

* Go1.13.4

## 安装
    go mod github.com/changzee/cache-accessor
    go mod vendor
 
## 使用
    accessor := NewCacheAccessor(
        NewMemoryCache(5*time.Second, 180),
    )
    
    var ret string
    accessor.LazyGet("cache_key", &ret, function() (interface{}, error) {
        return "cache_value", nil
    })

