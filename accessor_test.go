package cache

import (
	"fmt"
	"github.com/alicebob/miniredis"
	"github.com/stretchr/testify/assert"

	"os"
	"testing"
	"time"
)

func fatalTestError(fmtStr string, args ...interface{}) {
	_, _ = fmt.Fprintf(os.Stderr, fmtStr, args...)
	os.Exit(1)
}

func TestCacheAccessor_LazyGet(t *testing.T) {
	// mock redis
	miniRedisClient, err := miniredis.Run()
	if err != nil {
		fatalTestError("Error creating mini redis: %v\n", err)
	}

	accessor := NewCacheAccessor(
		NewMemoryCache(5*time.Second, 180),
		NewRedisCache(RedisOpts{
			Host: miniRedisClient.Addr(),
		}),
	)

	var ret string
	err = accessor.LazyGet("test", &ret, func() (interface{}, error) {
		return "test", nil
	})
	assert.NoError(t, err)
	assert.Equal(t, ret, "test")

	err = accessor.Get("test", &ret)
	assert.NoError(t, err)
	assert.Equal(t, ret, "test")

	type Foo struct{ A string }

	A := Foo{A: "asdasd"}
	err = accessor.LazyGet("A", &A, func() (interface{}, error) {
		return struct{ A string }{A: "adasda"}, nil
	})
	assert.NoError(t, err)

	var B *Foo
	err = accessor.LazyGet("B", &B, func() (interface{}, error) {
		return &Foo{A: "asdasd"}, nil
	})
	assert.NoError(t, err)

	var slice []int
	err = accessor.LazyGet("slice", &slice, func() (interface{}, error) {
		return []int{1, 2, 3}, nil
	})
	assert.NoError(t, err)
}
