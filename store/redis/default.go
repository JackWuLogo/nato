package redis

import (
	"bytes"
	"context"
	"fmt"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"micro-libs/app"
	"micro-libs/utils/dtype"
	"strconv"
	"sync"
	"time"
)

var (
	single *Store
	once   sync.Once
)

func S() *Store {
	once.Do(func() {
		single = NewStore()
	})
	return single
}

func Connect() error {
	return S().Connect()
}

func Disconnect() error {
	return S().Disconnect()
}

// Client 获取客户端
func Client() *redis.Client {
	return S().client
}

func Locker() *redislock.Client {
	return S().Locker()
}

func Do(args ...interface{}) *redis.Cmd {
	return S().Do(args...)
}

// 分布式锁 (默认生存周期, 默认重试时间)
func Lock(ctx context.Context, key string, index interface{}) (*redislock.Lock, error) {
	return LockBackoff(ctx, key, index, 3*time.Second, 100*time.Millisecond)
}

// 分布式锁 (指定生存周期)
func LockTTL(ctx context.Context, key string, index interface{}, ttl time.Duration) (*redislock.Lock, error) {
	return LockBackoff(ctx, key, index, ttl, 100*time.Millisecond)
}

// 分布式锁 (指定生存周期, 设置重试时间)
func LockBackoff(ctx context.Context, key string, index interface{}, ttl time.Duration, backoff time.Duration) (*redislock.Lock, error) {
	opt := &redislock.Options{}
	if backoff > 0 {
		opt.RetryStrategy = redislock.LinearBackoff(backoff)
	}
	keyName := fmt.Sprintf("LOCK:%s:%s", key, dtype.ParseStr(index))
	return Locker().Obtain(ctx, keyName, ttl, opt)
}

// GetCacheName 生成缓存名称
func GetCacheName(prefix string, tags ...interface{}) string {
	var buf = new(bytes.Buffer)
	buf.WriteString(strconv.FormatInt(app.Opts.ClusterId, 10))
	buf.WriteString(":")
	buf.WriteString(prefix)
	for _, tag := range tags {
		buf.WriteString(":")
		switch tag.(type) {
		case string:
			buf.WriteString(tag.(string))
		case int:
			buf.WriteString(strconv.FormatInt(int64(tag.(int)), 10))
		case int32:
			buf.WriteString(strconv.FormatInt(int64(tag.(int32)), 10))
		case int64:
			buf.WriteString(strconv.FormatInt(tag.(int64), 10))
		case uint:
			buf.WriteString(strconv.FormatUint(uint64(tag.(uint)), 10))
		case uint32:
			buf.WriteString(strconv.FormatUint(uint64(tag.(uint32)), 10))
		case uint64:
			buf.WriteString(strconv.FormatUint(tag.(uint64), 10))
		default:
			buf.WriteString(fmt.Sprint(tag))
		}
	}
	return buf.String()
}
