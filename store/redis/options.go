package redis

import "time"

const (
	DefaultMaxRetries      = 3
	DefaultMaxConnAge      = 5 * time.Minute
	DefaultReadTimeout     = 5 * time.Second
	DefaultWriteTimeout    = 5 * time.Second
	DefaultPoolTimeout     = 5 * time.Second
	DefaultIdleTimeout     = 5 * time.Minute        // 闲置超时，默认5分钟，-1表示取消闲置超时检查
	DefaultMinRetryBackoff = 8 * time.Millisecond   // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
	DefaultMaxRetryBackoff = 512 * time.Millisecond // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
)

type Option func(o *Options)

type Options struct {
	RawUrl          string
	Db              int           // 数据库
	PoolSize        int           // 最大连接数
	MinIdleConns    int           // 在启动阶段创建指定数量的Idle连接，并长期维持idle状态的连接数不少于指定数量
	MaxRetries      int           // 最大重试次数
	MaxConnAge      time.Duration // 最大存活时间
	ReadTimeout     time.Duration // 读取超时
	WriteTimeout    time.Duration // 写入超时
	PoolTimeout     time.Duration // 等待连接超时
	IdleTimeout     time.Duration // 闲置超时，默认5分钟，-1表示取消闲置超时检查
	MinRetryBackoff time.Duration // 每次计算重试间隔时间的下限，默认8毫秒，-1表示取消间隔
	MaxRetryBackoff time.Duration // 每次计算重试间隔时间的上限，默认512毫秒，-1表示取消间隔
}

func (o *Options) Init(opts ...Option) {
	for _, opt := range opts {
		opt(o)
	}
}

func newOptions(opts ...Option) *Options {
	o := &Options{
		RawUrl:          Opts.RedisUrl,
		Db:              Opts.RedisDb,
		PoolSize:        Opts.RedisMaxPool,
		MinIdleConns:    Opts.RedisIdeConns,
		MaxRetries:      DefaultMaxRetries,
		MaxConnAge:      DefaultMaxConnAge,
		ReadTimeout:     DefaultReadTimeout,
		WriteTimeout:    DefaultWriteTimeout,
		PoolTimeout:     DefaultPoolTimeout,
		IdleTimeout:     DefaultIdleTimeout,
		MinRetryBackoff: DefaultMinRetryBackoff,
		MaxRetryBackoff: DefaultMaxRetryBackoff,
	}
	o.Init(opts...)
	return o
}

func WithUrl(url string) Option {
	return func(o *Options) {
		o.RawUrl = url
	}
}

func WithDb(db int) Option {
	return func(o *Options) {
		o.Db = db
	}
}

func WithMinIdleConns(size int) Option {
	return func(o *Options) {
		o.MinIdleConns = size
	}
}

func WithPoolSize(size int) Option {
	return func(o *Options) {
		o.PoolSize = size
	}
}

func WithMaxRetries(size int) Option {
	return func(o *Options) {
		o.MaxRetries = size
	}
}

func WithReadTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.ReadTimeout = t
	}
}

func WithWriteTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.WriteTimeout = t
	}
}

func WithIdleTimeout(t time.Duration) Option {
	return func(o *Options) {
		o.IdleTimeout = t
	}
}

func WithMinRetryBackoff(t time.Duration) Option {
	return func(o *Options) {
		o.MinRetryBackoff = t
	}
}

func WithMaxRetryBackoff(t time.Duration) Option {
	return func(o *Options) {
		o.MaxRetryBackoff = t
	}
}
