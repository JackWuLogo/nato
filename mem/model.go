package mem

import (
	"micro-libs/utils/dtype"
	"reflect"
	"sync"
)

type MModel interface {
	RLock()                                   // 加读锁
	RUnlock()                                 // 释放读锁
	Lock()                                    // 加写锁
	Unlock()                                  // 释放写锁
	Table() *Table                            // 数据表
	Pk() string                               // 主键值
	Fk() string                               // 外键值
	Index() string                            // 缓存索引
	Ref() reflect.Value                       // 数据反射
	Data() interface{}                        // 原始数据
	Byte() []byte                             // 字节数据
	GetField(field string) reflect.Value      // 获取字段
	GetValue(field string) interface{}        // 获取值
	SetValue(field string, value interface{}) // 设置单个属性值
	SetValues(values map[string]interface{})  // 设置多个属性值
}

type ModelNil struct {
	sync.RWMutex
	table *Table // 数据表
}

func (mm *ModelNil) Table() *Table                            { return mm.table }
func (mm *ModelNil) Pk() string                               { return "" }
func (mm *ModelNil) Fk() string                               { return "" }
func (mm *ModelNil) Index() string                            { return "" }
func (mm *ModelNil) Ref() reflect.Value                       { return reflect.Value{} }
func (mm *ModelNil) Data() interface{}                        { return dtype.Ptr(mm.table.Model()).Interface() }
func (mm *ModelNil) Byte() []byte                             { return []byte{} }
func (mm *ModelNil) GetField(field string) reflect.Value      { return reflect.Value{} }
func (mm *ModelNil) GetValue(field string) interface{}        { return nil }
func (mm *ModelNil) SetValue(field string, value interface{}) {}
func (mm *ModelNil) SetValues(values map[string]interface{})  {}

// 空内存模型
func NewNilMModel(table *Table) *ModelNil {
	return &ModelNil{table: table}
}
