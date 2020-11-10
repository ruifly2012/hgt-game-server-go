// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: protos/protobuf/GameMessage.proto

package protobuf

import (
	proto "github.com/golang/protobuf/proto"
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

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

// 基础信息
type Message struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// 协议号
	Protocol int64 `protobuf:"varint,1,opt,name=protocol,proto3" json:"protocol,omitempty"`
	// 返回码
	Code int64 `protobuf:"varint,2,opt,name=code,proto3" json:"code,omitempty"`
	// 数据
	Data []byte `protobuf:"bytes,3,opt,name=data,proto3" json:"data,omitempty"`
}

func (x *Message) Reset() {
	*x = Message{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Message) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Message) ProtoMessage() {}

func (x *Message) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Message.ProtoReflect.Descriptor instead.
func (*Message) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{0}
}

func (x *Message) GetProtocol() int64 {
	if x != nil {
		return x.Protocol
	}
	return 0
}

func (x *Message) GetCode() int64 {
	if x != nil {
		return x.Code
	}
	return 0
}

func (x *Message) GetData() []byte {
	if x != nil {
		return x.Data
	}
	return nil
}

// 2 心跳
type HeartBeatReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HeartBeatReq) Reset() {
	*x = HeartBeatReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartBeatReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartBeatReq) ProtoMessage() {}

func (x *HeartBeatReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartBeatReq.ProtoReflect.Descriptor instead.
func (*HeartBeatReq) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{1}
}

// -2 心跳 返回
type HeartBeatRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *HeartBeatRes) Reset() {
	*x = HeartBeatRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *HeartBeatRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*HeartBeatRes) ProtoMessage() {}

func (x *HeartBeatRes) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use HeartBeatRes.ProtoReflect.Descriptor instead.
func (*HeartBeatRes) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{2}
}

// 1003 注册
type RegisterReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *RegisterReq) Reset() {
	*x = RegisterReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterReq) ProtoMessage() {}

func (x *RegisterReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterReq.ProtoReflect.Descriptor instead.
func (*RegisterReq) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{3}
}

func (x *RegisterReq) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *RegisterReq) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// -1003 注册返回
type RegisterRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id       string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	Username string `protobuf:"bytes,2,opt,name=username,proto3" json:"username,omitempty"`
}

func (x *RegisterRes) Reset() {
	*x = RegisterRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RegisterRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RegisterRes) ProtoMessage() {}

func (x *RegisterRes) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RegisterRes.ProtoReflect.Descriptor instead.
func (*RegisterRes) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{4}
}

func (x *RegisterRes) GetId() string {
	if x != nil {
		return x.Id
	}
	return ""
}

func (x *RegisterRes) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

// 1002 登录
type LoginReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
}

func (x *LoginReq) Reset() {
	*x = LoginReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginReq) ProtoMessage() {}

func (x *LoginReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginReq.ProtoReflect.Descriptor instead.
func (*LoginReq) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{5}
}

func (x *LoginReq) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *LoginReq) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

// -1002 登录返回
type LoginRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	UserId   string `protobuf:"bytes,2,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *LoginRes) Reset() {
	*x = LoginRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LoginRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LoginRes) ProtoMessage() {}

func (x *LoginRes) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LoginRes.ProtoReflect.Descriptor instead.
func (*LoginRes) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{6}
}

func (x *LoginRes) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *LoginRes) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type OfflineReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId string `protobuf:"bytes,1,opt,name=userId,proto3" json:"userId,omitempty"`
}

func (x *OfflineReq) Reset() {
	*x = OfflineReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OfflineReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OfflineReq) ProtoMessage() {}

func (x *OfflineReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OfflineReq.ProtoReflect.Descriptor instead.
func (*OfflineReq) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{7}
}

func (x *OfflineReq) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

type OfflineRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *OfflineRes) Reset() {
	*x = OfflineRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *OfflineRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*OfflineRes) ProtoMessage() {}

func (x *OfflineRes) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use OfflineRes.ProtoReflect.Descriptor instead.
func (*OfflineRes) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{8}
}

type TestReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Str string `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *TestReq) Reset() {
	*x = TestReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestReq) ProtoMessage() {}

func (x *TestReq) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestReq.ProtoReflect.Descriptor instead.
func (*TestReq) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{9}
}

func (x *TestReq) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

type TestRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Str string `protobuf:"bytes,1,opt,name=str,proto3" json:"str,omitempty"`
}

func (x *TestRes) Reset() {
	*x = TestRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_protos_protobuf_GameMessage_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TestRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestRes) ProtoMessage() {}

func (x *TestRes) ProtoReflect() protoreflect.Message {
	mi := &file_protos_protobuf_GameMessage_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestRes.ProtoReflect.Descriptor instead.
func (*TestRes) Descriptor() ([]byte, []int) {
	return file_protos_protobuf_GameMessage_proto_rawDescGZIP(), []int{10}
}

func (x *TestRes) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

var File_protos_protobuf_GameMessage_proto protoreflect.FileDescriptor

var file_protos_protobuf_GameMessage_proto_rawDesc = []byte{
	0x0a, 0x21, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2f, 0x47, 0x61, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x0b, 0x47, 0x61, 0x6d, 0x65, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x22, 0x4d, 0x0a, 0x07, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64,
	0x61, 0x74, 0x61, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x04, 0x64, 0x61, 0x74, 0x61, 0x22,
	0x0e, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x52, 0x65, 0x71, 0x22,
	0x0e, 0x0a, 0x0c, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x52, 0x65, 0x73, 0x22,
	0x45, 0x0a, 0x0b, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74, 0x65, 0x72, 0x52, 0x65, 0x71, 0x12, 0x1a,
	0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61,
	0x73, 0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x39, 0x0a, 0x0b, 0x52, 0x65, 0x67, 0x69, 0x73, 0x74,
	0x65, 0x72, 0x52, 0x65, 0x73, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x02, 0x69, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d,
	0x65, 0x22, 0x42, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73,
	0x73, 0x77, 0x6f, 0x72, 0x64, 0x22, 0x3e, 0x0a, 0x08, 0x4c, 0x6f, 0x67, 0x69, 0x6e, 0x52, 0x65,
	0x73, 0x12, 0x1a, 0x0a, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x08, 0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x24, 0x0a, 0x0a, 0x4f, 0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65,
	0x52, 0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x22, 0x0c, 0x0a, 0x0a, 0x4f,
	0x66, 0x66, 0x6c, 0x69, 0x6e, 0x65, 0x52, 0x65, 0x73, 0x22, 0x1b, 0x0a, 0x07, 0x54, 0x65, 0x73,
	0x74, 0x52, 0x65, 0x71, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x22, 0x1b, 0x0a, 0x07, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65,
	0x73, 0x12, 0x10, 0x0a, 0x03, 0x73, 0x74, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x73, 0x74, 0x72, 0x42, 0x2f, 0x0a, 0x21, 0x69, 0x6f, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x68, 0x64, 0x66, 0x67, 0x31, 0x35, 0x39, 0x2e, 0x67, 0x61, 0x6d, 0x65, 0x2e, 0x64, 0x6f,
	0x6d, 0x61, 0x69, 0x6e, 0x2e, 0x64, 0x74, 0x6f, 0x5a, 0x0a, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x75, 0x66, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_protos_protobuf_GameMessage_proto_rawDescOnce sync.Once
	file_protos_protobuf_GameMessage_proto_rawDescData = file_protos_protobuf_GameMessage_proto_rawDesc
)

func file_protos_protobuf_GameMessage_proto_rawDescGZIP() []byte {
	file_protos_protobuf_GameMessage_proto_rawDescOnce.Do(func() {
		file_protos_protobuf_GameMessage_proto_rawDescData = protoimpl.X.CompressGZIP(file_protos_protobuf_GameMessage_proto_rawDescData)
	})
	return file_protos_protobuf_GameMessage_proto_rawDescData
}

var file_protos_protobuf_GameMessage_proto_msgTypes = make([]protoimpl.MessageInfo, 11)
var file_protos_protobuf_GameMessage_proto_goTypes = []interface{}{
	(*Message)(nil),      // 0: GameMessage.Message
	(*HeartBeatReq)(nil), // 1: GameMessage.HeartBeatReq
	(*HeartBeatRes)(nil), // 2: GameMessage.HeartBeatRes
	(*RegisterReq)(nil),  // 3: GameMessage.RegisterReq
	(*RegisterRes)(nil),  // 4: GameMessage.RegisterRes
	(*LoginReq)(nil),     // 5: GameMessage.LoginReq
	(*LoginRes)(nil),     // 6: GameMessage.LoginRes
	(*OfflineReq)(nil),   // 7: GameMessage.OfflineReq
	(*OfflineRes)(nil),   // 8: GameMessage.OfflineRes
	(*TestReq)(nil),      // 9: GameMessage.TestReq
	(*TestRes)(nil),      // 10: GameMessage.TestRes
}
var file_protos_protobuf_GameMessage_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_protos_protobuf_GameMessage_proto_init() }
func file_protos_protobuf_GameMessage_proto_init() {
	if File_protos_protobuf_GameMessage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_protos_protobuf_GameMessage_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Message); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartBeatReq); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*HeartBeatRes); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterReq); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RegisterRes); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginReq); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LoginRes); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OfflineReq); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*OfflineRes); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestReq); i {
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
		file_protos_protobuf_GameMessage_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TestRes); i {
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
			RawDescriptor: file_protos_protobuf_GameMessage_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   11,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_protos_protobuf_GameMessage_proto_goTypes,
		DependencyIndexes: file_protos_protobuf_GameMessage_proto_depIdxs,
		MessageInfos:      file_protos_protobuf_GameMessage_proto_msgTypes,
	}.Build()
	File_protos_protobuf_GameMessage_proto = out.File
	file_protos_protobuf_GameMessage_proto_rawDesc = nil
	file_protos_protobuf_GameMessage_proto_goTypes = nil
	file_protos_protobuf_GameMessage_proto_depIdxs = nil
}
