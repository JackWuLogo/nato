package scheme

import (
	"context"
	"encoding/json"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"micro-libs/utils/dtype"
	"micro-libs/utils/log"
	"sort"
	"sync"
)

// Clients 客户端配置表管理器
type Clients struct {
	sync.RWMutex
	ver    *Version
	values map[string]*Client // 导出数据 map<数据表key>
}

func (c *Clients) Ver() *Version {
	return c.ver
}

// Versions 获取版本库信息
func (c *Clients) Versions() map[string]int64 {
	c.RLock()
	defer c.RUnlock()

	var values = make(map[string]int64, len(c.values))
	for k, v := range c.values {
		values[k] = v.Version
	}

	return values
}

// Value 获取指定配置表
func (c *Clients) Value(key string) *Client {
	c.RLock()
	defer c.RUnlock()

	if v, ok := c.values[key]; ok {
		return v
	}

	return &Client{Table: key}
}

// Values 获取所有配置表
func (c *Clients) Values() map[string]*Client {
	c.RLock()
	defer c.RUnlock()

	var values = make(map[string]*Client, len(c.values))
	for k, v := range c.values {
		values[k] = v
	}
	return values
}

// Sort 获取排序配置表
func (c *Clients) Sort() []*Client {
	c.RLock()
	defer c.RUnlock()

	var values = make(SortClient, 0, len(c.values))
	for _, v := range c.values {
		values = append(values, v)
	}

	sort.Sort(values)

	return values
}

// Load 加载导出数据
func (c *Clients) Load(keys ...string) error {
	ctx := context.TODO()

	// 获取版本信息
	version, err := c.ver.GetCache(ctx)
	if err != nil {
		return err
	}

	c.Lock()
	defer c.Unlock()

	if len(keys) == 0 {
		for key := range version {
			keys = append(keys, key)
		}
		c.values = make(map[string]*Client)
	}

	fn := func(sctx mongo.SessionContext) error {
		db := sctx.Client().Database(c.ver.dbName)

		for _, key := range keys {
			ver, ok := version[key]
			if !ok {
				log.Error("[%s] not found scheme table version info ...", key)
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
			var rows []map[string]interface{}
			if err := cur.All(sctx, &rows); err != nil {
				log.Error("[%s] parse mongodb data error: %s", key, err.Error())
				cur.Close(sctx)
				continue
			}
			cur.Close(sctx)

			var res = make(map[string]map[string]interface{})
			for _, row := range rows {
				id, ok := row["_id"]
				if !ok {
					id = row["id"]
				}
				res[dtype.ParseStr(id)] = row
			}

			b, err := json.Marshal(res)
			if err != nil {
				log.Error("[%s] marshal error: %s", key, err.Error())
				continue
			}

			c.values[key] = &Client{
				Table:   key,
				Version: ver,
				Attrs:   b,
			}

			log.Debug("[%s] scheme form mongodb load success, version: %d", key, version[key])
		}

		return nil
	}

	// 读取配置信息
	if err := c.ver.mongo.Client().UseSession(ctx, fn); err != nil {
		return err
	}

	return nil
}

// NewClients 实例化配置表
func NewClients(opts ...Option) *Clients {
	s := &Clients{
		ver:    newVersion(),
		values: make(map[string]*Client),
	}
	s.ver.Init(opts...)
	return s
}
