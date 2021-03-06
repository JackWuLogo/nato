package redis

import (
	"context"
	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v8"
	"micro-libs/utils/errors"
	"micro-libs/utils/log"
	"strings"
)

var (
	DefaultDb = 0
)

// Redis 存储
type Store struct {
	opts   *Options
	client *redis.Client
	locker *redislock.Client
}

func (s *Store) Init(opts ...Option) {
	for _, o := range opts {
		o(s.opts)
	}
}

func (s *Store) Opts() *Options {
	return s.opts
}

func (s *Store) Locker() *redislock.Client {
	return s.locker
}

func (s *Store) Client() *redis.Client {
	return s.client
}

func (s *Store) Connect() error {
	if s.client != nil {
		return nil
	}

	if s.opts.RawUrl == "" {
		s.opts.RawUrl = "redis://127.0.0.1:6379"
	} else if !strings.HasPrefix(s.opts.RawUrl, "redis://") {
		s.opts.RawUrl = "redis://" + s.opts.RawUrl
	}

	opts, err := redis.ParseURL(s.opts.RawUrl)
	if err != nil {
		return errors.Wrap(err, "redis url invalid")
	}

	// 设置默认参数
	opts.Username = ""
	opts.DB = s.opts.Db
	opts.PoolSize = s.opts.PoolSize
	opts.MinIdleConns = s.opts.MinIdleConns
	opts.MaxConnAge = s.opts.MaxConnAge
	opts.MaxRetries = s.opts.MaxRetries
	opts.MinRetryBackoff = s.opts.MinRetryBackoff
	opts.MaxRetryBackoff = s.opts.MaxRetryBackoff
	opts.ReadTimeout = s.opts.ReadTimeout
	opts.WriteTimeout = s.opts.WriteTimeout
	opts.PoolTimeout = s.opts.PoolTimeout
	opts.IdleTimeout = s.opts.IdleTimeout

	s.client = redis.NewClient(opts)
	if err := s.client.Ping(context.Background()).Err(); err != nil {
		return errors.Wrap(err, "redis connect error")
	}

	// 启用分布式锁
	s.locker = redislock.New(s.client)

	log.Debug("Store [redis] Connect to %s, DB: %d", opts.Addr, opts.DB)

	return nil
}

func (s *Store) Disconnect() error {
	if s.client != nil {
		if err := s.client.Close(); err != nil {
			return err
		}
		s.client = nil
	}
	return nil
}

// Do 执行命令
func (s *Store) Do(args ...interface{}) *redis.Cmd {
	ctx := context.Background()
	cmd := redis.NewCmd(ctx, args...)
	s.client.Process(ctx, cmd)
	return cmd
}

// HSetStruct 设置结构体
func (s *Store) HSetStruct(key string, result interface{}) error {
	ctx := context.Background()
	cmd := redis.NewStatusCmd(ctx, Args{CmdHSet}.Add(key).AddFlat(result)...)
	s.client.Process(ctx, cmd)
	_, err := cmd.Result()
	return err
}

func NewStore(opts ...Option) *Store {
	rs := &Store{
		opts: newOptions(opts...),
	}
	return rs
}
