package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"micro-libs/utils/log"
	"micro-libs/utils/tool"
	"sync"
)

// 数据库初始化
type Tables struct {
	sync.RWMutex
	tables map[string]*Table
}

// Tables 返回所有数据表
func (mts *Tables) Tables() map[string]*Table {
	return mts.tables
}

// Count 返回集合数量
func (mts *Tables) Count() int {
	mts.RLock()
	defer mts.RUnlock()

	return len(mts.tables)
}

// Get 获取指定集合
func (mts *Tables) Get(name string) *Table {
	mts.RLock()
	defer mts.RUnlock()

	for n, col := range mts.tables {
		if n == name {
			return col
		}
	}

	return nil
}

// Add 添加新集合
func (mts *Tables) Add(name, pkField, fkField string, model interface{}, index []mongo.IndexModel, data []interface{}) *Table {
	tab := NewTable(name, pkField, fkField, model)
	if len(index) > 0 {
		tab.SetIndex(index)
	}
	if len(data) > 0 {
		tab.SetData(data)
	}

	mts.Lock()
	mts.tables[name] = tab
	mts.Unlock()

	return tab
}

// 设置自增ID初始数据
func (mts *Tables) SetAutoIdData(id string, initial int64) {
	mts.Add(GetAutoIncName(id), "", "", AutoIncId{}, nil, []interface{}{
		&AutoIncId{id, initial},
	})
}

// Init 初始化数据库
func (mts *Tables) Check(mdb *Store, dbname string) error {
	if mts.Count() == 0 {
		return nil
	}

	log.Debug("[%s] check collections ...", dbname)

	// 获取集合列表
	names, err := mdb.ListCollectionNames(dbname)
	if err != nil {
		return err
	}

	var tables []string
	mdb.Client().UseSession(context.Background(), func(sctx mongo.SessionContext) error {
		cdb := mdb.Client().Database(dbname)

		for _, tab := range mts.tables {
			// 判断集合是否已存在
			if !tool.InStrSlice(tab.name, names) {
				// 获取MongoDB的集合对象
				col := cdb.Collection(tab.name)

				// 创建索引
				if len(tab.index) > 0 {
					if _, err := col.Indexes().CreateMany(sctx, tab.Index()); err != nil {
						log.Error("create table [ %s ] index failure, error: %s", tab.name, err.Error())
					} else {
						log.Info("init table [ %s ] index success. total: %d", tab.name, len(tab.index))
					}
				}

				// 初始化数据
				if len(tab.data) > 0 {
					if _, err := col.InsertMany(sctx, tab.data); err != nil {
						log.Error("init table [ %s ] data failure, error: %s", tab.name, err.Error())
					} else {
						log.Info("init table [ %s ] data success. total: %d", tab.name, len(tab.data))
					}
				}

				tables = append(tables, tab.name)
			}
		}
		return nil
	})

	if len(tables) > 0 {
		log.Info("[%s] successfully initialize tables %+v ...", dbname, tables)
	}

	return nil
}

// 实例化数据库模型
func NewTables() *Tables {
	d := &Tables{
		tables: make(map[string]*Table),
	}
	return d
}
