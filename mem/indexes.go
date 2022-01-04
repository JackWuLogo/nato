package mem

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
	"sync"
)

// 外键结构
type Index struct {
	sync.Mutex
	ma  *Indexes
	fk  string
	key string
}

// 外键管理器
type Indexes struct {
	sync.RWMutex
	table   *Table
	indexes map[string]*Index
}

func (ma *Indexes) Table() *Table {
	return ma.table
}

// 获取外键对象
func (ma *Indexes) Get(fk string) *Index {
	ma.RLock()
	if mfk, ok := ma.indexes[fk]; ok {
		ma.RUnlock()
		return mfk
	}
	ma.RUnlock()

	ma.Lock()
	mfk := &Index{
		ma:  ma,
		fk:  fk,
		key: rds.GetCacheName(ma.table.admin.opts.PrefixIndex, ma.table.Name(), fk),
	}
	ma.indexes[fk] = mfk
	ma.Unlock()

	return mfk
}

// 获取多个外键对象, map[fk]*Index
func (ma *Indexes) MultiGet(fks []string) map[string]*Index {
	var result = make(map[string]*Index, len(fks))
	for _, fk := range fks {
		result[fk] = ma.Get(fk)
	}
	return result
}

// 检查缓存是否存在
func (ma *Indexes) Exist(fk string) bool {
	res, _ := rds.Client().Exists(context.Background(), ma.Get(fk).key).Result()
	return res > 0
}

// 检查多个缓存是否存在, map[fk]bool
func (ma *Indexes) MultiExist(fks []string) map[string]bool {
	ctx := context.Background()
	var result = make(map[string]bool, len(fks))
	cmds, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, fk := range fks {
			pipe.Exists(ctx, ma.Get(fk).key)
		}
		return nil
	})
	if err != nil {
		return result
	}

	for i, cmd := range cmds {
		n, _ := cmd.(*redis.IntCmd).Result()
		result[fks[i]] = n > 0
	}

	return result
}

// 获取外键数据总数
func (ma *Indexes) Count(fk string) int64 {
	if !ma.Exist(fk) {
		indexes, err := ma.SetPks(fk)
		if err != nil {
			return 0
		}
		return int64(len(indexes))
	}

	mfk := ma.Get(fk)
	count, err := rds.Client().SCard(context.Background(), mfk.key).Result()
	if err != nil {
		log.Warn("[MFK] Get [%s] fk [%s] data count error: %s", ma.table.Name(), fk, err.Error())
		return 0
	}

	return count
}

// 获取多外键数据总数
func (ma *Indexes) MultiCount(fks []string) int64 {
	// 检查多外键缓存是否存在
	exists := ma.MultiExist(fks)
	var not []string
	for fk, ok := range exists {
		if !ok {
			not = append(not, fk)
		}
	}
	if len(not) > 0 {
		_, err := ma.MultiSetPks(not)
		if err != nil {
			log.Warn("[MFK][Multi] Set [%s] multi fk [%+v] data error: %s", ma.table.Name(), not, err.Error())
			return 0
		}
	}

	// 获取每个外键集合的总数
	ctx := context.Background()
	cmds, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, fk := range fks {
			pipe.SCard(ctx, ma.Get(fk).key)
		}
		return nil
	})
	if err != nil {
		log.Warn("[MFK][Multi] Get [%s] multi fk [%+v] data count error: %s", ma.table.Name(), fks, err.Error())
		return 0
	}

	var count int64
	for i, cmd := range cmds {
		if c, ok := cmd.(*redis.IntCmd); ok {
			val, err := c.Result()
			if err != nil {
				log.Warn("[MFK][Multi] Get [%s] multi fk [%s] data count error: %s", ma.table.Name(), fks[i], err.Error())
				continue
			}
			count += val
		}
	}

	return count
}

// 计算分页 分页结果, 偏移量
func (ma *Indexes) Scan(fk string, cur int64, limit int64) (*ScanResult, int64) {
	count := ma.Count(fk)

	res := NewScanResult(cur, count, limit)
	offset := (cur - 1) * limit

	return res, offset
}

// 计算分页 分页结果, 偏移量
func (ma *Indexes) MultiScan(fks []string, cur int64, limit int64) (*ScanResult, int64) {
	count := ma.MultiCount(fks)

	res := NewScanResult(cur, count, limit)
	offset := (cur - 1) * limit

	return res, offset
}

// 获取外键的主键列表
func (ma *Indexes) GetPks(fk string) ([]string, error) {
	if !ma.Exist(fk) {
		return ma.SetPks(fk)
	}

	// 从缓存读取
	ctx := context.Background()
	indexes := ma.Get(fk)
	cmds, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SMembers(ctx, indexes.key)
		pipe.Expire(ctx, indexes.key, ma.table.admin.opts.IndexExpireTime)
		return nil
	})
	if err != nil {
		return nil, err
	}

	pks, err := cmds[0].(*redis.StringSliceCmd).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	return pks, nil
}

// 获取多个外键的主键列表
func (ma *Indexes) MultiGetPks(fks []string) (map[string][]string, error) {
	exists := ma.MultiExist(fks)
	var not []string
	var in []string
	for fk, ok := range exists {
		if !ok {
			not = append(not, fk)
		} else {
			in = append(in, fk)
		}
	}

	var result = make(map[string][]string, len(fks))
	if len(not) > 0 {
		res, err := ma.MultiSetPks(not)
		if err != nil {
			return nil, err
		}
		for k, v := range res {
			result[k] = append(result[k], v...)
		}
	}

	// 从缓存读取已存在的
	ctx := context.Background()
	mfks := ma.MultiGet(fks)
	cmds, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		for _, mfk := range mfks {
			pipe.SMembers(ctx, mfk.key)
			pipe.Expire(ctx, mfk.key, ma.table.admin.opts.IndexExpireTime)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	if len(cmds)%2 != 0 {
		return nil, ErrRedisError
	}

	// 读取数据
	for i := 0; i < len(cmds); i += 2 {
		if c, ok := cmds[i].(*redis.StringSliceCmd); ok {
			res, err := c.Result()
			if err != nil {
				log.Debug("[MFK][Multi] load cache data error: %s", err.Error())
				continue
			}
			fk := fks[i/2]
			result[fk] = append(result[fk], res...)
		}
	}

	return result, nil
}

// 设置外键的主键列表缓存
func (ma *Indexes) SetPks(fk string) ([]string, error) {
	// 从数据库读取主键列表
	pks, err := ma.loadFkData(fk)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := ma.setCache([]string{fk}, map[string][]string{fk: pks}); err != nil {
		return nil, err
	}

	return pks, nil
}

// 设置多个外键的主键列表缓存
func (ma *Indexes) MultiSetPks(fks []string) (map[string][]string, error) {
	pks, err := ma.loadMultiFkData(fks)
	if err != nil {
		return nil, err
	}

	// 写入缓存
	if err := ma.setCache(fks, pks); err != nil {
		return nil, err
	}

	return pks, nil
}

// 增加主键
func (ma *Indexes) AddPk(fk string, pks ...string) error {
	if !ma.Exist(fk) {
		if _, err := ma.SetPks(fk); err != nil {
			return err
		}
	}
	return ma.addCache(fk, pks...)
}

// 删除主键
func (ma *Indexes) RemPk(fk string, pks ...string) error {
	if !ma.Exist(fk) {
		if _, err := ma.SetPks(fk); err != nil {
			return err
		}
	}
	return ma.delCache(fk, pks...)
}

// 从数据库读取外键的主键列表数据
func (ma *Indexes) loadFkData(fk string) ([]string, error) {
	// 从数据库读取主键列表
	ctx := context.TODO()
	filter := bson.M{
		ma.table.FkField(): FormatIndexVal(ma.table.FkKind(), fk),
	}
	opts := options.Find().SetProjection(bson.M{ma.table.PkField(): 1})
	cur, err := mgo.Col(ma.table.Name(), ma.table.admin.dbName).Find(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer cur.Close(ctx)

	var rows []map[string]interface{}
	if err := cur.All(ctx, &rows); err != nil {
		return nil, err
	}

	var pks []string
	for _, row := range rows {
		pks = append(pks, dtype.ParseStr(row[ma.table.PkField()]))
	}

	return pks, nil
}

// 从数据库读取外键的主键列表数据 (多外键)
func (ma *Indexes) loadMultiFkData(fks []string) (map[string][]string, error) {
	if ma.table.FkField() == "" {
		return nil, ErrInvalidFilter
	}

	var fFks []interface{}
	for _, fk := range fks {
		fFks = append(fFks, FormatIndexVal(ma.table.FkKind(), fk))
	}

	// 从数据库读取主键列表
	ctx := context.TODO()
	filter := bson.M{
		ma.table.FkField(): bson.M{"$in": fFks},
	}
	opts := options.Find().SetProjection(bson.M{ma.table.PkField(): 1, ma.table.FkField(): 1})
	cur, err := mgo.Col(ma.table.Name(), ma.table.admin.dbName).Find(ctx, filter, opts)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}
	defer cur.Close(ctx)

	var rows []map[string]interface{}
	if err := cur.All(ctx, &rows); err != nil {
		return nil, err
	}

	var pks = make(map[string][]string)
	for _, row := range rows {
		fk := dtype.ParseStr(row[ma.table.FkField()])
		pk := dtype.ParseStr(row[ma.table.PkField()])
		pks[fk] = append(pks[fk], pk)
	}

	return pks, nil
}

// 写入新缓存
func (ma *Indexes) setCache(fks []string, pks map[string][]string) error {
	if len(pks) == 0 {
		return nil
	}

	ctx := context.Background()
	mfks := ma.MultiGet(fks)
	rsFn := func(pipe redis.Pipeliner) error {
		for _, mfk := range mfks {
			pipe.Del(ctx, mfk.key)
			if len(pks[mfk.fk]) > 0 {
				pipe.SAdd(ctx, mfk.key, toSliceInterface(pks[mfk.fk])...)
				pipe.Expire(ctx, mfk.key, ma.table.admin.opts.IndexExpireTime)
			}
		}
		return nil
	}
	if _, err := rds.Client().Pipelined(ctx, rsFn); err != nil {
		return err
	}

	return nil
}

// 获取
func (ma *Indexes) allCache(fk string) ([]string, error) {
	ctx := context.Background()
	mfk := ma.Get(fk)
	cmds, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SMembers(ctx, mfk.key)
		pipe.Expire(ctx, mfk.key, ma.table.admin.opts.IndexExpireTime)
		return nil
	})
	if err != nil {
		return nil, err
	}

	pks, err := cmds[0].(*redis.StringSliceCmd).Result()
	if err != nil && err != redis.Nil {
		return nil, err
	}

	return pks, nil
}

// 增加
func (ma *Indexes) addCache(fk string, pks ...string) error {
	if len(pks) == 0 {
		return nil
	}

	ctx := context.Background()
	mfk := ma.Get(fk)
	_, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SAdd(ctx, mfk.key, toSliceInterface(pks)...)
		pipe.Expire(ctx, mfk.key, ma.table.admin.opts.IndexExpireTime)
		return nil
	})
	return err
}

// 删除
func (ma *Indexes) delCache(fk string, pks ...string) error {
	if len(pks) == 0 {
		return nil
	}

	ctx := context.Background()
	mfk := ma.Get(fk)
	_, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.SRem(ctx, mfk.key, toSliceInterface(pks)...)
		pipe.Expire(ctx, mfk.key, ma.table.admin.opts.IndexExpireTime)
		return nil
	})
	return err
}

// 实例化数据表外键管理器
func newIndexes(table *Table) *Indexes {
	return &Indexes{
		table:   table,
		indexes: make(map[string]*Index),
	}
}

func toSliceInterface(slice []string) []interface{} {
	var values = make([]interface{}, 0, len(slice))
	for _, str := range slice {
		values = append(values, str)
	}
	return values
}
