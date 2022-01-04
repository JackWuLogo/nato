// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.11.4
// source: mem/mem.proto

package mem

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// 无返回内容
type MemNone struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *MemNone) Reset() {
	*x = MemNone{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MemNone) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MemNone) ProtoMessage() {}

func (x *MemNone) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MemNone.ProtoReflect.Descriptor instead.
func (*MemNone) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{0}
}

// 获取指定主键数据
type InMemGet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Table string   `protobuf:"bytes,1,opt,name=Table,proto3" json:"Table,omitempty"`
	Pk    []string `protobuf:"bytes,2,rep,name=Pk,proto3" json:"Pk,omitempty"`
}

func (x *InMemGet) Reset() {
	*x = InMemGet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InMemGet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InMemGet) ProtoMessage() {}

func (x *InMemGet) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InMemGet.ProtoReflect.Descriptor instead.
func (*InMemGet) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{1}
}

func (x *InMemGet) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *InMemGet) GetPk() []string {
	if x != nil {
		return x.Pk
	}
	return nil
}

type OutMemGet struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result []byte `protobuf:"bytes,1,opt,name=Result,proto3" json:"Result,omitempty"`
}

func (x *OutMemGet) Reset() {
	*x = OutMemGet{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutMemGet) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutMemGet) ProtoMessage() {}

func (x *OutMemGet) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutMemGet.ProtoReflect.Descriptor instead.
func (*OutMemGet) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{2}
}

func (x *OutMemGet) GetResult() []byte {
	if x != nil {
		return x.Result
	}
	return nil
}

// 获取外键数据
type InMemGetFk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Table string   `protobuf:"bytes,1,opt,name=Table,proto3" json:"Table,omitempty"`
	Pk    []string `protobuf:"bytes,2,rep,name=Pk,proto3" json:"Pk,omitempty"`
}

func (x *InMemGetFk) Reset() {
	*x = InMemGetFk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InMemGetFk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InMemGetFk) ProtoMessage() {}

func (x *InMemGetFk) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InMemGetFk.ProtoReflect.Descriptor instead.
func (*InMemGetFk) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{3}
}

func (x *InMemGetFk) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *InMemGetFk) GetPk() []string {
	if x != nil {
		return x.Pk
	}
	return nil
}

type OutMemGetFk struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result map[string][]byte `protobuf:"bytes,1,rep,name=Result,proto3" json:"Result,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *OutMemGetFk) Reset() {
	*x = OutMemGetFk{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OutMemGetFk) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OutMemGetFk) ProtoMessage() {}

func (x *OutMemGetFk) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OutMemGetFk.ProtoReflect.Descriptor instead.
func (*OutMemGetFk) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{4}
}

func (x *OutMemGetFk) GetResult() map[string][]byte {
	if x != nil {
		return x.Result
	}
	return nil
}

// 新增数据
type InMemInsert struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Table string   `protobuf:"bytes,1,opt,name=Table,proto3" json:"Table,omitempty"`
	Data  []byte   `protobuf:"bytes,2,opt,name=Data,proto3" json:"Data,omitempty"`
	Pk    []string `protobuf:"bytes,3,rep,name=Pk,proto3" json:"Pk,omitempty"`
}

func (x *InMemInsert) Reset() {
	*x = InMemInsert{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InMemInsert) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InMemInsert) ProtoMessage() {}

func (x *InMemInsert) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InMemInsert.ProtoReflect.Descriptor instead.
func (*InMemInsert) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{5}
}

func (x *InMemInsert) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *InMemInsert) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

func (x *InMemInsert) GetPk() []string {
	if x != nil {
		return x.Pk
	}
	return nil
}

// 删除数据
type InMemDelete struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Table string   `protobuf:"bytes,1,opt,name=Table,proto3" json:"Table,omitempty"`
	Pk    []string `protobuf:"bytes,2,rep,name=Pk,proto3" json:"Pk,omitempty"`
}

func (x *InMemDelete) Reset() {
	*x = InMemDelete{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InMemDelete) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InMemDelete) ProtoMessage() {}

func (x *InMemDelete) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InMemDelete.ProtoReflect.Descriptor instead.
func (*InMemDelete) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{6}
}

func (x *InMemDelete) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *InMemDelete) GetPk() []string {
	if x != nil {
		return x.Pk
	}
	return nil
}

// 设置字段值
type InMemSetValue struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Table string   `protobuf:"bytes,1,opt,name=Table,proto3" json:"Table,omitempty"`
	Pk    []string `protobuf:"bytes,2,rep,name=Pk,proto3" json:"Pk,omitempty"`
	Field string   `protobuf:"bytes,3,opt,name=Field,proto3" json:"Field,omitempty"`
	Value []byte   `protobuf:"bytes,4,opt,name=Value,proto3" json:"Value,omitempty"`
}

func (x *InMemSetValue) Reset() {
	*x = InMemSetValue{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InMemSetValue) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InMemSetValue) ProtoMessage() {}

func (x *InMemSetValue) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InMemSetValue.ProtoReflect.Descriptor instead.
func (*InMemSetValue) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{7}
}

func (x *InMemSetValue) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *InMemSetValue) GetPk() []string {
	if x != nil {
		return x.Pk
	}
	return nil
}

func (x *InMemSetValue) GetField() string {
	if x != nil {
		return x.Field
	}
	return ""
}

func (x *InMemSetValue) GetValue() []byte {
	if x != nil {
		return x.Value
	}
	return nil
}

// 批量设置字段值
type InMemSetValues struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Table  string            `protobuf:"bytes,1,opt,name=Table,proto3" json:"Table,omitempty"`
	Pk     []string          `protobuf:"bytes,2,rep,name=Pk,proto3" json:"Pk,omitempty"`
	Values map[string][]byte `protobuf:"bytes,3,rep,name=Values,proto3" json:"Values,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *InMemSetValues) Reset() {
	*x = InMemSetValues{}
	if protoimpl.UnsafeEnabled {
		mi := &file_mem_mem_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *InMemSetValues) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*InMemSetValues) ProtoMessage() {}

func (x *InMemSetValues) ProtoReflect() protoreflect.Message {
	mi := &file_mem_mem_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use InMemSetValues.ProtoReflect.Descriptor instead.
func (*InMemSetValues) Descriptor() ([]byte, []int) {
	return file_mem_mem_proto_rawDescGZIP(), []int{8}
}

func (x *InMemSetValues) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *InMemSetValues) GetPk() []string {
	if x != nil {
		return x.Pk
	}
	return nil
}

func (x *InMemSetValues) GetValues() map[string][]byte {
	if x != nil {
		return x.Values
	}
	return nil
}

var File_mem_mem_proto protoreflect.FileDescriptor

var file_mem_mem_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x6d, 0x65, 0x6d, 0x2f, 0x6d, 0x65, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x03, 0x6d, 0x65, 0x6d, 0x22, 0x09, 0x0a, 0x07, 0x4d, 0x65, 0x6d, 0x4e, 0x6f, 0x6e, 0x65, 0x22,
	0x30, 0x0a, 0x08, 0x49, 0x6e, 0x4d, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x54,
	0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x61, 0x62, 0x6c,
	0x65, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x6b, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x50,
	0x6b, 0x22, 0x23, 0x0a, 0x09, 0x4f, 0x75, 0x74, 0x4d, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x12, 0x16,
	0x0a, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06,
	0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x22, 0x32, 0x0a, 0x0a, 0x49, 0x6e, 0x4d, 0x65, 0x6d, 0x47,
	0x65, 0x74, 0x46, 0x6b, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x6b,
	0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x50, 0x6b, 0x22, 0x7e, 0x0a, 0x0b, 0x4f, 0x75,
	0x74, 0x4d, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x46, 0x6b, 0x12, 0x34, 0x0a, 0x06, 0x52, 0x65, 0x73,
	0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x1c, 0x2e, 0x6d, 0x65, 0x6d, 0x2e,
	0x4f, 0x75, 0x74, 0x4d, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x46, 0x6b, 0x2e, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x1a,
	0x39, 0x0a, 0x0b, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10,
	0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79,
	0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x47, 0x0a, 0x0b, 0x49, 0x6e,
	0x4d, 0x65, 0x6d, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x12, 0x0a, 0x04, 0x44, 0x61, 0x74, 0x61, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x44,
	0x61, 0x74, 0x61, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x6b, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52,
	0x02, 0x50, 0x6b, 0x22, 0x33, 0x0a, 0x0b, 0x49, 0x6e, 0x4d, 0x65, 0x6d, 0x44, 0x65, 0x6c, 0x65,
	0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x6b, 0x18, 0x02,
	0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x50, 0x6b, 0x22, 0x61, 0x0a, 0x0d, 0x49, 0x6e, 0x4d, 0x65,
	0x6d, 0x53, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x54, 0x61, 0x62,
	0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x50, 0x6b, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x02, 0x50, 0x6b, 0x12,
	0x14, 0x0a, 0x05, 0x46, 0x69, 0x65, 0x6c, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05,
	0x46, 0x69, 0x65, 0x6c, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x22, 0xaa, 0x01, 0x0a, 0x0e,
	0x49, 0x6e, 0x4d, 0x65, 0x6d, 0x53, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x12, 0x14,
	0x0a, 0x05, 0x54, 0x61, 0x62, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x54,
	0x61, 0x62, 0x6c, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x50, 0x6b, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x02, 0x50, 0x6b, 0x12, 0x37, 0x0a, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x18, 0x03,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x1f, 0x2e, 0x6d, 0x65, 0x6d, 0x2e, 0x49, 0x6e, 0x4d, 0x65, 0x6d,
	0x53, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73,
	0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x06, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x1a, 0x39, 0x0a,
	0x0b, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x32, 0x8c, 0x02, 0x0a, 0x06, 0x52, 0x65, 0x6d,
	0x6f, 0x74, 0x65, 0x12, 0x24, 0x0a, 0x03, 0x47, 0x65, 0x74, 0x12, 0x0d, 0x2e, 0x6d, 0x65, 0x6d,
	0x2e, 0x49, 0x6e, 0x4d, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x1a, 0x0e, 0x2e, 0x6d, 0x65, 0x6d, 0x2e,
	0x4f, 0x75, 0x74, 0x4d, 0x65, 0x6d, 0x47, 0x65, 0x74, 0x12, 0x2a, 0x0a, 0x05, 0x47, 0x65, 0x74,
	0x46, 0x6b, 0x12, 0x0f, 0x2e, 0x6d, 0x65, 0x6d, 0x2e, 0x49, 0x6e, 0x4d, 0x65, 0x6d, 0x47, 0x65,
	0x74, 0x46, 0x6b, 0x1a, 0x10, 0x2e, 0x6d, 0x65, 0x6d, 0x2e, 0x4f, 0x75, 0x74, 0x4d, 0x65, 0x6d,
	0x47, 0x65, 0x74, 0x46, 0x6b, 0x12, 0x28, 0x0a, 0x06, 0x49, 0x6e, 0x73, 0x65, 0x72, 0x74, 0x12,
	0x10, 0x2e, 0x6d, 0x65, 0x6d, 0x2e, 0x49, 0x6e, 0x4d, 0x65, 0x6d, 0x49, 0x6e, 0x73, 0x65, 0x72,
	0x74, 0x1a, 0x0c, 0x2e, 0x6d, 0x65, 0x6d, 0x2e, 0x4d, 0x65, 0x6d, 0x4e, 0x6f, 0x6e, 0x65, 0x12,
	0x28, 0x0a, 0x06, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x12, 0x10, 0x2e, 0x6d, 0x65, 0x6d, 0x2e,
	0x49, 0x6e, 0x4d, 0x65, 0x6d, 0x44, 0x65, 0x6c, 0x65, 0x74, 0x65, 0x1a, 0x0c, 0x2e, 0x6d, 0x65,
	0x6d, 0x2e, 0x4d, 0x65, 0x6d, 0x4e, 0x6f, 0x6e, 0x65, 0x12, 0x2c, 0x0a, 0x08, 0x53, 0x65, 0x74,
	0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x12, 0x2e, 0x6d, 0x65, 0x6d, 0x2e, 0x49, 0x6e, 0x4d, 0x65,
	0x6d, 0x53, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x0c, 0x2e, 0x6d, 0x65, 0x6d, 0x2e,
	0x4d, 0x65, 0x6d, 0x4e, 0x6f, 0x6e, 0x65, 0x12, 0x2e, 0x0a, 0x09, 0x53, 0x65, 0x74, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x73, 0x12, 0x13, 0x2e, 0x6d, 0x65, 0x6d, 0x2e, 0x49, 0x6e, 0x4d, 0x65, 0x6d,
	0x53, 0x65, 0x74, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x73, 0x1a, 0x0c, 0x2e, 0x6d, 0x65, 0x6d, 0x2e,
	0x4d, 0x65, 0x6d, 0x4e, 0x6f, 0x6e, 0x65, 0x42, 0x10, 0x5a, 0x0e, 0x6d, 0x69, 0x63, 0x72, 0x6f,
	0x2d, 0x6c, 0x69, 0x62, 0x73, 0x2f, 0x6d, 0x65, 0x6d, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x33,
}

var (
	file_mem_mem_proto_rawDescOnce sync.Once
	file_mem_mem_proto_rawDescData = file_mem_mem_proto_rawDesc
)

func file_mem_mem_proto_rawDescGZIP() []byte {
	file_mem_mem_proto_rawDescOnce.Do(func() {
		file_mem_mem_proto_rawDescData = protoimpl.X.CompressGZIP(file_mem_mem_proto_rawDescData)
	})
	return file_mem_mem_proto_rawDescData
}

var file_mem_mem_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_mem_mem_proto_goTypes = []interface{}{
	(*MemNone)(nil),        // 0: mem.MemNone
	(*InMemGet)(nil),       // 1: mem.InMemGet
	(*OutMemGet)(nil),      // 2: mem.OutMemGet
	(*InMemGetFk)(nil),     // 3: mem.InMemGetFk
	(*OutMemGetFk)(nil),    // 4: mem.OutMemGetFk
	(*InMemInsert)(nil),    // 5: mem.InMemInsert
	(*InMemDelete)(nil),    // 6: mem.InMemDelete
	(*InMemSetValue)(nil),  // 7: mem.InMemSetValue
	(*InMemSetValues)(nil), // 8: mem.InMemSetValues
	nil,                    // 9: mem.OutMemGetFk.ResultEntry
	nil,                    // 10: mem.InMemSetValues.ValuesEntry
}
var file_mem_mem_proto_depIdxs = []int32{
	9,  // 0: mem.OutMemGetFk.Result:type_name -> mem.OutMemGetFk.ResultEntry
	10, // 1: mem.InMemSetValues.Values:type_name -> mem.InMemSetValues.ValuesEntry
	1,  // 2: mem.Remote.Get:input_type -> mem.InMemGet
	3,  // 3: mem.Remote.GetFk:input_type -> mem.InMemGetFk
	5,  // 4: mem.Remote.Insert:input_type -> mem.InMemInsert
	6,  // 5: mem.Remote.Delete:input_type -> mem.InMemDelete
	7,  // 6: mem.Remote.SetValue:input_type -> mem.InMemSetValue
	8,  // 7: mem.Remote.SetValues:input_type -> mem.InMemSetValues
	2,  // 8: mem.Remote.Get:output_type -> mem.OutMemGet
	4,  // 9: mem.Remote.GetFk:output_type -> mem.OutMemGetFk
	0,  // 10: mem.Remote.Insert:output_type -> mem.MemNone
	0,  // 11: mem.Remote.Delete:output_type -> mem.MemNone
	0,  // 12: mem.Remote.SetValue:output_type -> mem.MemNone
	0,  // 13: mem.Remote.SetValues:output_type -> mem.MemNone
	8,  // [8:14] is the sub-list for method output_type
	2,  // [2:8] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_mem_mem_proto_init() }
func file_mem_mem_proto_init() {
	if File_mem_mem_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_mem_mem_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MemNone); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mem_mem_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InMemGet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mem_mem_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OutMemGet); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mem_mem_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InMemGetFk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mem_mem_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OutMemGetFk); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mem_mem_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InMemInsert); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mem_mem_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InMemDelete); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mem_mem_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InMemSetValue); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_mem_mem_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*InMemSetValues); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_mem_mem_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_mem_mem_proto_goTypes,
		DependencyIndexes: file_mem_mem_proto_depIdxs,
		MessageInfos:      file_mem_mem_proto_msgTypes,
	}.Build()
	File_mem_mem_proto = out.File
	file_mem_mem_proto_rawDesc = nil
	file_mem_mem_proto_goTypes = nil
	file_mem_mem_proto_depIdxs = nil
}