// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.24.4
// source: in_toto_attestation/predicates/test_result/v0/test_result.proto

package v0

import (
	v1 "github.com/in-toto/attestation/go/v1"
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

type TestResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Result        string                   `protobuf:"bytes,1,opt,name=result,proto3" json:"result,omitempty"`
	Configuration []*v1.ResourceDescriptor `protobuf:"bytes,2,rep,name=configuration,proto3" json:"configuration,omitempty"`
	Url           string                   `protobuf:"bytes,3,opt,name=url,proto3" json:"url,omitempty"`
	PassedTests   []string                 `protobuf:"bytes,4,rep,name=passed_tests,json=passedTests,proto3" json:"passed_tests,omitempty"`
	WarnedTests   []string                 `protobuf:"bytes,5,rep,name=warned_tests,json=warnedTests,proto3" json:"warned_tests,omitempty"`
	FailedTests   []string                 `protobuf:"bytes,6,rep,name=failed_tests,json=failedTests,proto3" json:"failed_tests,omitempty"`
}

func (x *TestResult) Reset() {
	*x = TestResult{}
	mi := &file_in_toto_attestation_predicates_test_result_v0_test_result_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *TestResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TestResult) ProtoMessage() {}

func (x *TestResult) ProtoReflect() protoreflect.Message {
	mi := &file_in_toto_attestation_predicates_test_result_v0_test_result_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TestResult.ProtoReflect.Descriptor instead.
func (*TestResult) Descriptor() ([]byte, []int) {
	return file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDescGZIP(), []int{0}
}

func (x *TestResult) GetResult() string {
	if x != nil {
		return x.Result
	}
	return ""
}

func (x *TestResult) GetConfiguration() []*v1.ResourceDescriptor {
	if x != nil {
		return x.Configuration
	}
	return nil
}

func (x *TestResult) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

func (x *TestResult) GetPassedTests() []string {
	if x != nil {
		return x.PassedTests
	}
	return nil
}

func (x *TestResult) GetWarnedTests() []string {
	if x != nil {
		return x.WarnedTests
	}
	return nil
}

func (x *TestResult) GetFailedTests() []string {
	if x != nil {
		return x.FailedTests
	}
	return nil
}

var File_in_toto_attestation_predicates_test_result_v0_test_result_proto protoreflect.FileDescriptor

var file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDesc = []byte{
	0x0a, 0x3f, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73,
	0x2f, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2f, 0x76, 0x30, 0x2f,
	0x74, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x12, 0x2d, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x73, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x76, 0x30,
	0x1a, 0x30, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76, 0x31, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x22, 0xf1, 0x01, 0x0a, 0x0a, 0x54, 0x65, 0x73, 0x74, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x50, 0x0a, 0x0d, 0x63, 0x6f, 0x6e,
	0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x2a, 0x2e, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x0d, 0x63, 0x6f,
	0x6e, 0x66, 0x69, 0x67, 0x75, 0x72, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x10, 0x0a, 0x03, 0x75,
	0x72, 0x6c, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x75, 0x72, 0x6c, 0x12, 0x21, 0x0a,
	0x0c, 0x70, 0x61, 0x73, 0x73, 0x65, 0x64, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x73, 0x18, 0x04, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x0b, 0x70, 0x61, 0x73, 0x73, 0x65, 0x64, 0x54, 0x65, 0x73, 0x74, 0x73,
	0x12, 0x21, 0x0a, 0x0c, 0x77, 0x61, 0x72, 0x6e, 0x65, 0x64, 0x5f, 0x74, 0x65, 0x73, 0x74, 0x73,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x77, 0x61, 0x72, 0x6e, 0x65, 0x64, 0x54, 0x65,
	0x73, 0x74, 0x73, 0x12, 0x21, 0x0a, 0x0c, 0x66, 0x61, 0x69, 0x6c, 0x65, 0x64, 0x5f, 0x74, 0x65,
	0x73, 0x74, 0x73, 0x18, 0x06, 0x20, 0x03, 0x28, 0x09, 0x52, 0x0b, 0x66, 0x61, 0x69, 0x6c, 0x65,
	0x64, 0x54, 0x65, 0x73, 0x74, 0x73, 0x42, 0x75, 0x0a, 0x36, 0x69, 0x6f, 0x2e, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x69, 0x6e, 0x74, 0x6f, 0x74, 0x6f, 0x2e, 0x61, 0x74, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65,
	0x73, 0x2e, 0x74, 0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2e, 0x76, 0x30,
	0x5a, 0x3b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x2d,
	0x74, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x2f, 0x74,
	0x65, 0x73, 0x74, 0x5f, 0x72, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x2f, 0x76, 0x30, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDescOnce sync.Once
	file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDescData = file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDesc
)

func file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDescGZIP() []byte {
	file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDescOnce.Do(func() {
		file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDescData = protoimpl.X.CompressGZIP(file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDescData)
	})
	return file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDescData
}

var file_in_toto_attestation_predicates_test_result_v0_test_result_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_in_toto_attestation_predicates_test_result_v0_test_result_proto_goTypes = []any{
	(*TestResult)(nil),            // 0: in_toto_attestation.predicates.test_result.v0.TestResult
	(*v1.ResourceDescriptor)(nil), // 1: in_toto_attestation.v1.ResourceDescriptor
}
var file_in_toto_attestation_predicates_test_result_v0_test_result_proto_depIdxs = []int32{
	1, // 0: in_toto_attestation.predicates.test_result.v0.TestResult.configuration:type_name -> in_toto_attestation.v1.ResourceDescriptor
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_in_toto_attestation_predicates_test_result_v0_test_result_proto_init() }
func file_in_toto_attestation_predicates_test_result_v0_test_result_proto_init() {
	if File_in_toto_attestation_predicates_test_result_v0_test_result_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_in_toto_attestation_predicates_test_result_v0_test_result_proto_goTypes,
		DependencyIndexes: file_in_toto_attestation_predicates_test_result_v0_test_result_proto_depIdxs,
		MessageInfos:      file_in_toto_attestation_predicates_test_result_v0_test_result_proto_msgTypes,
	}.Build()
	File_in_toto_attestation_predicates_test_result_v0_test_result_proto = out.File
	file_in_toto_attestation_predicates_test_result_v0_test_result_proto_rawDesc = nil
	file_in_toto_attestation_predicates_test_result_v0_test_result_proto_goTypes = nil
	file_in_toto_attestation_predicates_test_result_v0_test_result_proto_depIdxs = nil
}
