package mem

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	rds "micro-libs/store/redis"
	"micro-libs/utils/dtype"
	"time"
)

// MCache 同步缓存数据字段
const (
	MCacheFieldPk    = "pk"
	MCacheFieldFk    = "fk"
	MCacheFieldIsNil = "is_nil"
	MCacheFieldData  = "data"
	MCacheFieldSync  = "sync"
	MCacheFieldTime  = "time"
)

// 缓存数据结构
type MCache struct {
	table *Table
	cache map[string]string
}

func (m *MCache) IsValid() bool {
	return len(m.cache) > 0 && m.cache[MCacheFieldPk] != ""
}

func (m *MCache) IsNil() bool {
	return dtype.ParseBool(m.cache[MCacheFieldIsNil])
}

func (m *MCache) Index() string {
	return GetCacheIndex(m.Pk(), m.Fk())
}

func (m *MCache) Pk() string {
	return m.cache[MCacheFieldPk]
}

func (m *MCache) Fk() string {
	return m.cache[MCacheFieldFk]
}

func (m *MCache) Data() []byte {
	return []byte(m.cache[MCacheFieldData])
}

func (m *MCache) Sync() bool {
	return dtype.ParseBool(m.cache[MCacheFieldSync])
}

func (m *MCache) Time() int64 {
	return dtype.ParseInt64(m.cache[MCacheFieldTime])
}

func (m *MCache) Values() map[string]string {
	return m.cache
}

func (m *MCache) Parse(v interface{}) error {
	data := m.Data()
	if len(data) == 0 {
		return nil
	}
	return json.Unmarshal(data, v)
}

func (m *MCache) Set(mtPk, mtFk string, data []byte) {
	m.SetPk(mtPk)
	m.SetFk(mtFk)
	m.SetData(data)
	m.SetTime(time.Now().Unix())
}

func (m *MCache) SetPk(val string) {
	m.cache[MCacheFieldPk] = val
}

func (m *MCache) SetFk(val string) {
	m.cache[MCacheFieldFk] = val
}

func (m *MCache) SetData(val []byte) {
	if len(val) == 0 {
		m.cache[MCacheFieldIsNil] = "true"
	} else {
		m.cache[MCacheFieldIsNil] = "false"
	}
	m.cache[MCacheFieldData] = string(val)
}

func (m *MCache) SetSync(sync bool) {
	m.cache[MCacheFieldSync] = dtype.ParseStr(sync)
}

func (m *MCache) SetTime(ts int64) {
	m.cache[MCacheFieldTime] = dtype.ParseStr(ts)
}

func (m *MCache) SetActive() error {
	if m.Time() == 0 {
		m.SetTime(time.Now().Unix())
	}
	return rds.Client().HSet(context.TODO(), GetCacheName(m.table, m.Index()), MCacheFieldTime, m.Time()).Err()
}

func (m *MCache) Save() error {
	key := GetCacheName(m.table, m.Index())
	ctx := context.TODO()
	_, err := rds.Client().Pipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.HMSet(ctx, key, rds.Args{}.AddFlat(m.Values())...)
		if m.IsNil() {
			pipe.Expire(ctx, key, m.table.admin.Options().NilDataTTL)
		}
		return nil
	})
	return err
}

func (m *MCache) Del() error {
	return rds.Client().Del(context.TODO(), GetCacheName(m.table, m.Index())).Err()
}

// 生成缓存数据结构
func FromMCache(table *Table, mtPk, mtFk string, data []byte) *MCache {
	mc := &MCache{
		table: table,
		cache: make(map[string]string, 6),
	}
	mc.Set(mtPk, mtFk, data)
	mc.SetSync(false)
	return mc
}

// 新建缓存数据结构
func NewMCache(table *Table, mm *ModelLocal) *MCache {
	mc := &MCache{
		table: table,
		cache: make(map[string]string, 6),
	}
	mc.SetPk(mm.pk)
	mc.SetFk(mm.fk)
	mc.SetData(mm.Byte())
	mc.SetSync(mm.update)
	mc.SetTime(time.Now().Unix())
	return mc
}

// 生成缓存数据结构
func ToMCache(table *Table, cache map[string]string) *MCache {
	mc := &MCache{
		table: table,
		cache: cache,
	}
	return mc
}

// 获取数据索引
func GetCacheIndex(pk string, fk string) string {
	if fk == "" {
		fk = "NOFK"
	}
	return fmt.Sprintf("%s:%s", fk, pk)
}

// 生成缓存名称
func GetCacheName(tab *Table, index string) string {
	return rds.GetCacheName(tab.Admin().Options().PrefixData, tab.Name(), index)
}

// 格式化主键/外键值类型, 索引仅支持 string, int, int32, int64
func FormatIndexVal(kind string, val interface{}) interface{} {
	switch kind {
	case "int":
		return dtype.ParseInt(val)
	case "int32":
		return dtype.ParseInt32(val)
	case "int64":
		return dtype.ParseInt64(val)
	}
	return dtype.ParseStr(val)
}
