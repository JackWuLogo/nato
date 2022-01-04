package scheme

import (
	"encoding/json"
	"fmt"
	"micro-libs/utils/dtype"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// 格式化数据类型
func FormatDataType(t string) string {
	switch t {
	case "bool":
		return "bool"
	case "string", "str":
		return "string"
	case "int32", "int":
		return "int32"
	case "int64":
		return "int64"
	case "float32":
		return "float32"
	case "float64", "float":
		return "float64"
	case "any":
		return "interface{}"
	case "date", "datetime":
		return t
	case "[]bool", "bool[]", "bool1":
		return "[]bool"
	case "[][]bool", "bool[][]", "bool2":
		return "[][]bool"
	case "[]string", "string[]", "str[]", "string1", "str1":
		return "[]string"
	case "[][]string", "string[][]", "str[][]", "string2", "str2":
		return "[][]string"
	case "[]int32", "[]int", "int32[]", "int[]", "int1":
		return "[]int32"
	case "[][]int32", "[][]int", "int32[][]", "int[][]", "int2":
		return "[][]int32"
	case "[]int64", "int64[]":
		return "[]int64"
	case "[][]int64", "int64[][]":
		return "[][]int64"
	case "[]float32", "float32[]":
		return "[]float32"
	case "[][]float32", "float32[][]":
		return "[][]float32"
	case "[]float64", "float64[]", "[]float", "float[]", "float1":
		return "[]float64"
	case "[][]float64", "float64[][]", "[][]float", "float[][]", "float2":
		return "[][]float64"
	case "[]interface{}", "any[]", "any1", "array1":
		return "[]interface{}"
	case "[][]interface{}", "any[][]", "any2", "array2":
		return "[][]interface{}"
	default:
		return "string"
	}
}

// 格式化数据
func FormatValue(typ string, val string) interface{} {
	validType := FormatDataType(typ)
	switch validType {
	case "string":
		return val
	case "bool":
		v, _ := strconv.ParseBool(val)
		return v
	case "int32":
		v, _ := strconv.Atoi(val)
		return int32(v)
	case "int64":
		v, _ := strconv.ParseInt(val, 10, 64)
		return v
	case "float32":
		v, _ := strconv.ParseFloat(val, 32)
		return float32(v)
	case "float64":
		v, _ := strconv.ParseFloat(val, 64)
		return v
	case "date":
		if val != "" {
			t, err := time.ParseInLocation("2006-01-02", val, time.Local)
			if err == nil {
				return t.Unix()
			}
		}
		return 0
	case "datetime":
		if val != "" {
			t, err := time.ParseInLocation("2006-01-02 15:04:05", val, time.Local)
			if err == nil {
				return t.Unix()
			}
		}
		return 0
	case "[]string", "[]bool", "[]int32", "[]int64", "[]float32", "[]float64":
		if !strings.HasPrefix(val, "[") && strings.Contains(val, ",") {
			if typ == "[]string" {
				return strings.Split(val, ",")
			}
			val = fmt.Sprintf("[%s]", val)
		}
		fallthrough
	default:
		var elem reflect.Value
		switch validType {
		case "[]string":
			elem = dtype.StrSliceElem()
		case "[]bool":
			elem = dtype.BoolSliceElem()
		case "[]int32":
			elem = dtype.Int32SliceElem()
		case "[]int64":
			elem = dtype.Int64SliceElem()
		case "[]float32":
			elem = dtype.Float32SliceElem()
		case "[]float64":
			elem = dtype.Float64SliceElem()
		case "[][]string":
			elem = dtype.SliceTwoElem(dtype.RefTypeStr)
		case "[][]bool":
			elem = dtype.SliceTwoElem(dtype.RefTypeBool)
		case "[][]int32":
			elem = dtype.SliceTwoElem(dtype.RefTypeInt32)
		case "[][]int64":
			elem = dtype.SliceTwoElem(dtype.RefTypeInt64)
		case "[][]float32":
			elem = dtype.SliceTwoElem(dtype.RefTypeFloat32)
		case "[][]float64":
			elem = dtype.SliceTwoElem(dtype.RefTypeFloat64)
		default:
			return nil
		}

		_ = json.Unmarshal([]byte(val), elem.Addr().Interface())
		return elem.Interface()
	}
}
