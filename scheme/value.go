package scheme

import (
	"micro-libs/utils/dtype"
	"micro-libs/utils/tool"
	"reflect"
)

type RefValue struct {
	isNil bool
	ref   reflect.Value
}

// 是否空字段
func (v *RefValue) IsNil() bool {
	return v.isNil
}

// 标记当前数据为空
func (v *RefValue) SetNil() {
	v.isNil = true
}

// 获取字段反射
func (v *RefValue) GetField(field string) reflect.Value {
	return v.ref.FieldByName(tool.UnderscoreToCamelCase(field))
}

// 获取反射值
func (v *RefValue) Get(field string) interface{} {
	vf := v.GetField(field)
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

func (v *RefValue) String(field string) string {
	return dtype.ParseStr(v.Get(field))
}

func (v *RefValue) Bool(field string) bool {
	return dtype.ParseBool(v.Get(field))
}

func (v *RefValue) Int32(field string) int32 {
	return dtype.ParseInt32(v.Get(field))
}

func (v *RefValue) Int64(field string) int64 {
	return dtype.ParseInt64(v.Get(field))
}

func (v *RefValue) Float32(field string) float32 {
	return dtype.ParseFloat32(v.Get(field))
}

func (v *RefValue) Float64(field string) float64 {
	return dtype.ParseFloat64(v.Get(field))
}

// 获取[]interface{}数据
func (v *RefValue) Slice(field string) []interface{} {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	var res = make([]interface{}, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		res = append(res, vf.Index(i).Interface())
	}

	return res
}

// 获取[][]interface{}数据
func (v *RefValue) SliceTwo(field string) [][]interface{} {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	var res = make([][]interface{}, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		row := vf.Index(i)
		if !row.IsValid() || row.Type().Kind() != reflect.Slice {
			continue
		}

		var two = make([]interface{}, 0, row.Len())
		for j := 0; j < row.Len(); j++ {
			two = append(two, row.Index(j).Interface())
		}

		res = append(res, two)
	}

	return res
}

// 获取[]string数据
func (v *RefValue) StrSlice(field string) []string {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([]string); ok {
		return v
	}

	var res = make([]string, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		res = append(res, dtype.ParseStr(vf.Index(i).Interface()))
	}

	return res
}

// 获取[][]string数据
func (v *RefValue) StrSliceTwo(field string) [][]string {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([][]string); ok {
		return v
	}

	var res = make([][]string, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		row := vf.Index(i)
		if !row.IsValid() || row.Type().Kind() != reflect.Slice {
			continue
		}

		var two = make([]string, 0, row.Len())
		for j := 0; j < row.Len(); j++ {
			two = append(two, dtype.ParseStr(row.Index(j).Interface()))
		}

		res = append(res, two)
	}

	return res
}

// 获取[]bool数据
func (v *RefValue) BoolSlice(field string) []bool {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([]bool); ok {
		return v
	}

	var res = make([]bool, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		res = append(res, dtype.ParseBool(vf.Index(i).Interface()))
	}

	return res
}

// 获取[][]bool数据
func (v *RefValue) BoolSliceTwo(field string) [][]bool {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([][]bool); ok {
		return v
	}

	var res = make([][]bool, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		row := vf.Index(i)
		if !row.IsValid() || row.Type().Kind() != reflect.Slice {
			continue
		}

		var two = make([]bool, 0, row.Len())
		for j := 0; j < row.Len(); j++ {
			two = append(two, dtype.ParseBool(row.Index(j).Interface()))
		}

		res = append(res, two)
	}

	return res
}

// 获取[]int32数据
func (v *RefValue) Int32Slice(field string) []int32 {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([]int32); ok {
		return v
	}

	var res = make([]int32, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		res = append(res, dtype.ParseInt32(vf.Index(i).Interface()))
	}

	return res
}

// 获取[][]int32数据
func (v *RefValue) Int32SliceTwo(field string) [][]int32 {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([][]int32); ok {
		return v
	}

	var res = make([][]int32, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		row := vf.Index(i)
		if !row.IsValid() || row.Type().Kind() != reflect.Slice {
			continue
		}

		var two = make([]int32, 0, row.Len())
		for j := 0; j < row.Len(); j++ {
			two = append(two, dtype.ParseInt32(row.Index(j).Interface()))
		}

		res = append(res, two)
	}

	return res
}

// 获取[]int64数据
func (v *RefValue) Int64Slice(field string) []int64 {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([]int64); ok {
		return v
	}

	var res = make([]int64, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		res = append(res, dtype.ParseInt64(vf.Index(i).Interface()))
	}

	return res
}

// 获取[][]int64数据
func (v *RefValue) Int64SliceTwo(field string) [][]int64 {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([][]int64); ok {
		return v
	}

	var res = make([][]int64, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		row := vf.Index(i)
		if !row.IsValid() || row.Type().Kind() != reflect.Slice {
			continue
		}

		var two = make([]int64, 0, row.Len())
		for j := 0; j < row.Len(); j++ {
			two = append(two, dtype.ParseInt64(row.Index(j).Interface()))
		}

		res = append(res, two)
	}

	return res
}

// 获取[]float32数据
func (v *RefValue) Float32Slice(field string) []float32 {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([]float32); ok {
		return v
	}

	var res = make([]float32, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		res = append(res, dtype.ParseFloat32(vf.Index(i).Interface()))
	}

	return res
}

// 获取[][]float32数据
func (v *RefValue) Float32SliceTwo(field string) [][]float32 {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([][]float32); ok {
		return v
	}

	var res = make([][]float32, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		row := vf.Index(i)
		if !row.IsValid() || row.Type().Kind() != reflect.Slice {
			continue
		}

		var two = make([]float32, 0, row.Len())
		for j := 0; j < row.Len(); j++ {
			two = append(two, dtype.ParseFloat32(row.Index(j).Interface()))
		}

		res = append(res, two)
	}

	return res
}

// 获取[]float64数据
func (v *RefValue) Float64Slice(field string) []float64 {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([]float64); ok {
		return v
	}

	var res = make([]float64, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		res = append(res, dtype.ParseFloat64(vf.Index(i).Interface()))
	}

	return res
}

// 获取[][]float64数据
func (v *RefValue) Float64SliceTwo(field string) [][]float64 {
	vf := v.GetField(field)
	if !vf.IsValid() || vf.Type().Kind() != reflect.Slice {
		return nil
	}

	if v, ok := vf.Interface().([][]float64); ok {
		return v
	}

	var res = make([][]float64, 0, vf.Len())
	for i := 0; i < vf.Len(); i++ {
		row := vf.Index(i)
		if !row.IsValid() || row.Type().Kind() != reflect.Slice {
			continue
		}

		var two = make([]float64, 0, row.Len())
		for j := 0; j < row.Len(); j++ {
			two = append(two, dtype.ParseFloat64(row.Index(j).Interface()))
		}

		res = append(res, two)
	}

	return res
}

func NewRefValue(val interface{}) *RefValue {
	return &RefValue{
		ref: reflect.Indirect(reflect.ValueOf(val)),
	}
}
