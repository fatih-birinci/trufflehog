// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.20.0
// source: custom_detectors.proto

package custom_detectorspb

import (
	_ "github.com/envoyproxy/protoc-gen-validate/validate"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	anypb "google.golang.org/protobuf/types/known/anypb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type CustomDetector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type       string     `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Name       string     `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Connection *anypb.Any `protobuf:"bytes,3,opt,name=connection,proto3" json:"connection,omitempty"`
}

func (x *CustomDetector) Reset() {
	*x = CustomDetector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_custom_detectors_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CustomDetector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CustomDetector) ProtoMessage() {}

func (x *CustomDetector) ProtoReflect() protoreflect.Message {
	mi := &file_custom_detectors_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CustomDetector.ProtoReflect.Descriptor instead.
func (*CustomDetector) Descriptor() ([]byte, []int) {
	return file_custom_detectors_proto_rawDescGZIP(), []int{0}
}

func (x *CustomDetector) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *CustomDetector) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CustomDetector) GetConnection() *anypb.Any {
	if x != nil {
		return x.Connection
	}
	return nil
}

type CustomRegex struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keywords []string          `protobuf:"bytes,1,rep,name=keywords,proto3" json:"keywords,omitempty"`
	Regex    map[string]string `protobuf:"bytes,2,rep,name=regex,proto3" json:"regex,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	Verify   []*VerifierConfig `protobuf:"bytes,3,rep,name=verify,proto3" json:"verify,omitempty"`
}

func (x *CustomRegex) Reset() {
	*x = CustomRegex{}
	if protoimpl.UnsafeEnabled {
		mi := &file_custom_detectors_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CustomRegex) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CustomRegex) ProtoMessage() {}

func (x *CustomRegex) ProtoReflect() protoreflect.Message {
	mi := &file_custom_detectors_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CustomRegex.ProtoReflect.Descriptor instead.
func (*CustomRegex) Descriptor() ([]byte, []int) {
	return file_custom_detectors_proto_rawDescGZIP(), []int{1}
}

func (x *CustomRegex) GetKeywords() []string {
	if x != nil {
		return x.Keywords
	}
	return nil
}

func (x *CustomRegex) GetRegex() map[string]string {
	if x != nil {
		return x.Regex
	}
	return nil
}

func (x *CustomRegex) GetVerify() []*VerifierConfig {
	if x != nil {
		return x.Verify
	}
	return nil
}

type VerifierConfig struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Endpoint      string   `protobuf:"bytes,1,opt,name=endpoint,proto3" json:"endpoint,omitempty"`
	Unsafe        bool     `protobuf:"varint,2,opt,name=unsafe,proto3" json:"unsafe,omitempty"`
	Headers       []string `protobuf:"bytes,3,rep,name=headers,proto3" json:"headers,omitempty"`
	SuccessRanges []string `protobuf:"bytes,4,rep,name=successRanges,proto3" json:"successRanges,omitempty"`
}

func (x *VerifierConfig) Reset() {
	*x = VerifierConfig{}
	if protoimpl.UnsafeEnabled {
		mi := &file_custom_detectors_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VerifierConfig) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VerifierConfig) ProtoMessage() {}

func (x *VerifierConfig) ProtoReflect() protoreflect.Message {
	mi := &file_custom_detectors_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VerifierConfig.ProtoReflect.Descriptor instead.
func (*VerifierConfig) Descriptor() ([]byte, []int) {
	return file_custom_detectors_proto_rawDescGZIP(), []int{2}
}

func (x *VerifierConfig) GetEndpoint() string {
	if x != nil {
		return x.Endpoint
	}
	return ""
}

func (x *VerifierConfig) GetUnsafe() bool {
	if x != nil {
		return x.Unsafe
	}
	return false
}

func (x *VerifierConfig) GetHeaders() []string {
	if x != nil {
		return x.Headers
	}
	return nil
}

func (x *VerifierConfig) GetSuccessRanges() []string {
	if x != nil {
		return x.SuccessRanges
	}
	return nil
}

var File_custom_detectors_proto protoreflect.FileDescriptor

var file_custom_detectors_proto_rawDesc = []byte{
	0x0a, 0x16, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f,
	0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x10, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d,
	0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x1a, 0x19, 0x67, 0x6f, 0x6f, 0x67,
	0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x61, 0x6e, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2f,
	0x76, 0x61, 0x6c, 0x69, 0x64, 0x61, 0x74, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x6e,
	0x0a, 0x0e, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x44, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x12, 0x12, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x34, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x6e,
	0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x67,
	0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x41,
	0x6e, 0x79, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0xdd,
	0x01, 0x0a, 0x0b, 0x43, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x52, 0x65, 0x67, 0x65, 0x78, 0x12, 0x1a,
	0x0a, 0x08, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x08, 0x6b, 0x65, 0x79, 0x77, 0x6f, 0x72, 0x64, 0x73, 0x12, 0x3e, 0x0a, 0x05, 0x72, 0x65,
	0x67, 0x65, 0x78, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x28, 0x2e, 0x63, 0x75, 0x73, 0x74,
	0x6f, 0x6d, 0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x43, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x52, 0x65, 0x67, 0x65, 0x78, 0x2e, 0x52, 0x65, 0x67, 0x65, 0x78, 0x45, 0x6e,
	0x74, 0x72, 0x79, 0x52, 0x05, 0x72, 0x65, 0x67, 0x65, 0x78, 0x12, 0x38, 0x0a, 0x06, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x20, 0x2e, 0x63, 0x75, 0x73,
	0x74, 0x6f, 0x6d, 0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x2e, 0x56, 0x65,
	0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69, 0x67, 0x52, 0x06, 0x76, 0x65,
	0x72, 0x69, 0x66, 0x79, 0x1a, 0x38, 0x0a, 0x0a, 0x52, 0x65, 0x67, 0x65, 0x78, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x8e,
	0x01, 0x0a, 0x0e, 0x56, 0x65, 0x72, 0x69, 0x66, 0x69, 0x65, 0x72, 0x43, 0x6f, 0x6e, 0x66, 0x69,
	0x67, 0x12, 0x24, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x42, 0x08, 0xfa, 0x42, 0x05, 0x72, 0x03, 0x90, 0x01, 0x01, 0x52, 0x08, 0x65,
	0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x75, 0x6e, 0x73, 0x61, 0x66,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x06, 0x75, 0x6e, 0x73, 0x61, 0x66, 0x65, 0x12,
	0x18, 0x0a, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x07, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x73, 0x12, 0x24, 0x0a, 0x0d, 0x73, 0x75, 0x63,
	0x63, 0x65, 0x73, 0x73, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x09,
	0x52, 0x0d, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x42,
	0x44, 0x5a, 0x42, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x74, 0x72,
	0x75, 0x66, 0x66, 0x6c, 0x65, 0x73, 0x65, 0x63, 0x75, 0x72, 0x69, 0x74, 0x79, 0x2f, 0x74, 0x72,
	0x75, 0x66, 0x66, 0x6c, 0x65, 0x68, 0x6f, 0x67, 0x2f, 0x76, 0x33, 0x2f, 0x70, 0x6b, 0x67, 0x2f,
	0x70, 0x62, 0x2f, 0x63, 0x75, 0x73, 0x74, 0x6f, 0x6d, 0x5f, 0x64, 0x65, 0x74, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x73, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_custom_detectors_proto_rawDescOnce sync.Once
	file_custom_detectors_proto_rawDescData = file_custom_detectors_proto_rawDesc
)

func file_custom_detectors_proto_rawDescGZIP() []byte {
	file_custom_detectors_proto_rawDescOnce.Do(func() {
		file_custom_detectors_proto_rawDescData = protoimpl.X.CompressGZIP(file_custom_detectors_proto_rawDescData)
	})
	return file_custom_detectors_proto_rawDescData
}

var file_custom_detectors_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_custom_detectors_proto_goTypes = []interface{}{
	(*CustomDetector)(nil), // 0: custom_detectors.CustomDetector
	(*CustomRegex)(nil),    // 1: custom_detectors.CustomRegex
	(*VerifierConfig)(nil), // 2: custom_detectors.VerifierConfig
	nil,                    // 3: custom_detectors.CustomRegex.RegexEntry
	(*anypb.Any)(nil),      // 4: google.protobuf.Any
}
var file_custom_detectors_proto_depIdxs = []int32{
	4, // 0: custom_detectors.CustomDetector.connection:type_name -> google.protobuf.Any
	3, // 1: custom_detectors.CustomRegex.regex:type_name -> custom_detectors.CustomRegex.RegexEntry
	2, // 2: custom_detectors.CustomRegex.verify:type_name -> custom_detectors.VerifierConfig
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_custom_detectors_proto_init() }
func file_custom_detectors_proto_init() {
	if File_custom_detectors_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_custom_detectors_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CustomDetector); i {
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
		file_custom_detectors_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CustomRegex); i {
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
		file_custom_detectors_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VerifierConfig); i {
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
			RawDescriptor: file_custom_detectors_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_custom_detectors_proto_goTypes,
		DependencyIndexes: file_custom_detectors_proto_depIdxs,
		MessageInfos:      file_custom_detectors_proto_msgTypes,
	}.Build()
	File_custom_detectors_proto = out.File
	file_custom_detectors_proto_rawDesc = nil
	file_custom_detectors_proto_goTypes = nil
	file_custom_detectors_proto_depIdxs = nil
}
