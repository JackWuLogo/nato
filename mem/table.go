package mem

import (
	"context"
	"encoding/json"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"micro-libs/meta"
	mgo "micro-libs/store/mongo"
	rds "micro-libs/store/redis"
	"micro-libs/utils/dtype"
	"micro-libs/utils/errors"
	"micro-libs/utils/log"
	"micro-libs/utils/multi"
	"reflect"
	"sync"
)

// Table 数据表管理器
type Table struct {
	*mgo.Table
	sync.RWMutex
	wg      sync.WaitGroup
	admin   *Admin                 // 数据管理器
	indexes *Indexes               // 外键索引
	models  map[string]*ModelLocal // 数据记录
}

// 管理器
func (t *Table) Admin() *Admin {
	return t.admin
}

// 外键管理器
func (t *Table) Indexes() *Indexes {
	return t.indexes
}

// HasModel 检查内存数据是否存在
func (t *Table) HasModel(index string) bool {
	t.RLock()
	defer t.RUnlock()

	_, ok := t.models[index]

	return ok
}

// Has 获取内存数据
func (t *Table) GetModel(index string) *ModelLocal {
	t.RLock()
	defer t.RUnlock()

	mm, ok := t.models[index]
	if !ok {
		return nil
	}

	return mm
}

// 删除内存模型
func (t *Table) DelModel(indexes ...string) {
	t.Lock()
	for _, index := range indexes {
		delete(t.models, index)
	}
	t.Unlock()
}

// Get 获取主键数据
// 如果pk为空, 表示获取无外键的主键数据
// 如果pk不为空, 表示获取有外键的主键数据
func (t *Table) Get(mt meta.DataNodeMeta, pk ...string) (MModel, error) {
	var mtPk, mtFk string
	if len(pk) > 0 && pk[0] != "" {
		if pk[0] == "0" {
			return nil, errors.Invalid("the primary key of foreign key data cannot be zero: %s", pk[0])
		}
		mtPk = pk[0]
		mtFk = mt.DataId()
	} else {
		mtPk = mt.DataId()
		mtFk = ""
	}

	index := GetCacheIndex(mtPk, mtFk)

	// 生成数据索引
	t.RLock()
	if mm, ok := t.models[index]; ok {
		t.RUnlock()
		return mm, nil
	}
	t.RUnlock()

	// 如果数据是本节点的, 则本地加载. 否则远程加载
	if meta.IsSelf(mt) {
		// 缓存查找数据
		ctx := context.TODO()
		mc, err := rds.Client().HGetAll(ctx, GetCacheName(t, index)).Result()
		if err != nil && err != redis.Nil {
			return nil, err
		}

		cache := ToMCache(t, mc)

		var mm *ModelLocal
		if cache.IsValid() {
			// 缓存有效
			if cache.IsNil() {
				log.Debug("[LoadData] form redis load [%s][%s][%s] data error: not found", t.Name(), mtPk, mtFk)
				return nil, ErrDataNotFound
			}

			if len(cache.Data()) > 0 {
				mm = NewModelLocal(t, cache.Data())
			}
		} else {
			// 没有数据, 需要从数据库读取
			var filter = bson.M{
				t.PkField(): FormatIndexVal(t.PkKind(), mtPk),
			}
			if t.FkField() != "" {
				filter[t.FkField()] = FormatIndexVal(t.FkKind(), mtFk)
			}
			var res = reflect.New(t.Model()).Interface()
			if err := mgo.Col(t.Name(), t.Admin().DbName()).FindOne(ctx, filter).Decode(res); err != nil {
				if err == mongo.ErrNoDocuments {
					if err := FromMCache(t, mtPk, mtFk, nil).Save(); err != nil {
						return nil, err
					}
					log.Debug("[LoadData] form redis load [%s][%s][%s] data error: not found", t.Name(), mtPk, mtFk)
					return nil, ErrDataNotFound
				}
				log.Debug("[LoadData] form database load [%s][%s] data error: %s", t.Name(), mtPk, err)
				return nil, err
			}

			mm = NewModelLocal(t, res)

			// 设置缓存
			if err := mm.save(); err != nil {
				return nil, err
			}
		}

		// 注册模型
		if mm != nil {
			t.Lock()
			t.models[index] = mm
			t.Unlock()

			return mm, nil
		}

		return nil, ErrDataNotFound
	} else {
		// 远程节点
		res, err := t.admin.client.Get(mt, &InMemGet{Table: t.Name(), Pk: pk})
		if err != nil {
			log.Error("[Get][Remote] get [%s][%s] data error: %s", t.Name(), mt.DataId(), err)
			return nil, err
		}
		if len(res.Result) == 0 {
			return nil, ErrDataNotFound
		}
		return NewModelRemote(t, mt, res.Result), nil
	}
}

// GetFk 获取外键数据
// 如果指定了主键, 则使用指定主键列表
func (t *Table) GetFk(mt meta.DataNodeMeta, pk ...string) (map[string]MModel, error) {
	fk := mt.DataId()
	if fk == "" {
		return nil, errors.Invalid("invalid data meta, fk is nil")
	}

	// 过滤主键
	var pks []string
	if len(pk) > 0 {
		pks = pk
	} else {
		// 获取外键的主键列表
		res, err := t.indexes.GetPks(fk)
		if err != nil {
			return nil, err
		} else if len(res) == 0 {
			return nil, nil
		}
		pks = res
	}

	// 生成索引列表
	var indexes = make([]string, 0, len(pks))
	var indexPk = make(map[string]interface{}, len(pks))
	for _, pk := range pks {
		index := GetCacheIndex(pk, fk)
		indexes = append(indexes, index)
		indexPk[index] = FormatIndexVal(t.PkKind(), pk)
	}

	if meta.IsSelf(mt) {
		// 本节点
		var result = make(map[string]MModel, len(pks))

		// 从内存获取数据
		var no1 []string

		t.RLock()
		for _, index := range indexes {
			if mm, ok := t.models[index]; ok {
				result[mm.Pk()] = mm
			} else {
				no1 = append(no1, index)
			}
		}
		t.RUnlock()

		// 从缓存获取数据
		if len(no1) > 0 {
			ctx := context.Background()
			cmds, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
				for _, index := range no1 {
					pipe.HGetAll(ctx, GetCacheName(t, index))
				}
				return nil
			})
			if err != nil {
				return nil, err
			}

			var no2 []string //缓存也不存在的数据
			for i, index := range no1 {
				res, err := cmds[i].(*redis.StringStringMapCmd).Result()
				if err != nil && err != redis.Nil {
					return nil, err
				}

				cache := ToMCache(t, res)
				if cache.IsValid() && !cache.IsNil() && len(cache.Data()) > 0 {
					mm := NewModelLocal(t, cache.Data())

					result[mm.Pk()] = mm

					t.Lock()
					t.models[index] = mm // 注册模型
					t.Unlock()

					continue
				}

				no2 = append(no2, index)
			}

			// 从数据库获取剩余数据
			if len(no2) > 0 {
				var residue []interface{}
				for _, index := range no2 {
					residue = append(residue, indexPk[index])
				}

				rows, err := mgo.SelectAll(mgo.Col(t.Name(), t.admin.dbName), bson.M{t.PkField(): bson.M{"$in": residue}}, t.Model())
				if err != nil && err != mongo.ErrNoDocuments {
					return nil, err
				}

				for _, row := range rows {
					mm := NewModelLocal(t, row)
					if err := mm.save(); err != nil {
						return nil, err
					}

					result[mm.Pk()] = mm

					t.Lock()
					t.models[mm.Index()] = mm // 注册模型
					t.Unlock()
				}
			}
		}

		return result, nil
	} else {
		// 远程节点
		res, err := t.admin.client.GetFk(mt, &InMemGetFk{Table: t.Name(), Pk: pk})
		if err != nil {
			log.Error("[GetFk][Remote] get [%s][%s] fk data error: %s", t.Data(), mt.DataId(), err)
			return nil, err
		}

		var result = make(map[string]MModel)
		for key, b := range res.Result {
			result[key] = NewModelRemote(t, mt, b)
		}

		return result, nil
	}
}

// 获取外键数据总数
func (t *Table) GetFkCount(fk string) int64 {
	return t.indexes.Count(fk)
}

// 获取多外键数据总数
func (t *Table) GetFkMultiCount(fks []string) int64 {
	return t.indexes.MultiCount(fks)
}

// GetMulti 获取多条无外键的主键数据
func (t *Table) GetMulti(nodes []meta.DataNodeMeta) (map[string]MModel, error) {
	work := multi.NewWorks()
	for _, node := range nodes {
		node := node
		work.Do(func() (interface{}, error) {
			return t.Get(node)
		})
	}

	if err := work.Run(); err != nil {
		return nil, err
	}

	var result = make(map[string]MModel)
	for _, row := range work.Result() {
		if row == nil {
			continue
		}

		mm := row.(MModel)
		result[mm.Pk()] = mm
	}

	return result, nil
}

// Insert 写入新数据
func (t *Table) Insert(mt meta.DataNodeMeta, model interface{}, pk ...string) (MModel, error) {
	var data interface{}
	if b, ok := model.([]byte); ok {
		data = reflect.New(t.Model()).Interface()
		_ = json.Unmarshal(b, data)
	} else {
		data = model
	}

	// 创建数据模型
	if meta.IsSelf(mt) {
		// 数据库写入数据
		if _, err := mgo.Col(t.Name(), t.admin.dbName).InsertOne(nil, data); err != nil {
			return nil, err
		}

		var newPk, newFk string
		if len(pk) > 0 {
			if pk[0] == mt.DataId() || pk[0] == "" {
				newPk = mt.DataId()
				newFk = ""
			} else {
				newPk = pk[0]
				newFk = mt.DataId()
			}
		} else {
			newPk = mt.DataId()
			newFk = ""
		}

		// 本地数据mData
		mm := NewModelLocal(t, data)
		// 写入缓存数据
		if err := mm.save(); err != nil {
			return nil, err
		}

		t.Lock()
		t.models[GetCacheIndex(newPk, newFk)] = mm // 注册模型
		t.Unlock()

		// 记录外键信息
		if t.FkField() != "" {
			if err := t.indexes.AddPk(newFk, newPk); err != nil {
				log.Error("[Insert] add [%s] fk [%s] index [%s] error: %s", t.Name(), mm.Pk(), mm.Fk(), err.Error())
			}
		}

		return mm, nil
	} else {
		b, _ := json.Marshal(data)
		if _, err := t.admin.client.Insert(mt, &InMemInsert{Table: t.Name(), Data: b, Pk: pk}); err != nil {
			log.Error("[Insert][Remote] insert [%s] data error: %s", t.Name(), err)
			return nil, err
		}
		return NewModelRemote(t, mt, data), nil
	}
}

// 从数据库删除数据模型 (数据库立即删除)
// 如果pk为空, 表示删除无外键的主键数据
// 如果pk不为空, 表示删除有外键的主键数据
func (t *Table) Delete(mt meta.DataNodeMeta, pk ...string) error {
	// 本节点
	if meta.IsSelf(mt) {
		// 生成数据索引
		var mtPk, mtFk string
		if len(pk) > 0 && pk[0] != "" {
			mtPk = pk[0]
			mtFk = mt.DataId()
		} else {
			mtPk = mt.DataId()
			mtFk = ""
		}

		index := GetCacheIndex(mtPk, mtFk)

		t.RLock()
		mm, ok := t.models[index]
		t.RUnlock()
		if ok {
			// 清理内存&缓存
			if err := mm.Clean(); err != nil {
				return err
			}
		}

		// 删除数据库记录
		if _, err := mgo.Col(t.Name(), t.admin.dbName).DeleteOne(nil, bson.M{t.PkField(): FormatIndexVal(t.PkKind(), mtPk)}); err != nil {
			return err
		}

		// 删除外键信息
		if mtFk != "" {
			if err := t.indexes.RemPk(mtFk, mtPk); err != nil {
				log.Error("[Delete] delete [%s] fk [%s] index [%s] error: %s", t.Name(), mtPk, mtFk, err.Error())
			}
		}
	} else {
		// 远程节点
		if _, err := t.admin.client.Delete(mt, &InMemDelete{Table: t.Name(), Pk: pk}); err != nil {
			log.Error("[Delete][Remote] delete [%s] %v error: %s", t.Name(), pk, err)
			return err
		}
	}

	return nil
}

// 筛选数据主键 (可排序, 直接查询数据库获取主键列表. 慎用)
func (t *Table) Filter(filter bson.M, sort bson.D, limit int64, skip int64) ([]string, error) {
	opts := options.Find().SetProjection(bson.M{t.PkField(): 1})
	if len(sort) > 0 {
		opts.SetSort(sort)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}

	// 默认数据读取限制
	if opts.Limit == nil || *opts.Limit == 0 {
		opts.SetLimit(20)
	}

	ctx := context.TODO()
	cur, err := mgo.Col(t.Name(), t.admin.dbName).Find(ctx, filter, opts)
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
		pks = append(pks, dtype.ParseStr(row[t.PkField()]))
	}

	return pks, nil
}

// 筛选数据主键和外键 (可排序, 直接查询数据库获取主键&外键列表. 慎用), map[FK][]PK
func (t *Table) FilterFks(filter bson.M, sort bson.D, limit int64, skip int64) ([]*PkFkData, error) {
	if t.FkField() == "" {
		return nil, ErrInvalidFilter
	}

	opts := options.Find().SetProjection(bson.M{t.PkField(): 1, t.FkField(): 1})
	if len(sort) > 0 {
		opts.SetSort(sort)
	}
	if limit > 0 {
		opts.SetLimit(limit)
	}
	if skip > 0 {
		opts.SetSkip(skip)
	}

	// 默认数据读取限制
	if opts.Limit == nil || *opts.Limit == 0 {
		opts.SetLimit(20)
	}

	ctx := context.TODO()
	cur, err := mgo.Col(t.Name(), t.admin.dbName).Find(ctx, filter, opts)
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

	// 整理数据
	var pks []*PkFkData
	for _, row := range rows {
		pk := dtype.ParseStr(row[t.PkField()])
		fk := dtype.ParseStr(row[t.FkField()])
		pks = append(pks, &PkFkData{
			Pk: pk,
			Fk: fk,
		})
	}

	return pks, nil
}

// 外键数据分页
func (t *Table) FkScan(fk string, sort bson.D, cur int64, limit int64) (*ScanResult, error) {
	return t.FkScanFilter(fk, nil, sort, cur, limit)
}

// 外键数据分页 (多外键)
func (t *Table) MultiFkScan(fks []string, sort bson.D, cur int64, limit int64) (*ScanResult, error) {
	return t.MultiFkScanFilter(fks, nil, sort, cur, limit)
}

// 外键数据分页 (可过滤)
func (t *Table) FkScanFilter(fk string, filter bson.M, sort bson.D, cur int64, limit int64) (*ScanResult, error) {
	if filter == nil {
		filter = bson.M{}
	}
	filter[t.FkField()] = FormatIndexVal(t.FkKind(), fk)

	res, skip := t.indexes.Scan(fk, cur, limit)

	pks, err := t.Filter(filter, sort, limit, skip)
	if err != nil {
		return nil, err
	}

	res.Rows = pks

	return res, nil
}

// 外键数据分页 (多外键, 可过滤)
func (t *Table) MultiFkScanFilter(fks []string, filter bson.M, sort bson.D, cur int64, limit int64) (*ScanResult, error) {
	if filter == nil {
		filter = bson.M{}
	}
	var fmtFks []interface{}
	for _, fk := range fks {
		fmtFks = append(fmtFks, FormatIndexVal(t.FkKind(), fk))
	}
	filter[t.FkField()] = bson.M{"$in": fmtFks}

	res, skip := t.indexes.MultiScan(fks, cur, limit)

	pks, err := t.FilterFks(filter, sort, limit, skip)
	if err != nil {
		return nil, err
	}

	res.Rows = pks

	return res, nil
}

// SyncAll 同步所有数据到缓存
func (t *Table) SyncAll() {
	models := make(map[string]*ModelLocal, len(t.models))

	t.RLock()
	for k, v := range t.models {
		models[k] = v
	}
	t.RUnlock()

	for _, mm := range models {
		mm := mm
		t.wg.Add(1)

		go func() {
			defer t.wg.Done()

			// 检查数据状态
			if err := mm.CheckState(); err != nil {
				log.Error("[MemSync] sync table [%s][%s] data error: %s", t.Name(), mm.Index(), err)
			}
		}()
	}

	t.wg.Wait()
}

// NewTable 实例化数据表管理器
func NewTable(admin *Admin, table *mgo.Table) *Table {
	tab := &Table{
		Table:  table,
		admin:  admin,
		models: make(map[string]*ModelLocal),
	}
	tab.indexes = newIndexes(tab)
	return tab
}
