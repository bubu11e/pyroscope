// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0-devel
// 	protoc        (unknown)
// source: push/v1/push.proto

package pushv1

import (
	v1 "github.com/grafana/fire/pkg/gen/common/v1"
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

type PushResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushResponse) Reset() {
	*x = PushResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_push_v1_push_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushResponse) ProtoMessage() {}

func (x *PushResponse) ProtoReflect() protoreflect.Message {
	mi := &file_push_v1_push_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushResponse.ProtoReflect.Descriptor instead.
func (*PushResponse) Descriptor() ([]byte, []int) {
	return file_push_v1_push_proto_rawDescGZIP(), []int{0}
}

// WriteRawRequest writes a pprof profile
type PushRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// series is a set raw pprof profiles and accompanying labels
	Series []*RawProfileSeries `protobuf:"bytes,1,rep,name=series,proto3" json:"series,omitempty"`
}

func (x *PushRequest) Reset() {
	*x = PushRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_push_v1_push_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushRequest) ProtoMessage() {}

func (x *PushRequest) ProtoReflect() protoreflect.Message {
	mi := &file_push_v1_push_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushRequest.ProtoReflect.Descriptor instead.
func (*PushRequest) Descriptor() ([]byte, []int) {
	return file_push_v1_push_proto_rawDescGZIP(), []int{1}
}

func (x *PushRequest) GetSeries() []*RawProfileSeries {
	if x != nil {
		return x.Series
	}
	return nil
}

// RawProfileSeries represents the pprof profile and its associated labels
type RawProfileSeries struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// LabelPair is the key value pairs to identify the corresponding profile
	Labels []*v1.LabelPair `protobuf:"bytes,1,rep,name=labels,proto3" json:"labels,omitempty"`
	// samples are the set of profile bytes
	Samples []*RawSample `protobuf:"bytes,2,rep,name=samples,proto3" json:"samples,omitempty"`
}

func (x *RawProfileSeries) Reset() {
	*x = RawProfileSeries{}
	if protoimpl.UnsafeEnabled {
		mi := &file_push_v1_push_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawProfileSeries) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawProfileSeries) ProtoMessage() {}

func (x *RawProfileSeries) ProtoReflect() protoreflect.Message {
	mi := &file_push_v1_push_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawProfileSeries.ProtoReflect.Descriptor instead.
func (*RawProfileSeries) Descriptor() ([]byte, []int) {
	return file_push_v1_push_proto_rawDescGZIP(), []int{2}
}

func (x *RawProfileSeries) GetLabels() []*v1.LabelPair {
	if x != nil {
		return x.Labels
	}
	return nil
}

func (x *RawProfileSeries) GetSamples() []*RawSample {
	if x != nil {
		return x.Samples
	}
	return nil
}

// RawSample is the set of bytes that correspond to a pprof profile
type RawSample struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// raw_profile is the set of bytes of the pprof profile
	RawProfile []byte `protobuf:"bytes,1,opt,name=raw_profile,json=rawProfile,proto3" json:"raw_profile,omitempty"`
}

func (x *RawSample) Reset() {
	*x = RawSample{}
	if protoimpl.UnsafeEnabled {
		mi := &file_push_v1_push_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *RawSample) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*RawSample) ProtoMessage() {}

func (x *RawSample) ProtoReflect() protoreflect.Message {
	mi := &file_push_v1_push_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use RawSample.ProtoReflect.Descriptor instead.
func (*RawSample) Descriptor() ([]byte, []int) {
	return file_push_v1_push_proto_rawDescGZIP(), []int{3}
}

func (x *RawSample) GetRawProfile() []byte {
	if x != nil {
		return x.RawProfile
	}
	return nil
}

var File_push_v1_push_proto protoreflect.FileDescriptor

var file_push_v1_push_proto_rawDesc = []byte{
	0x0a, 0x12, 0x70, 0x75, 0x73, 0x68, 0x2f, 0x76, 0x31, 0x2f, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x76, 0x31, 0x1a, 0x16, 0x63,
	0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x63, 0x6f, 0x6d, 0x6d, 0x6f, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0e, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x40, 0x0a, 0x0b, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71,
	0x75, 0x65, 0x73, 0x74, 0x12, 0x31, 0x0a, 0x06, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x19, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x52,
	0x61, 0x77, 0x50, 0x72, 0x6f, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x52,
	0x06, 0x73, 0x65, 0x72, 0x69, 0x65, 0x73, 0x22, 0x6e, 0x0a, 0x10, 0x52, 0x61, 0x77, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x53, 0x65, 0x72, 0x69, 0x65, 0x73, 0x12, 0x2c, 0x0a, 0x06, 0x6c,
	0x61, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x63, 0x6f,
	0x6d, 0x6d, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x4c, 0x61, 0x62, 0x65, 0x6c, 0x50, 0x61, 0x69,
	0x72, 0x52, 0x06, 0x6c, 0x61, 0x62, 0x65, 0x6c, 0x73, 0x12, 0x2c, 0x0a, 0x07, 0x73, 0x61, 0x6d,
	0x70, 0x6c, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x75, 0x73,
	0x68, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x61, 0x77, 0x53, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x52, 0x07,
	0x73, 0x61, 0x6d, 0x70, 0x6c, 0x65, 0x73, 0x22, 0x2c, 0x0a, 0x09, 0x52, 0x61, 0x77, 0x53, 0x61,
	0x6d, 0x70, 0x6c, 0x65, 0x12, 0x1f, 0x0a, 0x0b, 0x72, 0x61, 0x77, 0x5f, 0x70, 0x72, 0x6f, 0x66,
	0x69, 0x6c, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x0a, 0x72, 0x61, 0x77, 0x50, 0x72,
	0x6f, 0x66, 0x69, 0x6c, 0x65, 0x32, 0x3f, 0x0a, 0x06, 0x50, 0x75, 0x73, 0x68, 0x65, 0x72, 0x12,
	0x35, 0x0a, 0x04, 0x50, 0x75, 0x73, 0x68, 0x12, 0x14, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x76,
	0x31, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e,
	0x70, 0x75, 0x73, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x32, 0x41, 0x0a, 0x08, 0x49, 0x6e, 0x67, 0x65, 0x73, 0x74,
	0x65, 0x72, 0x12, 0x35, 0x0a, 0x04, 0x50, 0x75, 0x73, 0x68, 0x12, 0x14, 0x2e, 0x70, 0x75, 0x73,
	0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x15, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x76, 0x31, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x22, 0x00, 0x42, 0x85, 0x01, 0x0a, 0x0b, 0x63, 0x6f,
	0x6d, 0x2e, 0x70, 0x75, 0x73, 0x68, 0x2e, 0x76, 0x31, 0x42, 0x09, 0x50, 0x75, 0x73, 0x68, 0x50,
	0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63,
	0x6f, 0x6d, 0x2f, 0x67, 0x72, 0x61, 0x66, 0x61, 0x6e, 0x61, 0x2f, 0x66, 0x69, 0x72, 0x65, 0x2f,
	0x70, 0x6b, 0x67, 0x2f, 0x67, 0x65, 0x6e, 0x2f, 0x70, 0x75, 0x73, 0x68, 0x2f, 0x76, 0x31, 0x3b,
	0x70, 0x75, 0x73, 0x68, 0x76, 0x31, 0xa2, 0x02, 0x03, 0x50, 0x58, 0x58, 0xaa, 0x02, 0x07, 0x50,
	0x75, 0x73, 0x68, 0x2e, 0x56, 0x31, 0xca, 0x02, 0x07, 0x50, 0x75, 0x73, 0x68, 0x5c, 0x56, 0x31,
	0xe2, 0x02, 0x13, 0x50, 0x75, 0x73, 0x68, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65,
	0x74, 0x61, 0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x08, 0x50, 0x75, 0x73, 0x68, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_push_v1_push_proto_rawDescOnce sync.Once
	file_push_v1_push_proto_rawDescData = file_push_v1_push_proto_rawDesc
)

func file_push_v1_push_proto_rawDescGZIP() []byte {
	file_push_v1_push_proto_rawDescOnce.Do(func() {
		file_push_v1_push_proto_rawDescData = protoimpl.X.CompressGZIP(file_push_v1_push_proto_rawDescData)
	})
	return file_push_v1_push_proto_rawDescData
}

var file_push_v1_push_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_push_v1_push_proto_goTypes = []interface{}{
	(*PushResponse)(nil),     // 0: push.v1.PushResponse
	(*PushRequest)(nil),      // 1: push.v1.PushRequest
	(*RawProfileSeries)(nil), // 2: push.v1.RawProfileSeries
	(*RawSample)(nil),        // 3: push.v1.RawSample
	(*v1.LabelPair)(nil),     // 4: common.v1.LabelPair
}
var file_push_v1_push_proto_depIdxs = []int32{
	2, // 0: push.v1.PushRequest.series:type_name -> push.v1.RawProfileSeries
	4, // 1: push.v1.RawProfileSeries.labels:type_name -> common.v1.LabelPair
	3, // 2: push.v1.RawProfileSeries.samples:type_name -> push.v1.RawSample
	1, // 3: push.v1.Pusher.Push:input_type -> push.v1.PushRequest
	1, // 4: push.v1.Ingester.Push:input_type -> push.v1.PushRequest
	0, // 5: push.v1.Pusher.Push:output_type -> push.v1.PushResponse
	0, // 6: push.v1.Ingester.Push:output_type -> push.v1.PushResponse
	5, // [5:7] is the sub-list for method output_type
	3, // [3:5] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_push_v1_push_proto_init() }
func file_push_v1_push_proto_init() {
	if File_push_v1_push_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_push_v1_push_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushResponse); i {
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
		file_push_v1_push_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushRequest); i {
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
		file_push_v1_push_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RawProfileSeries); i {
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
		file_push_v1_push_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*RawSample); i {
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
			RawDescriptor: file_push_v1_push_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   2,
		},
		GoTypes:           file_push_v1_push_proto_goTypes,
		DependencyIndexes: file_push_v1_push_proto_depIdxs,
		MessageInfos:      file_push_v1_push_proto_msgTypes,
	}.Build()
	File_push_v1_push_proto = out.File
	file_push_v1_push_proto_rawDesc = nil
	file_push_v1_push_proto_goTypes = nil
	file_push_v1_push_proto_depIdxs = nil
}
