package mem

import (
	"context"
	"encoding/json"
	rds "micro-libs/store/redis"
	"micro-libs/utils/dtype"
	"micro-libs/utils/log"
	"micro-libs/utils/tool"
	"reflect"
	"sync"
	"time"
)

type ModelLocal struct {
	sync.RWMutex
	table  *Table        // 数据表
	watch  *Watch        // 数据变化监听
	pk     string        // 主键值
	fk     string        // 外键值
	ref    reflect.Value // 数据反射
	data   interface{}   // 数据结构
	update bool          // 是否有更新
	active time.Time     // 活跃时间
}

func (mm *ModelLocal) Table() *Table {
	return mm.table
}

func (mm *ModelLocal) Watch() *Watch {
	return mm.watch
}

func (mm *ModelLocal) Pk() string {
	return mm.pk
}

func (mm *ModelLocal) Fk() string {
	return mm.fk
}

func (mm *ModelLocal) Index() string {
	return GetCacheIndex(mm.pk, mm.fk)
}

func (mm *ModelLocal) Ref() reflect.Value {
	return mm.ref
}

func (mm *ModelLocal) Data() interface{} {
	return mm.data
}

func (mm *ModelLocal) Byte() []byte {
	if mm.data == nil {
		return nil
	}
	b, _ := json.Marshal(mm.data)
	return b
}

func (mm *ModelLocal) GetField(field string) reflect.Value {
	mm.active = time.Now()
	return mm.ref.FieldByName(tool.UnderscoreToCamelCase(field))
}

func (mm *ModelLocal) GetValue(field string) interface{} {
	vf := mm.GetField(field)
	if !vf.IsValid() {
		return nil
	}

	switch vf.Kind() {
	case reflect.Ptr, reflect.Map, reflect.Slice, reflect.Interface, reflect.UnsafePointer:
		if vf.IsNil() {
			return nil
		}
	}

	return vf.Interface()
}

// 设置单个属性值
func (mm *ModelLocal) SetValue(field string, value interface{}) {
	if field == mm.table.PkField() {
		return
	}

	vf := mm.GetField(field)
	if !vf.IsValid() || !vf.CanSet() {
		log.Error("[MModel][%s] Field [%s] Invalid or not can set ...", mm.table.Name(), field)
		return
	}

	vf.Set(reflect.ValueOf(value))
	mm.watch.Trigger(field)

	mm.update = true
}

// 设置多个属性值
func (mm *ModelLocal) SetValues(values map[string]interface{}) {
	for field, value := range values {
		mm.SetValue(field, value)
	}
}

// 设置单个属性值
func (mm *ModelLocal) SetRemoteValue(field string, value []byte) {
	if field == mm.table.PkField() {
		return
	}

	vf := mm.GetField(field)
	if !vf.IsValid() || !vf.CanSet() {
		log.Error("[MModel][%s] Field [%s] Invalid or not can set ...", mm.table.Name(), field)
		return
	}

	if vf.Kind() == reflect.Ptr {
		_ = json.Unmarshal(value, vf.Interface())
	} else {
		_ = json.Unmarshal(value, vf.Addr().Interface())
	}

	mm.update = true
	mm.watch.Trigger(field)
}

// 设置多个属性值
func (mm *ModelLocal) SetRemoteValues(values map[string][]byte) {
	for field, value := range values {
		mm.SetRemoteValue(field, value)
	}
}

// 导出缓存对象
func (mm *ModelLocal) MCache() *MCache {
	return NewMCache(mm.table, mm)
}

// 检查数据状态
func (mm *ModelLocal) CheckState() error {
	mm.Lock()
	defer mm.Unlock()

	// 如果数据过期了, 则立即同步到缓存, 并从内存中删除
	if time.Now().After(mm.active.Add(mm.table.admin.opts.StateExpireTime)) {
		return mm.flush()
	}

	return mm.sync()
}

// 刷新数据缓存
func (mm *ModelLocal) flush() error {
	// 关闭数据监听
	mm.watch.Close()

	// 立即同步数据
	if err := mm.sync(); err != nil {
		return err
	}

	// 删除内存数据
	mm.table.DelModel(mm.Index())

	return nil
}

// 同步数据
func (mm *ModelLocal) sync() error {
	if !mm.update {
		return mm.MCache().SetActive()
	}

	if err := mm.save(); err != nil {
		return err
	}

	mm.update = false

	return nil
}

func (mm *ModelLocal) save() error {
	return mm.MCache().Save()
}

// Clean 清理内存&缓存数据
func (mm *ModelLocal) Clean() error {
	mm.Lock()
	defer mm.Unlock()

	// 关闭数据监听
	mm.watch.Close()

	// 删除内存数据
	mm.table.DelModel(mm.Index())

	return rds.Client().Del(context.Background(), GetCacheName(mm.Table(), mm.Index())).Err()
}

func NewModelLocal(table *Table, model interface{}) *ModelLocal {
	var data interface{}
	if b, ok := model.([]byte); ok {
		data = reflect.New(table.Model()).Interface()
		_ = json.Unmarshal(b, data)
	} else {
		data = model
	}

	mm := &ModelLocal{
		table:  table,
		ref:    reflect.Indirect(reflect.ValueOf(data)),
		data:   data,
		active: time.Now(),
	}

	// 事件监听
	mm.watch = NewWatch(mm)

	// 设置主键值
	mm.pk = dtype.ParseStr(mm.GetField(mm.table.PkField()).Interface())

	// 设置外键值
	if mm.table.FkField() != "" {
		mm.fk = dtype.ParseStr(mm.GetField(mm.table.FkField()).Interface())
	}

	return mm
}
