// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        v4.24.4
// source: in_toto_attestation/v1/statement.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	structpb "google.golang.org/protobuf/types/known/structpb"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// Proto representation of the in-toto v1 Statement.
// https://github.com/in-toto/attestation/tree/main/spec/v1
// Validation of all fields is left to the users of this proto.
type Statement struct {
	state protoimpl.MessageState `protogen:"open.v1"`
	// Expected to always be "https://in-toto.io/Statement/v1"
	Type          string                `protobuf:"bytes,1,opt,name=type,json=_type,proto3" json:"type,omitempty"`
	Subject       []*ResourceDescriptor `protobuf:"bytes,2,rep,name=subject,proto3" json:"subject,omitempty"`
	PredicateType string                `protobuf:"bytes,3,opt,name=predicate_type,json=predicateType,proto3" json:"predicate_type,omitempty"`
	Predicate     *structpb.Struct      `protobuf:"bytes,4,opt,name=predicate,proto3" json:"predicate,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Statement) Reset() {
	*x = Statement{}
	mi := &file_in_toto_attestation_v1_statement_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Statement) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Statement) ProtoMessage() {}

func (x *Statement) ProtoReflect() protoreflect.Message {
	mi := &file_in_toto_attestation_v1_statement_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Statement.ProtoReflect.Descriptor instead.
func (*Statement) Descriptor() ([]byte, []int) {
	return file_in_toto_attestation_v1_statement_proto_rawDescGZIP(), []int{0}
}

func (x *Statement) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Statement) GetSubject() []*ResourceDescriptor {
	if x != nil {
		return x.Subject
	}
	return nil
}

func (x *Statement) GetPredicateType() string {
	if x != nil {
		return x.PredicateType
	}
	return ""
}

func (x *Statement) GetPredicate() *structpb.Struct {
	if x != nil {
		return x.Predicate
	}
	return nil
}

var File_in_toto_attestation_v1_statement_proto protoreflect.FileDescriptor

const file_in_toto_attestation_v1_statement_proto_rawDesc = "" +
	"\n" +
	"&in_toto_attestation/v1/statement.proto\x12\x16in_toto_attestation.v1\x1a\x1cgoogle/protobuf/struct.proto\x1a0in_toto_attestation/v1/resource_descriptor.proto\"\xc4\x01\n" +
	"\tStatement\x12\x13\n" +
	"\x04type\x18\x01 \x01(\tR\x05_type\x12D\n" +
	"\asubject\x18\x02 \x03(\v2*.in_toto_attestation.v1.ResourceDescriptorR\asubject\x12%\n" +
	"\x0epredicate_type\x18\x03 \x01(\tR\rpredicateType\x125\n" +
	"\tpredicate\x18\x04 \x01(\v2\x17.google.protobuf.StructR\tpredicateBG\n" +
	"\x1fio.github.intoto.attestation.v1Z$github.com/in-toto/attestation/go/v1b\x06proto3"

var (
	file_in_toto_attestation_v1_statement_proto_rawDescOnce sync.Once
	file_in_toto_attestation_v1_statement_proto_rawDescData []byte
)

func file_in_toto_attestation_v1_statement_proto_rawDescGZIP() []byte {
	file_in_toto_attestation_v1_statement_proto_rawDescOnce.Do(func() {
		file_in_toto_attestation_v1_statement_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_in_toto_attestation_v1_statement_proto_rawDesc), len(file_in_toto_attestation_v1_statement_proto_rawDesc)))
	})
	return file_in_toto_attestation_v1_statement_proto_rawDescData
}

var file_in_toto_attestation_v1_statement_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_in_toto_attestation_v1_statement_proto_goTypes = []any{
	(*Statement)(nil),          // 0: in_toto_attestation.v1.Statement
	(*ResourceDescriptor)(nil), // 1: in_toto_attestation.v1.ResourceDescriptor
	(*structpb.Struct)(nil),    // 2: google.protobuf.Struct
}
var file_in_toto_attestation_v1_statement_proto_depIdxs = []int32{
	1, // 0: in_toto_attestation.v1.Statement.subject:type_name -> in_toto_attestation.v1.ResourceDescriptor
	2, // 1: in_toto_attestation.v1.Statement.predicate:type_name -> google.protobuf.Struct
	2, // [2:2] is the sub-list for method output_type
	2, // [2:2] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_in_toto_attestation_v1_statement_proto_init() }
func file_in_toto_attestation_v1_statement_proto_init() {
	if File_in_toto_attestation_v1_statement_proto != nil {
		return
	}
	file_in_toto_attestation_v1_resource_descriptor_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_in_toto_attestation_v1_statement_proto_rawDesc), len(file_in_toto_attestation_v1_statement_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_in_toto_attestation_v1_statement_proto_goTypes,
		DependencyIndexes: file_in_toto_attestation_v1_statement_proto_depIdxs,
		MessageInfos:      file_in_toto_attestation_v1_statement_proto_msgTypes,
	}.Build()
	File_in_toto_attestation_v1_statement_proto = out.File
	file_in_toto_attestation_v1_statement_proto_goTypes = nil
	file_in_toto_attestation_v1_statement_proto_depIdxs = nil
}
