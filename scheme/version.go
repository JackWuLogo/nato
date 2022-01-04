package scheme

import (
	"context"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	mgo "micro-libs/store/mongo"
	rds "micro-libs/store/redis"
	"micro-libs/utils/dtype"
	"micro-libs/utils/log"
	"micro-libs/utils/tool"
)

const (
	DefaultSchemeName        = "game_scheme"
	DefaultSchemeVersionName = "scheme_version"
)

type Version struct {
	dbName  string     // 数据库名称
	verName string     // 版本库名称
	mongo   *mgo.Store // 数据库
	redis   *rds.Store // 缓存
}

func (v *Version) Init(opts ...Option) {
	for _, opt := range opts {
		opt(v)
	}
}

// DbName 数据库名称
func (v *Version) DbName() string {
	return v.dbName
}

// VerName 版本表名称
func (v *Version) VerName() string {
	return v.verName
}

// GetCacheName 缓存名称
func (v *Version) GetCacheName() string {
	return rds.GetCacheName("SCHEME", v.dbName, v.verName)
}

// Check 版本库检查
func (v *Version) Check() error {
	cols, err := v.mongo.ListCollectionNames(v.dbName)
	if err != nil {
		return err
	}

	if !tool.InStrSlice(v.verName, cols) {
		// 创建索引
		if _, err := v.mongo.C(v.verName, v.dbName).Indexes().CreateMany(context.TODO(), VersionModelIndex); err != nil {
			log.Error("create scheme index failure, error: %s", err)
			return err
		}
	}

	log.Info("[%s] check collections success ...", v.dbName)

	return nil
}

// SetCache 设置版本缓存
func (v *Version) SetCache(ctx context.Context) (map[string]int64, error) {
	cur, err := v.mongo.C(v.verName, v.dbName).Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"key": 1, "version": 1}))
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer cur.Close(ctx)

	var rows []*VersionModel
	if err := cur.All(ctx, &rows); err != nil {
		return nil, err
	}

	var version = make(map[string]int64, len(rows))
	var values []interface{}
	for _, row := range rows {
		version[row.Key] = row.Version
		values = append(values, row.Key, row.Version)
	}

	if err := v.redis.Client().Del(ctx, v.GetCacheName()).Err(); err != nil {
		return nil, err
	}
	if err := v.redis.Client().HSet(ctx, v.GetCacheName(), values).Err(); err != nil {
		return nil, err
	}

	return version, nil
}

// GetCache 获取版本缓存
func (v *Version) GetCache(ctx context.Context) (map[string]int64, error) {
	res, err := v.redis.Client().HGetAll(ctx, v.GetCacheName()).Result()
	if err != nil {
		if err != redis.Nil {
			return nil, err
		}
	}

	if len(res) == 0 {
		return v.SetCache(ctx)
	}

	var version = make(map[string]int64, len(res))
	for k, v := range res {
		version[k] = dtype.ParseInt64(v)
	}

	return version, nil
}

// DelCache 删除指定缓存
func (v *Version) DelCache(ctx context.Context, keys ...string) error {
	if err := v.redis.Client().HDel(ctx, v.GetCacheName(), keys...).Err(); err != nil && err != redis.Nil {
		return err
	}
	return nil
}

func newVersion() *Version {
	return &Version{
		dbName:  DefaultSchemeName,
		verName: DefaultSchemeVersionName,
	}
}

type Option func(s *Version)

func WithDbName(name string) Option {
	return func(s *Version) {
		s.dbName = name
	}
}

func WithVerName(name string) Option {
	return func(s *Version) {
		s.verName = name
	}
}

func WithMongoDB(client *mgo.Store) Option {
	return func(s *Version) {
		s.mongo = client
	}
}

func WithRedis(client *rds.Store) Option {
	return func(s *Version) {
		s.redis = client
	}
}
