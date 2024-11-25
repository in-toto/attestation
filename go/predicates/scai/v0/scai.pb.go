// Protobuf definition for the SCAI Attribute Report predicate
// (predicateType = https://in-toto.io/attestation/scai/attribute-report/v0.2)
//
// Validation of all fields is left to the users of this proto.

// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.24.4
// source: in_toto_attestation/predicates/scai/v0/scai.proto

package v0

import (
	v1 "github.com/in-toto/attestation/go/v1"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type AttributeAssertion struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attribute  string                 `protobuf:"bytes,1,opt,name=attribute,proto3" json:"attribute,omitempty"` // required
	Target     *v1.ResourceDescriptor `protobuf:"bytes,2,opt,name=target,proto3" json:"target,omitempty"`
	Conditions *structpb.Struct       `protobuf:"bytes,3,opt,name=conditions,proto3" json:"conditions,omitempty"`
	Evidence   *v1.ResourceDescriptor `protobuf:"bytes,4,opt,name=evidence,proto3" json:"evidence,omitempty"`
}

func (x *AttributeAssertion) Reset() {
	*x = AttributeAssertion{}
	mi := &file_in_toto_attestation_predicates_scai_v0_scai_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AttributeAssertion) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttributeAssertion) ProtoMessage() {}

func (x *AttributeAssertion) ProtoReflect() protoreflect.Message {
	mi := &file_in_toto_attestation_predicates_scai_v0_scai_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttributeAssertion.ProtoReflect.Descriptor instead.
func (*AttributeAssertion) Descriptor() ([]byte, []int) {
	return file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescGZIP(), []int{0}
}

func (x *AttributeAssertion) GetAttribute() string {
	if x != nil {
		return x.Attribute
	}
	return ""
}

func (x *AttributeAssertion) GetTarget() *v1.ResourceDescriptor {
	if x != nil {
		return x.Target
	}
	return nil
}

func (x *AttributeAssertion) GetConditions() *structpb.Struct {
	if x != nil {
		return x.Conditions
	}
	return nil
}

func (x *AttributeAssertion) GetEvidence() *v1.ResourceDescriptor {
	if x != nil {
		return x.Evidence
	}
	return nil
}

type AttributeReport struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Attributes []*AttributeAssertion  `protobuf:"bytes,1,rep,name=attributes,proto3" json:"attributes,omitempty"`
	Producer   *v1.ResourceDescriptor `protobuf:"bytes,2,opt,name=producer,proto3" json:"producer,omitempty"`
}

func (x *AttributeReport) Reset() {
	*x = AttributeReport{}
	mi := &file_in_toto_attestation_predicates_scai_v0_scai_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *AttributeReport) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AttributeReport) ProtoMessage() {}

func (x *AttributeReport) ProtoReflect() protoreflect.Message {
	mi := &file_in_toto_attestation_predicates_scai_v0_scai_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AttributeReport.ProtoReflect.Descriptor instead.
func (*AttributeReport) Descriptor() ([]byte, []int) {
	return file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescGZIP(), []int{1}
}

func (x *AttributeReport) GetAttributes() []*AttributeAssertion {
	if x != nil {
		return x.Attributes
	}
	return nil
}

func (x *AttributeReport) GetProducer() *v1.ResourceDescriptor {
	if x != nil {
		return x.Producer
	}
	return nil
}

var File_in_toto_attestation_predicates_scai_v0_scai_proto protoreflect.FileDescriptor

var file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDesc = []byte{
	0x0a, 0x31, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73,
	0x2f, 0x73, 0x63, 0x61, 0x69, 0x2f, 0x76, 0x30, 0x2f, 0x73, 0x63, 0x61, 0x69, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x26, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74,
	0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61,
	0x74, 0x65, 0x73, 0x2e, 0x73, 0x63, 0x61, 0x69, 0x2e, 0x76, 0x30, 0x1a, 0x1c, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x73, 0x74, 0x72,
	0x75, 0x63, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x30, 0x69, 0x6e, 0x5f, 0x74, 0x6f,
	0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2f, 0x76,
	0x31, 0x2f, 0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x5f, 0x64, 0x65, 0x73, 0x63, 0x72,
	0x69, 0x70, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xf7, 0x01, 0x0a, 0x12,
	0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69,
	0x6f, 0x6e, 0x12, 0x1c, 0x0a, 0x09, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x61, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x12, 0x42, 0x0a, 0x06, 0x74, 0x61, 0x72, 0x67, 0x65, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b,
	0x32, 0x2a, 0x2e, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73,
	0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72,
	0x63, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x06, 0x74, 0x61,
	0x72, 0x67, 0x65, 0x74, 0x12, 0x37, 0x0a, 0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x73, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x17, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x75, 0x63,
	0x74, 0x52, 0x0a, 0x63, 0x6f, 0x6e, 0x64, 0x69, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x12, 0x46, 0x0a,
	0x08, 0x65, 0x76, 0x69, 0x64, 0x65, 0x6e, 0x63, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x2a, 0x2e, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74,
	0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31, 0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63,
	0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x6f, 0x72, 0x52, 0x08, 0x65, 0x76, 0x69,
	0x64, 0x65, 0x6e, 0x63, 0x65, 0x22, 0xb5, 0x01, 0x0a, 0x0f, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62,
	0x75, 0x74, 0x65, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x12, 0x5a, 0x0a, 0x0a, 0x61, 0x74, 0x74,
	0x72, 0x69, 0x62, 0x75, 0x74, 0x65, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x3a, 0x2e,
	0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74, 0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74,
	0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x2e, 0x73,
	0x63, 0x61, 0x69, 0x2e, 0x76, 0x30, 0x2e, 0x41, 0x74, 0x74, 0x72, 0x69, 0x62, 0x75, 0x74, 0x65,
	0x41, 0x73, 0x73, 0x65, 0x72, 0x74, 0x69, 0x6f, 0x6e, 0x52, 0x0a, 0x61, 0x74, 0x74, 0x72, 0x69,
	0x62, 0x75, 0x74, 0x65, 0x73, 0x12, 0x46, 0x0a, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x2a, 0x2e, 0x69, 0x6e, 0x5f, 0x74, 0x6f, 0x74,
	0x6f, 0x5f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x76, 0x31,
	0x2e, 0x52, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x44, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70,
	0x74, 0x6f, 0x72, 0x52, 0x08, 0x70, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x72, 0x42, 0x67, 0x0a,
	0x2f, 0x69, 0x6f, 0x2e, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x69, 0x6e, 0x74, 0x6f, 0x74,
	0x6f, 0x2e, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x2e, 0x70, 0x72,
	0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x2e, 0x73, 0x63, 0x61, 0x69, 0x2e, 0x76, 0x30,
	0x5a, 0x34, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x69, 0x6e, 0x2d,
	0x74, 0x6f, 0x74, 0x6f, 0x2f, 0x61, 0x74, 0x74, 0x65, 0x73, 0x74, 0x61, 0x74, 0x69, 0x6f, 0x6e,
	0x2f, 0x67, 0x6f, 0x2f, 0x70, 0x72, 0x65, 0x64, 0x69, 0x63, 0x61, 0x74, 0x65, 0x73, 0x2f, 0x73,
	0x63, 0x61, 0x69, 0x2f, 0x76, 0x30, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescOnce sync.Once
	file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescData = file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDesc
)

func file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescGZIP() []byte {
	file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescOnce.Do(func() {
		file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescData = protoimpl.X.CompressGZIP(file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescData)
	})
	return file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDescData
}

var file_in_toto_attestation_predicates_scai_v0_scai_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_in_toto_attestation_predicates_scai_v0_scai_proto_goTypes = []any{
	(*AttributeAssertion)(nil),    // 0: in_toto_attestation.predicates.scai.v0.AttributeAssertion
	(*AttributeReport)(nil),       // 1: in_toto_attestation.predicates.scai.v0.AttributeReport
	(*v1.ResourceDescriptor)(nil), // 2: in_toto_attestation.v1.ResourceDescriptor
	(*structpb.Struct)(nil),       // 3: google.protobuf.Struct
}
var file_in_toto_attestation_predicates_scai_v0_scai_proto_depIdxs = []int32{
	2, // 0: in_toto_attestation.predicates.scai.v0.AttributeAssertion.target:type_name -> in_toto_attestation.v1.ResourceDescriptor
	3, // 1: in_toto_attestation.predicates.scai.v0.AttributeAssertion.conditions:type_name -> google.protobuf.Struct
	2, // 2: in_toto_attestation.predicates.scai.v0.AttributeAssertion.evidence:type_name -> in_toto_attestation.v1.ResourceDescriptor
	0, // 3: in_toto_attestation.predicates.scai.v0.AttributeReport.attributes:type_name -> in_toto_attestation.predicates.scai.v0.AttributeAssertion
	2, // 4: in_toto_attestation.predicates.scai.v0.AttributeReport.producer:type_name -> in_toto_attestation.v1.ResourceDescriptor
	5, // [5:5] is the sub-list for method output_type
	5, // [5:5] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_in_toto_attestation_predicates_scai_v0_scai_proto_init() }
func file_in_toto_attestation_predicates_scai_v0_scai_proto_init() {
	if File_in_toto_attestation_predicates_scai_v0_scai_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_in_toto_attestation_predicates_scai_v0_scai_proto_goTypes,
		DependencyIndexes: file_in_toto_attestation_predicates_scai_v0_scai_proto_depIdxs,
		MessageInfos:      file_in_toto_attestation_predicates_scai_v0_scai_proto_msgTypes,
	}.Build()
	File_in_toto_attestation_predicates_scai_v0_scai_proto = out.File
	file_in_toto_attestation_predicates_scai_v0_scai_proto_rawDesc = nil
	file_in_toto_attestation_predicates_scai_v0_scai_proto_goTypes = nil
	file_in_toto_attestation_predicates_scai_v0_scai_proto_depIdxs = nil
}
