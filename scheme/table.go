package scheme

import (
	"fmt"
	"micro-libs/utils/dtype"
	"micro-libs/utils/errors"
	"reflect"
	"strconv"
	"sync"
)

var (
	nilTable = &Table{
		values: make(map[string]interface{}),
	}
)

// Field 配置表字段信息
type Field struct {
	Name   string `json:"name"`   // 字段名
	Key    string `json:"key"`    // 字段标识
	Title  string `json:"title"`  // 显示名称
	Type   string `json:"type"`   // 数据类型
	Export bool   `json:"export"` // 是否导出
}

// Table 配置表信息
type Table struct {
	sync.RWMutex
	scheme  *Scheme
	Key     string                 // 标识信息, 即Excel的 sheet 名称
	Name    string                 // 数据表名, 即Excel的文件名
	Ref     reflect.Type           // 数据模型反射
	values  map[string]interface{} // 数据信息
	Version int64                  // 版本号
}

func (s *Table) IsNil() bool {
	return s.scheme == nil
}

// New 使用反射创建新的对象 (指针对象)
func (s *Table) New() reflect.Value {
	return dtype.Ptr(s.Ref)
}

// 使用反射创建新的对象 (指针数组)
func (s *Table) NewSlice() reflect.Value {
	return dtype.SliceElem(s.Ref)
}

// 使用反射创建新的对象 (指针数组)
func (s *Table) NewMap() reflect.Value {
	return dtype.StrMapElem(s.Ref)
}

// 设置当前数据
func (s *Table) SetValues(ver int64, rows reflect.Value) error {
	s.Lock()
	defer s.Unlock()

	rows = reflect.Indirect(rows)
	if rows.Type().Kind() != reflect.Slice {
		return errors.ParamInvalid("Unsupported data types")
	}

	s.Version = ver

	s.values = make(map[string]interface{})
	for i := 0; i < rows.Len(); i++ {
		row := rows.Index(i)
		id := row.Elem().FieldByName("Id")
		var key string
		switch id.Kind() {
		case reflect.Int32, reflect.Int64:
			key = strconv.FormatInt(id.Int(), 10)
		case reflect.String:
			key = id.String()
		default:
			continue
		}

		ref := row.Elem().FieldByName("RefValue")
		if ref.CanSet() {
			ref.Set(reflect.ValueOf(NewRefValue(row.Interface())))
		}

		s.values[key] = row.Interface()
	}

	return nil
}

func (s *Table) Fields() []*Field {
	var fields []*Field
	for i := 0; i < s.Ref.NumField(); i++ {
		ef := s.Ref.Field(i)
		key := ef.Tag.Get("bson")
		if key == "" || key == "-" {
			continue
		}
		fields = append(fields, &Field{
			Name:   ef.Name,
			Key:    key,
			Title:  ef.Tag.Get("title"),
			Type:   ef.Type.String(),
			Export: dtype.ParseBool(ef.Tag.Get("export")),
		})
	}
	return fields
}

// 获取数据总数
func (s *Table) Total() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.values)
}

// 获取全部数据
func (s *Table) Values() map[string]interface{} {
	s.RLock()
	defer s.RUnlock()

	var values = make(map[string]interface{}, len(s.values))
	for k, v := range s.values {
		values[k] = v
	}

	return values
}

// 迭代数据
func (s *Table) Range(fn func(k string, v interface{}) bool) {
	s.RLock()
	defer s.RUnlock()

	for k, v := range s.values {
		if b := fn(k, v); !b {
			break
		}
	}
}

// 获取数据
func (s *Table) Get(id interface{}) interface{} {
	if id == nil {
		return nil
	}

	var key string
	switch v := id.(type) {
	case string:
		key = v
	case int:
		key = strconv.Itoa(v)
	case int32:
		key = strconv.Itoa(int(v))
	case int64:
		key = strconv.FormatInt(v, 10)
	default:
		key = fmt.Sprint(v)
	}

	s.RLock()
	if val, ok := s.values[key]; ok {
		s.RUnlock()
		return val
	}
	s.RUnlock()

	return nil
}

// 数据大小
func (s *Table) Len() int {
	s.RLock()
	defer s.RUnlock()

	return len(s.values)
}

// 实例化数据表
func NewTable(s *Scheme, key, name string, model interface{}) *Table {
	t := &Table{
		scheme: s,
		Key:    key,
		Name:   name,
		Ref:    reflect.TypeOf(model),
		values: make(map[string]interface{}),
	}
	return t
}
