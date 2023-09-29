// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.13.0
// source: api/v1/diagnostic.proto

package v1

import (
	reflect "reflect"
	sync "sync"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type ScreenshotRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Authn  *Authn  `protobuf:"bytes,1,opt,name=authn,proto3" json:"authn,omitempty"`
	Vendor *Vendor `protobuf:"bytes,2,opt,name=vendor,proto3" json:"vendor,omitempty"`
}

func (x *ScreenshotRequest) Reset() {
	*x = ScreenshotRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_diagnostic_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScreenshotRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScreenshotRequest) ProtoMessage() {}

func (x *ScreenshotRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_diagnostic_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScreenshotRequest.ProtoReflect.Descriptor instead.
func (*ScreenshotRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_diagnostic_proto_rawDescGZIP(), []int{0}
}

func (x *ScreenshotRequest) GetAuthn() *Authn {
	if x != nil {
		return x.Authn
	}
	return nil
}

func (x *ScreenshotRequest) GetVendor() *Vendor {
	if x != nil {
		return x.Vendor
	}
	return nil
}

type ScreenshotResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Image    []byte `protobuf:"bytes,1,opt,name=image,proto3" json:"image,omitempty"`
	Filetype string `protobuf:"bytes,2,opt,name=filetype,proto3" json:"filetype,omitempty"`
}

func (x *ScreenshotResponse) Reset() {
	*x = ScreenshotResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_diagnostic_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ScreenshotResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ScreenshotResponse) ProtoMessage() {}

func (x *ScreenshotResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_diagnostic_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ScreenshotResponse.ProtoReflect.Descriptor instead.
func (*ScreenshotResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_diagnostic_proto_rawDescGZIP(), []int{1}
}

func (x *ScreenshotResponse) GetImage() []byte {
	if x != nil {
		return x.Image
	}
	return nil
}

func (x *ScreenshotResponse) GetFiletype() string {
	if x != nil {
		return x.Filetype
	}
	return ""
}

type ClearSystemEventLogRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Authn  *Authn  `protobuf:"bytes,1,opt,name=authn,proto3" json:"authn,omitempty"`
	Vendor *Vendor `protobuf:"bytes,2,opt,name=vendor,proto3" json:"vendor,omitempty"`
}

func (x *ClearSystemEventLogRequest) Reset() {
	*x = ClearSystemEventLogRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_diagnostic_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClearSystemEventLogRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClearSystemEventLogRequest) ProtoMessage() {}

func (x *ClearSystemEventLogRequest) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_diagnostic_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClearSystemEventLogRequest.ProtoReflect.Descriptor instead.
func (*ClearSystemEventLogRequest) Descriptor() ([]byte, []int) {
	return file_api_v1_diagnostic_proto_rawDescGZIP(), []int{2}
}

func (x *ClearSystemEventLogRequest) GetAuthn() *Authn {
	if x != nil {
		return x.Authn
	}
	return nil
}

func (x *ClearSystemEventLogRequest) GetVendor() *Vendor {
	if x != nil {
		return x.Vendor
	}
	return nil
}

type ClearSystemEventLogResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TaskId string `protobuf:"bytes,1,opt,name=task_id,json=taskId,proto3" json:"task_id,omitempty"`
}

func (x *ClearSystemEventLogResponse) Reset() {
	*x = ClearSystemEventLogResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_api_v1_diagnostic_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ClearSystemEventLogResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ClearSystemEventLogResponse) ProtoMessage() {}

func (x *ClearSystemEventLogResponse) ProtoReflect() protoreflect.Message {
	mi := &file_api_v1_diagnostic_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ClearSystemEventLogResponse.ProtoReflect.Descriptor instead.
func (*ClearSystemEventLogResponse) Descriptor() ([]byte, []int) {
	return file_api_v1_diagnostic_proto_rawDescGZIP(), []int{3}
}

func (x *ClearSystemEventLogResponse) GetTaskId() string {
	if x != nil {
		return x.TaskId
	}
	return ""
}

var File_api_v1_diagnostic_proto protoreflect.FileDescriptor

var file_api_v1_diagnostic_proto_rawDesc = []byte{
	0x0a, 0x17, 0x61, 0x70, 0x69, 0x2f, 0x76, 0x31, 0x2f, 0x64, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73,
	0x74, 0x69, 0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x21, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x69, 0x6e, 0x6b, 0x65, 0x72, 0x62, 0x65, 0x6c, 0x6c,
	0x2e, 0x70, 0x62, 0x6e, 0x6a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x1a, 0x13, 0x61, 0x70,
	0x69, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x22, 0x96, 0x01, 0x0a, 0x11, 0x53, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x73, 0x68, 0x6f, 0x74,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x3e, 0x0a, 0x05, 0x61, 0x75, 0x74, 0x68, 0x6e,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x69, 0x6e, 0x6b, 0x65, 0x72, 0x62, 0x65, 0x6c, 0x6c, 0x2e, 0x70,
	0x62, 0x6e, 0x6a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6e,
	0x52, 0x05, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x12, 0x41, 0x0a, 0x06, 0x76, 0x65, 0x6e, 0x64, 0x6f,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x29, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x69, 0x6e, 0x6b, 0x65, 0x72, 0x62, 0x65, 0x6c, 0x6c, 0x2e,
	0x70, 0x62, 0x6e, 0x6a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x6e, 0x64,
	0x6f, 0x72, 0x52, 0x06, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x22, 0x46, 0x0a, 0x12, 0x53, 0x63,
	0x72, 0x65, 0x65, 0x6e, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52,
	0x05, 0x69, 0x6d, 0x61, 0x67, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x74, 0x79,
	0x70, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x66, 0x69, 0x6c, 0x65, 0x74, 0x79,
	0x70, 0x65, 0x22, 0x9f, 0x01, 0x0a, 0x1a, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x53, 0x79, 0x73, 0x74,
	0x65, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x3e, 0x0a, 0x05, 0x61, 0x75, 0x74, 0x68, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x28, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x69,
	0x6e, 0x6b, 0x65, 0x72, 0x62, 0x65, 0x6c, 0x6c, 0x2e, 0x70, 0x62, 0x6e, 0x6a, 0x2e, 0x61, 0x70,
	0x69, 0x2e, 0x76, 0x31, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x6e, 0x52, 0x05, 0x61, 0x75, 0x74, 0x68,
	0x6e, 0x12, 0x41, 0x0a, 0x06, 0x76, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x29, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x74,
	0x69, 0x6e, 0x6b, 0x65, 0x72, 0x62, 0x65, 0x6c, 0x6c, 0x2e, 0x70, 0x62, 0x6e, 0x6a, 0x2e, 0x61,
	0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x65, 0x6e, 0x64, 0x6f, 0x72, 0x52, 0x06, 0x76, 0x65,
	0x6e, 0x64, 0x6f, 0x72, 0x22, 0x36, 0x0a, 0x1b, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x53, 0x79, 0x73,
	0x74, 0x65, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x17, 0x0a, 0x07, 0x74, 0x61, 0x73, 0x6b, 0x5f, 0x69, 0x64, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x74, 0x61, 0x73, 0x6b, 0x49, 0x64, 0x32, 0x9e, 0x02, 0x0a,
	0x0a, 0x44, 0x69, 0x61, 0x67, 0x6e, 0x6f, 0x73, 0x74, 0x69, 0x63, 0x12, 0x79, 0x0a, 0x0a, 0x53,
	0x63, 0x72, 0x65, 0x65, 0x6e, 0x73, 0x68, 0x6f, 0x74, 0x12, 0x34, 0x2e, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x69, 0x6e, 0x6b, 0x65, 0x72, 0x62, 0x65, 0x6c,
	0x6c, 0x2e, 0x70, 0x62, 0x6e, 0x6a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63,
	0x72, 0x65, 0x65, 0x6e, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a,
	0x35, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x69, 0x6e,
	0x6b, 0x65, 0x72, 0x62, 0x65, 0x6c, 0x6c, 0x2e, 0x70, 0x62, 0x6e, 0x6a, 0x2e, 0x61, 0x70, 0x69,
	0x2e, 0x76, 0x31, 0x2e, 0x53, 0x63, 0x72, 0x65, 0x65, 0x6e, 0x73, 0x68, 0x6f, 0x74, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x94, 0x01, 0x0a, 0x13, 0x43, 0x6c, 0x65, 0x61, 0x72,
	0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x45, 0x76, 0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x12, 0x3d,
	0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x69, 0x6e, 0x6b,
	0x65, 0x72, 0x62, 0x65, 0x6c, 0x6c, 0x2e, 0x70, 0x62, 0x6e, 0x6a, 0x2e, 0x61, 0x70, 0x69, 0x2e,
	0x76, 0x31, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x45, 0x76,
	0x65, 0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x3e, 0x2e,
	0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2e, 0x74, 0x69, 0x6e, 0x6b, 0x65,
	0x72, 0x62, 0x65, 0x6c, 0x6c, 0x2e, 0x70, 0x62, 0x6e, 0x6a, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x76,
	0x31, 0x2e, 0x43, 0x6c, 0x65, 0x61, 0x72, 0x53, 0x79, 0x73, 0x74, 0x65, 0x6d, 0x45, 0x76, 0x65,
	0x6e, 0x74, 0x4c, 0x6f, 0x67, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x33, 0x5a,
	0x21, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x69, 0x6e, 0x6b,
	0x65, 0x72, 0x62, 0x65, 0x6c, 0x6c, 0x2f, 0x70, 0x62, 0x6e, 0x6a, 0x2f, 0x61, 0x70, 0x69, 0x2f,
	0x76, 0x31, 0xea, 0x02, 0x0d, 0x50, 0x62, 0x6e, 0x6a, 0x3a, 0x3a, 0x41, 0x70, 0x69, 0x3a, 0x3a,
	0x56, 0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_api_v1_diagnostic_proto_rawDescOnce sync.Once
	file_api_v1_diagnostic_proto_rawDescData = file_api_v1_diagnostic_proto_rawDesc
)

func file_api_v1_diagnostic_proto_rawDescGZIP() []byte {
	file_api_v1_diagnostic_proto_rawDescOnce.Do(func() {
		file_api_v1_diagnostic_proto_rawDescData = protoimpl.X.CompressGZIP(file_api_v1_diagnostic_proto_rawDescData)
	})
	return file_api_v1_diagnostic_proto_rawDescData
}

var file_api_v1_diagnostic_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_api_v1_diagnostic_proto_goTypes = []interface{}{
	(*ScreenshotRequest)(nil),           // 0: github.com.tinkerbell.pbnj.api.v1.ScreenshotRequest
	(*ScreenshotResponse)(nil),          // 1: github.com.tinkerbell.pbnj.api.v1.ScreenshotResponse
	(*ClearSystemEventLogRequest)(nil),  // 2: github.com.tinkerbell.pbnj.api.v1.ClearSystemEventLogRequest
	(*ClearSystemEventLogResponse)(nil), // 3: github.com.tinkerbell.pbnj.api.v1.ClearSystemEventLogResponse
	(*Authn)(nil),                       // 4: github.com.tinkerbell.pbnj.api.v1.Authn
	(*Vendor)(nil),                      // 5: github.com.tinkerbell.pbnj.api.v1.Vendor
}
var file_api_v1_diagnostic_proto_depIdxs = []int32{
	4, // 0: github.com.tinkerbell.pbnj.api.v1.ScreenshotRequest.authn:type_name -> github.com.tinkerbell.pbnj.api.v1.Authn
	5, // 1: github.com.tinkerbell.pbnj.api.v1.ScreenshotRequest.vendor:type_name -> github.com.tinkerbell.pbnj.api.v1.Vendor
	4, // 2: github.com.tinkerbell.pbnj.api.v1.ClearSystemEventLogRequest.authn:type_name -> github.com.tinkerbell.pbnj.api.v1.Authn
	5, // 3: github.com.tinkerbell.pbnj.api.v1.ClearSystemEventLogRequest.vendor:type_name -> github.com.tinkerbell.pbnj.api.v1.Vendor
	0, // 4: github.com.tinkerbell.pbnj.api.v1.Diagnostic.Screenshot:input_type -> github.com.tinkerbell.pbnj.api.v1.ScreenshotRequest
	2, // 5: github.com.tinkerbell.pbnj.api.v1.Diagnostic.ClearSystemEventLog:input_type -> github.com.tinkerbell.pbnj.api.v1.ClearSystemEventLogRequest
	1, // 6: github.com.tinkerbell.pbnj.api.v1.Diagnostic.Screenshot:output_type -> github.com.tinkerbell.pbnj.api.v1.ScreenshotResponse
	3, // 7: github.com.tinkerbell.pbnj.api.v1.Diagnostic.ClearSystemEventLog:output_type -> github.com.tinkerbell.pbnj.api.v1.ClearSystemEventLogResponse
	6, // [6:8] is the sub-list for method output_type
	4, // [4:6] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_api_v1_diagnostic_proto_init() }
func file_api_v1_diagnostic_proto_init() {
	if File_api_v1_diagnostic_proto != nil {
		return
	}
	file_api_v1_common_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_api_v1_diagnostic_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScreenshotRequest); i {
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
		file_api_v1_diagnostic_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ScreenshotResponse); i {
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
		file_api_v1_diagnostic_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClearSystemEventLogRequest); i {
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
		file_api_v1_diagnostic_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ClearSystemEventLogResponse); i {
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
			RawDescriptor: file_api_v1_diagnostic_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_api_v1_diagnostic_proto_goTypes,
		DependencyIndexes: file_api_v1_diagnostic_proto_depIdxs,
		MessageInfos:      file_api_v1_diagnostic_proto_msgTypes,
	}.Build()
	File_api_v1_diagnostic_proto = out.File
	file_api_v1_diagnostic_proto_rawDesc = nil
	file_api_v1_diagnostic_proto_goTypes = nil
	file_api_v1_diagnostic_proto_depIdxs = nil
}
