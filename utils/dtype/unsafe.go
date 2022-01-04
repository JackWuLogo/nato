package dtype

import (
	"reflect"
	"unsafe"
)

// UnsafeStructToByte 仅支持转换结构体, 结构体字段只能是基础数据类型
func UnsafeStructToByte(data interface{}) []byte {
	return *(*[]byte)(unsafe.Pointer(reflect.ValueOf(data).Pointer()))
}

// UnsafeByteToStruct 仅支持转换结构体, 结构体字段只能是基础数据类型
func UnsafeByteToStruct(typ reflect.Type, data []byte) reflect.Value {
	return reflect.NewAt(typ, unsafe.Pointer(&data))
}

// 深度拷贝结构体, 结构体字段只能是基础数据类型
func UnsafeCopyStruct(old interface{}) interface{} {
	rf := reflect.TypeOf(old)
	if rf.Kind() == reflect.Ptr {
		rf = rf.Elem()
	}
	bytes := *(*[]byte)(unsafe.Pointer(reflect.ValueOf(old).Pointer()))
	return reflect.NewAt(rf, unsafe.Pointer(&bytes)).Interface()
}

func UnsafeByteToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func UnsafeStringToByte(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
