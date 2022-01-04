package mem

import (
	"micro-libs/utils/errors"
	"time"
)

var (
	ErrTableNotFound = errors.NotFound("not found mem table")
	ErrDataNotFound  = errors.NotFound("not found data")
	ErrInvalidFilter = errors.Invalid("invalid filter")
	ErrRedisError    = errors.Server("redis error")
)

const (
	DefaultPrefixData  = "CACHE_DATA"  // 数据缓存前缀
	DefaultPrefixIndex = "CACHE_INDEX" // 索引缓存
)

const (
	DefaultNilDataTTL         = 1800 * time.Second // Nil数据缓存生成周期
	DefaultStateCheckInterval = 60 * time.Second   // 内存数据状态监测
	DefaultStateExpireTime    = 3600 * time.Second // 内存数据过期时间
	DefaultIndexExpireTime    = 3600 * time.Second // 外键缓存过期时间
)

type Option func(o *Options)

type Options struct {
	PrefixData         string        // 数据缓存前缀
	PrefixIndex        string        // 索引缓存前缀
	NilDataTTL         time.Duration // Nil数据缓存生成周期
	StateCheckInterval time.Duration // 内存数据状态监测
	StateExpireTime    time.Duration // 内存数据过期时间
	IndexExpireTime    time.Duration // 外键缓存过期时间
}

func newOptions(opts ...Option) *Options {
	options := &Options{
		PrefixData:         DefaultPrefixData,
		PrefixIndex:        DefaultPrefixIndex,
		NilDataTTL:         DefaultNilDataTTL,
		StateCheckInterval: DefaultStateCheckInterval,
		StateExpireTime:    DefaultStateExpireTime,
		IndexExpireTime:    DefaultIndexExpireTime,
	}

	for _, o := range opts {
		o(options)
	}

	return options
}

func WithPrefixData(prefix string) Option {
	return func(o *Options) {
		o.PrefixData = prefix
	}
}

func WithPrefixIndex(prefix string) Option {
	return func(o *Options) {
		o.PrefixIndex = prefix
	}
}

func WithNilDataTTL(ttl time.Duration) Option {
	return func(o *Options) {
		o.NilDataTTL = ttl
	}
}

func WithStateCheckInterval(s time.Duration) Option {
	return func(o *Options) {
		o.StateCheckInterval = s
	}
}

func WithStateExpireTime(s time.Duration) Option {
	return func(o *Options) {
		o.StateExpireTime = s
	}
}

func WithIndexExpireTime(s time.Duration) Option {
	return func(o *Options) {
		o.IndexExpireTime = s
	}
}
