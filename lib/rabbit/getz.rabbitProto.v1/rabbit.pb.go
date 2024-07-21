// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.26.1
// source: rabbit.proto

package rabbitDefines

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

type ISBNMessage struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Isbn string `protobuf:"bytes,1,opt,name=isbn,proto3" json:"isbn,omitempty"`
}

func (x *ISBNMessage) Reset() {
	*x = ISBNMessage{}
	if protoimpl.UnsafeEnabled {
		mi := &file_rabbit_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ISBNMessage) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ISBNMessage) ProtoMessage() {}

func (x *ISBNMessage) ProtoReflect() protoreflect.Message {
	mi := &file_rabbit_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ISBNMessage.ProtoReflect.Descriptor instead.
func (*ISBNMessage) Descriptor() ([]byte, []int) {
	return file_rabbit_proto_rawDescGZIP(), []int{0}
}

func (x *ISBNMessage) GetIsbn() string {
	if x != nil {
		return x.Isbn
	}
	return ""
}

var File_rabbit_proto protoreflect.FileDescriptor

var file_rabbit_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x72, 0x61, 0x62, 0x62, 0x69, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x02,
	0x70, 0x62, 0x22, 0x21, 0x0a, 0x0b, 0x49, 0x53, 0x42, 0x4e, 0x4d, 0x65, 0x73, 0x73, 0x61, 0x67,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x69, 0x73, 0x62, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x69, 0x73, 0x62, 0x6e, 0x42, 0x23, 0x5a, 0x21, 0x67, 0x65, 0x74, 0x7a, 0x2e, 0x72, 0x61,
	0x62, 0x62, 0x69, 0x74, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x76, 0x31, 0x3b, 0x72, 0x61, 0x62,
	0x62, 0x69, 0x74, 0x44, 0x65, 0x66, 0x69, 0x6e, 0x65, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_rabbit_proto_rawDescOnce sync.Once
	file_rabbit_proto_rawDescData = file_rabbit_proto_rawDesc
)

func file_rabbit_proto_rawDescGZIP() []byte {
	file_rabbit_proto_rawDescOnce.Do(func() {
		file_rabbit_proto_rawDescData = protoimpl.X.CompressGZIP(file_rabbit_proto_rawDescData)
	})
	return file_rabbit_proto_rawDescData
}

var file_rabbit_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_rabbit_proto_goTypes = []any{
	(*ISBNMessage)(nil), // 0: pb.ISBNMessage
}
var file_rabbit_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_rabbit_proto_init() }
func file_rabbit_proto_init() {
	if File_rabbit_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_rabbit_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*ISBNMessage); i {
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
			RawDescriptor: file_rabbit_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_rabbit_proto_goTypes,
		DependencyIndexes: file_rabbit_proto_depIdxs,
		MessageInfos:      file_rabbit_proto_msgTypes,
	}.Build()
	File_rabbit_proto = out.File
	file_rabbit_proto_rawDesc = nil
	file_rabbit_proto_goTypes = nil
	file_rabbit_proto_depIdxs = nil
}
