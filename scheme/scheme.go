package scheme

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"micro-libs/utils/errors"
	"micro-libs/utils/log"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"
)

// Scheme 配置表管理器
type Scheme struct {
	sync.RWMutex
	ver    *Version
	tables map[string]*Table // 配置表
}

func (s *Scheme) Ver() *Version {
	return s.ver
}

// Add 注册
func (s *Scheme) Add(key, name string, model interface{}) {
	s.Lock()
	defer s.Unlock()

	s.tables[key] = NewTable(s, key, name, model)
}

// Keys 获取所有配置表名
func (s *Scheme) Keys() []string {
	s.RLock()
	defer s.RUnlock()

	var keys = make([]string, 0, len(s.tables))
	for _, v := range s.tables {
		keys = append(keys, v.Key)
	}

	return keys
}

// Table 获取指定配置表
func (s *Scheme) Table(key string) *Table {
	s.RLock()
	defer s.RUnlock()

	if v, ok := s.tables[key]; ok {
		return v
	}

	return nilTable
}

// Tables 获取所有配置表
func (s *Scheme) Tables() map[string]*Table {
	s.RLock()
	defer s.RUnlock()

	var tables = make(map[string]*Table, len(s.tables))
	for k, v := range s.tables {
		tables[k] = v
	}

	return tables
}

// Sort 获取配置表
func (s *Scheme) Sort() []*Table {
	s.RLock()
	defer s.RUnlock()

	var tables = make(SortTable, 0, len(s.tables))
	for _, v := range s.tables {
		tables = append(tables, v)
	}

	sort.Sort(tables)

	return tables
}

// Values 获取全部配置项
func (s *Scheme) Values(key string) map[string]interface{} {
	return s.Table(key).Values()
}

// Value 获取单个配置项
func (s *Scheme) Value(key string, id interface{}) interface{} {
	return s.Table(key).Get(id)
}

// Range 迭代配置项
func (s *Scheme) Range(key string, fn func(k string, v interface{}) bool) {
	s.Table(key).Range(fn)
}

// Update 更新配置数据
func (s *Scheme) Update(values map[string]string, aid int64) (tables []string, errs []string) {
	if len(values) == 0 {
		return nil, nil
	}

	var updates = make(map[string]reflect.Value)
	var nowTime = time.Now().Unix()

	// 解析数据
	for key, val := range values {
		key = strings.TrimSpace(strings.ToLower(key))
		tab := s.Table(key)
		if tab.IsNil() {
			errs = append(errs, fmt.Sprintf("[%s] 未生成数据结构, 请先生成数据结构后再导入数据 ...", key))
			continue
		}

		slice := tab.NewSlice()
		if err := json.Unmarshal([]byte(val), slice.Addr().Interface()); err != nil {
			errs = append(errs, fmt.Sprintf("[%s] 数据解析失败, 请检查是否有误: %s ...", key, err.Error()))
			continue
		}

		updates[key] = slice
	}

	// 更新数据
	_ = s.ver.mongo.Client().UseSession(context.TODO(), func(sctx mongo.SessionContext) error {
		db := sctx.Client().Database(s.ver.dbName)
		vCol := db.Collection(s.ver.verName)

		for key, val := range updates {
			tab := s.Table(key)

			// 检查key是否存在
			var origin = new(VersionModel)
			if err := vCol.FindOne(sctx, bson.M{"key": tab.Key}).Decode(origin); err != nil && err != mongo.ErrNoDocuments {
				errs = append(errs, fmt.Sprintf("[%s] 读取版本信息失败: %s", tab.Key, err.Error()))
				continue
			}

			if origin.Id.IsZero() {
				// 写入新配置数据
				insert := &VersionModel{
					Id:         primitive.NewObjectID(),
					Key:        tab.Key,
					Name:       tab.Name,
					Total:      val.Len(),
					Version:    1,
					UpdateUser: aid,
					UpdateTime: nowTime,
					CreateUser: aid,
					CreateTime: nowTime,
				}
				if _, err := vCol.InsertOne(sctx, insert); err != nil {
					errs = append(errs, fmt.Sprintf("[%s] 写入版本信息失败: %s", tab.Key, err.Error()))
					continue
				}
			} else {
				// 更新旧配置数据
				update := bson.M{
					"$set": bson.M{
						"key":         tab.Key,
						"name":        tab.Name,
						"total":       val.Len(),
						"update_user": aid,
						"update_time": nowTime,
					},
					"$inc": bson.M{
						"version": int64(1),
					},
				}
				if _, err := vCol.UpdateOne(sctx, bson.M{"_id": origin.Id}, update); err != nil {
					errs = append(errs, fmt.Sprintf("[%s] 更新版本信息失败: %s", tab.Key, err.Error()))
					continue
				}
			}

			// 准备写入配置表数据
			dCol := db.Collection(tab.Key)

			_ = dCol.Drop(sctx) // 先清空旧数据

			// 写入新数据
			var rows []interface{}
			for i := 0; i < val.Len(); i++ {
				rows = append(rows, val.Index(i).Interface())
			}
			if _, err := dCol.InsertMany(sctx, rows); err != nil {
				errs = append(errs, fmt.Sprintf("[%s] 写入配置数据失败: %s", tab.Key, err.Error()))
				continue
			}

			log.Debug("[%s] scheme update success, total: %d ...", key, len(rows))

			tables = append(tables, tab.Key)
		}

		return nil
	})

	return tables, errs
}

// Load 从数据库加载数据
func (s *Scheme) Load(keys ...string) error {
	if len(keys) == 0 {
		keys = s.Keys()
	}

	ctx := context.TODO()

	// 获取版本信息
	version, err := s.ver.GetCache(ctx)
	if err != nil {
		return err
	}

	var valid []string
	for _, key := range keys {
		if _, ok := version[key]; ok {
			valid = append(valid, key)
		} else {
			log.Warn("[%s] not found scheme table version info ...", key)
		}
	}

	fn := func(sctx mongo.SessionContext) error {
		db := sctx.Client().Database(s.ver.dbName)

		for _, key := range valid {
			tab, ok := s.tables[key]
			if !ok {
				log.Error("[%s] not found scheme table register info ...", key)
				continue
			}

			cur, err := db.Collection(key).Find(sctx, bson.M{})
			if err != nil {
				if err == mongo.ErrNilDocument {
					continue
				}
				log.Error("[%s] load mongodb error: %s", key, err.Error())
				continue
			}

			// 解析数据
			var rows = tab.NewSlice()
			if err := cur.All(sctx, rows.Addr().Interface()); err != nil {
				log.Error("[%s] parse mongodb data error: %s", key, err.Error())
				cur.Close(sctx)
				continue
			}
			cur.Close(sctx)

			if err := tab.SetValues(version[key], rows); err != nil {
				log.Error("[%s] update scheme table data error: %s", key, err.Error())
				continue
			}

			log.Debug("[%s] scheme form mongodb load success, version: %d, total: %d", key, version[key], rows.Len())
		}

		return nil
	}

	// 读取配置信息
	if err := s.ver.mongo.Client().UseSession(ctx, fn); err != nil {
		return err
	}

	return nil
}

// Delete 删除数据
func (s *Scheme) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		keys = s.Keys()
	}

	// 先删除缓存
	if err := s.ver.DelCache(ctx, keys...); err != nil {
		return err
	}

	// 再删除数据库
	delFn := func(sctx mongo.SessionContext) error {
		db := sctx.Client().Database(s.ver.dbName)
		verCol := db.Collection(s.ver.verName)

		// 删除版本库
		if _, err := verCol.DeleteMany(sctx, bson.M{"key": bson.M{"$in": keys}}); err != nil {
			return err
		}

		// 删除配置表
		for _, key := range keys {
			_ = db.Collection(key).Drop(sctx)
		}

		return nil
	}
	if err := s.ver.mongo.Client().UseSession(ctx, delFn); err != nil {
		return errors.Wrap(err, "删除配置表失败")
	}

	return nil
}

// NewScheme 实例化配置表
func NewScheme(opts ...Option) *Scheme {
	s := &Scheme{
		ver:    newVersion(),
		tables: make(map[string]*Table),
	}
	s.ver.Init(opts...)
	return s
}
