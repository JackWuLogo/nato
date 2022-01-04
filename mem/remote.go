package mem

import (
	"context"
	"encoding/json"
	"micro-libs/meta"
	"micro-libs/utils/dtype"
	"micro-libs/utils/errors"
	"reflect"
)

type Remote struct {
	admin *Admin
}

// 获取主键数据.
// 如果pk为空, 表示获取无外键的主键数据;
// 如果pk不为空, 表示获取有外键的主键数据
func (s *Remote) Get(ctx context.Context, req *InMemGet, rsp *OutMemGet) error {
	mt, err := meta.FromDataMeta(ctx)
	if err != nil {
		return err
	}

	tab, err := s.admin.Table(req.Table)
	if err != nil {
		return err
	}

	mm, err := tab.Get(mt, req.Pk...)
	if err != nil {
		return err
	}

	mm.RLock()
	rsp.Result = mm.Byte()
	mm.RUnlock()

	return nil
}

// 获取所有外键数据
func (s *Remote) GetFk(ctx context.Context, req *InMemGetFk, rsp *OutMemGetFk) error {
	mt, err := meta.FromDataMeta(ctx)
	if err != nil {
		return err
	}

	tab, err := s.admin.Table(req.Table)
	if err != nil {
		return err
	}

	rows, err := tab.GetFk(mt, req.Pk...)
	if err != nil {
		return err
	}

	rsp.Result = make(map[string][]byte)
	for key, row := range rows {
		row.RLock()
		rsp.Result[key] = row.Byte()
		row.RUnlock()
	}

	return nil
}

// 写入新数据
func (s *Remote) Insert(ctx context.Context, req *InMemInsert, _ *MemNone) error {
	mt, err := meta.FromDataMeta(ctx)
	if err != nil {
		return err
	}

	tab, err := s.admin.Table(req.Table)
	if err != nil {
		return err
	}

	_, err2 := tab.Insert(mt, req.Data, req.Pk...)
	return err2
}

// 从数据库删除数据模型 (数据库立即删除)
func (s *Remote) Delete(ctx context.Context, req *InMemDelete, _ *MemNone) error {
	mt, err := meta.FromDataMeta(ctx)
	if err != nil {
		return err
	}

	tab, err := s.admin.Table(req.Table)
	if err != nil {
		return err
	}

	return tab.Delete(mt, req.Pk...)
}

// 更新单个数据
func (s *Remote) SetValue(ctx context.Context, req *InMemSetValue, _ *MemNone) error {
	mt, err := meta.FromDataMeta(ctx)
	if err != nil {
		return err
	}

	tab, err := s.admin.Table(req.Table)
	if err != nil {
		return err
	}

	mm, err := tab.Get(mt, req.Pk...)
	if err != nil {
		return err
	}

	local, ok := mm.(*ModelLocal)
	if !ok {
		return errors.Unavailable("[Remote] update value meta node error")
	}

	// 更新数据
	local.Lock()
	local.SetRemoteValue(req.Field, req.Value)
	local.Unlock()

	return nil
}

func (s *Remote) SetValues(ctx context.Context, req *InMemSetValues, _ *MemNone) error {
	mt, err := meta.FromDataMeta(ctx)
	if err != nil {
		return err
	}

	tab, err := s.admin.Table(req.Table)
	if err != nil {
		return err
	}

	mm, err := tab.Get(mt, req.Pk...)
	if err != nil {
		return err
	}

	local, ok := mm.(*ModelLocal)
	if !ok {
		return errors.Unavailable("[Remote] update value meta node error")
	}

	// 更新数据
	local.Lock()
	local.SetRemoteValues(req.Values)
	local.Unlock()

	return nil
}

func NewRemote(admin *Admin) *Remote {
	return &Remote{admin: admin}
}

// 解析数据
func ParseFieldValue(mm MModel, field string, value []byte) (interface{}, error) {
	vf := mm.GetField(field)
	if !vf.IsValid() || !vf.CanSet() {
		return nil, errors.Unavailable("invalid field %s", field)
	}

	var isPtr bool
	var res reflect.Value

	var elem = vf
	var typ = vf.Type()
	if typ.Kind() == reflect.Ptr {
		elem = elem.Elem()
		isPtr = true
	}

	switch typ.Kind() {
	case reflect.String:
		res = dtype.StrElem()
	case reflect.Int:
		res = dtype.IntElem()
	case reflect.Int32:
		res = dtype.Int32Elem()
	case reflect.Int64:
		res = dtype.Int64Elem()
	case reflect.Float32:
		res = dtype.Float32Elem()
	case reflect.Float64:
		res = dtype.Float64Elem()
	case reflect.Bool:
		res = dtype.BoolElem()
	case reflect.Struct:
		res = dtype.Elem(typ)
	case reflect.Slice:
		res = dtype.SliceElem(typ.Elem())
	case reflect.Map:
		res = dtype.MapElem(typ.Key(), typ.Elem())
	}

	if err := json.Unmarshal(value, res.Addr().Interface()); err != nil {
		return nil, err
	}

	if isPtr {
		return res.Addr().Interface(), nil
	}
	return res.Interface(), nil
}
