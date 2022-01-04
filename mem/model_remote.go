package mem

import (
	"encoding/json"
	"micro-libs/meta"
	"micro-libs/utils/dtype"
	"micro-libs/utils/log"
	"micro-libs/utils/tool"
	"reflect"
	"sync"
)

// 远程数据对象
type ModelRemote struct {
	sync.RWMutex
	table *Table            // 数据表
	meta  meta.DataNodeMeta // 数据节点信息
	pk    string            // 主键值
	fk    string            // 外键值
	ref   reflect.Value     // 数据反射
	data  interface{}       // 数据结构
}

func (mm *ModelRemote) Table() *Table {
	return mm.table
}

func (mm *ModelRemote) Pk() string {
	return mm.pk
}

func (mm *ModelRemote) Fk() string {
	return mm.fk
}

func (mm *ModelRemote) Index() string {
	return GetCacheIndex(mm.pk, mm.fk)
}

func (mm *ModelRemote) Ref() reflect.Value {
	return mm.ref
}

func (mm *ModelRemote) Data() interface{} {
	return mm.data
}

func (mm *ModelRemote) Byte() []byte {
	if mm.data == nil {
		return nil
	}
	b, _ := json.Marshal(mm.data)
	return b
}

func (mm *ModelRemote) GetField(field string) reflect.Value {
	return mm.ref.FieldByName(tool.UnderscoreToCamelCase(field))
}

func (mm *ModelRemote) GetValue(field string) interface{} {
	vf := mm.GetField(field)
	if !vf.IsValid() {
		return nil
	}
	return vf.Interface()
}

// 设置属性值
func (mm *ModelRemote) SetValue(field string, value interface{}) {
	if field == mm.table.PkField() {
		return
	}

	vf := mm.GetField(field)
	if !vf.IsValid() || !vf.CanSet() {
		return
	}

	vf.Set(reflect.ValueOf(value))

	// 远程更新数据
	var pk []string
	if mm.table.FkField() != "" {
		pk = append(pk, mm.pk)
	}
	val, _ := json.Marshal(value)
	if _, err := mm.table.admin.client.SetValue(mm.meta, &InMemSetValue{
		Table: mm.table.Name(),
		Pk:    pk,
		Field: field,
		Value: val,
	}); err != nil {
		log.Error("[ModelRemote] Remote Set Value Failure: %s", err)
	}
}

// 批量设置属性值
func (mm *ModelRemote) SetValues(values map[string]interface{}) {
	for field, value := range values {
		if field == mm.table.PkField() {
			return
		}

		vf := mm.GetField(field)
		if !vf.IsValid() || !vf.CanSet() {
			return
		}

		vf.Set(reflect.ValueOf(value))
	}

	// 远程更新数据
	var pk []string
	if mm.table.FkField() != "" {
		pk = append(pk, mm.pk)
	}
	var data = make(map[string][]byte)
	for k, v := range values {
		if b, err := json.Marshal(v); err == nil {
			data[k] = b
		}
	}
	if _, err := mm.table.admin.client.SetValues(mm.meta, &InMemSetValues{
		Table:  mm.table.Name(),
		Pk:     pk,
		Values: data,
	}); err != nil {
		log.Error("[ModelRemote] Remote Set Values Failure: %s", err)
	}
}

func NewModelRemote(table *Table, mt meta.DataNodeMeta, model interface{}) *ModelRemote {
	var data interface{}
	if b, ok := model.([]byte); ok {
		data = reflect.New(table.Model()).Interface()
		_ = json.Unmarshal(b, data)
	} else {
		data = model
	}

	mm := &ModelRemote{
		table: table,
		meta:  mt,
		ref:   reflect.Indirect(reflect.ValueOf(data)),
		data:  data,
	}

	// 设置主键值
	mm.pk = dtype.ParseStr(mm.GetField(mm.table.PkField()).Interface())

	// 设置外键值
	if mm.table.FkField() != "" {
		mm.fk = dtype.ParseStr(mm.GetField(mm.table.FkField()).Interface())
	}

	return mm
}
