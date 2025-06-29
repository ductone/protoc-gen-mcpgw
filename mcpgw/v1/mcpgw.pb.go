// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.36.6
// 	protoc        (unknown)
// source: mcpgw/v1/mcpgw.proto

package v1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	descriptorpb "google.golang.org/protobuf/types/descriptorpb"
	_ "google.golang.org/protobuf/types/gofeaturespb"
	reflect "reflect"
	unsafe "unsafe"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type MessageOptions struct {
	state         protoimpl.MessageState `protogen:"opaque.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *MessageOptions) Reset() {
	*x = MessageOptions{}
	mi := &file_mcpgw_v1_mcpgw_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MessageOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MessageOptions) ProtoMessage() {}

func (x *MessageOptions) ProtoReflect() protoreflect.Message {
	mi := &file_mcpgw_v1_mcpgw_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

type MessageOptions_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

}

func (b0 MessageOptions_builder) Build() *MessageOptions {
	m0 := &MessageOptions{}
	b, x := &b0, m0
	_, _ = b, x
	return m0
}

type FieldOptions struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Description *string                `protobuf:"bytes,1,opt,name=description"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *FieldOptions) Reset() {
	*x = FieldOptions{}
	mi := &file_mcpgw_v1_mcpgw_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FieldOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FieldOptions) ProtoMessage() {}

func (x *FieldOptions) ProtoReflect() protoreflect.Message {
	mi := &file_mcpgw_v1_mcpgw_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *FieldOptions) GetDescription() string {
	if x != nil {
		if x.xxx_hidden_Description != nil {
			return *x.xxx_hidden_Description
		}
		return ""
	}
	return ""
}

func (x *FieldOptions) SetDescription(v string) {
	x.xxx_hidden_Description = &v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *FieldOptions) HasDescription() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *FieldOptions) ClearDescription() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Description = nil
}

type FieldOptions_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Description *string
}

func (b0 FieldOptions_builder) Build() *FieldOptions {
	m0 := &FieldOptions{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Description != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Description = b.Description
	}
	return m0
}

type MethodOptions struct {
	state                      protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Title           *string                `protobuf:"bytes,1,opt,name=title"`
	xxx_hidden_Description     *string                `protobuf:"bytes,2,opt,name=description"`
	xxx_hidden_ReadOnlyHint    bool                   `protobuf:"varint,3,opt,name=read_only_hint,json=readOnlyHint"`
	xxx_hidden_DestructiveHint bool                   `protobuf:"varint,4,opt,name=destructive_hint,json=destructiveHint"`
	xxx_hidden_IdempotentHint  bool                   `protobuf:"varint,5,opt,name=idempotent_hint,json=idempotentHint"`
	xxx_hidden_OpenWorldHint   bool                   `protobuf:"varint,6,opt,name=open_world_hint,json=openWorldHint"`
	XXX_raceDetectHookData     protoimpl.RaceDetectHookData
	XXX_presence               [1]uint32
	unknownFields              protoimpl.UnknownFields
	sizeCache                  protoimpl.SizeCache
}

func (x *MethodOptions) Reset() {
	*x = MethodOptions{}
	mi := &file_mcpgw_v1_mcpgw_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *MethodOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MethodOptions) ProtoMessage() {}

func (x *MethodOptions) ProtoReflect() protoreflect.Message {
	mi := &file_mcpgw_v1_mcpgw_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *MethodOptions) GetTitle() string {
	if x != nil {
		if x.xxx_hidden_Title != nil {
			return *x.xxx_hidden_Title
		}
		return ""
	}
	return ""
}

func (x *MethodOptions) GetDescription() string {
	if x != nil {
		if x.xxx_hidden_Description != nil {
			return *x.xxx_hidden_Description
		}
		return ""
	}
	return ""
}

func (x *MethodOptions) GetReadOnlyHint() bool {
	if x != nil {
		return x.xxx_hidden_ReadOnlyHint
	}
	return false
}

func (x *MethodOptions) GetDestructiveHint() bool {
	if x != nil {
		return x.xxx_hidden_DestructiveHint
	}
	return false
}

func (x *MethodOptions) GetIdempotentHint() bool {
	if x != nil {
		return x.xxx_hidden_IdempotentHint
	}
	return false
}

func (x *MethodOptions) GetOpenWorldHint() bool {
	if x != nil {
		return x.xxx_hidden_OpenWorldHint
	}
	return false
}

func (x *MethodOptions) SetTitle(v string) {
	x.xxx_hidden_Title = &v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 6)
}

func (x *MethodOptions) SetDescription(v string) {
	x.xxx_hidden_Description = &v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 1, 6)
}

func (x *MethodOptions) SetReadOnlyHint(v bool) {
	x.xxx_hidden_ReadOnlyHint = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 2, 6)
}

func (x *MethodOptions) SetDestructiveHint(v bool) {
	x.xxx_hidden_DestructiveHint = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 3, 6)
}

func (x *MethodOptions) SetIdempotentHint(v bool) {
	x.xxx_hidden_IdempotentHint = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 4, 6)
}

func (x *MethodOptions) SetOpenWorldHint(v bool) {
	x.xxx_hidden_OpenWorldHint = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 5, 6)
}

func (x *MethodOptions) HasTitle() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *MethodOptions) HasDescription() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 1)
}

func (x *MethodOptions) HasReadOnlyHint() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 2)
}

func (x *MethodOptions) HasDestructiveHint() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 3)
}

func (x *MethodOptions) HasIdempotentHint() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 4)
}

func (x *MethodOptions) HasOpenWorldHint() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 5)
}

func (x *MethodOptions) ClearTitle() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Title = nil
}

func (x *MethodOptions) ClearDescription() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 1)
	x.xxx_hidden_Description = nil
}

func (x *MethodOptions) ClearReadOnlyHint() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 2)
	x.xxx_hidden_ReadOnlyHint = false
}

func (x *MethodOptions) ClearDestructiveHint() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 3)
	x.xxx_hidden_DestructiveHint = false
}

func (x *MethodOptions) ClearIdempotentHint() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 4)
	x.xxx_hidden_IdempotentHint = false
}

func (x *MethodOptions) ClearOpenWorldHint() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 5)
	x.xxx_hidden_OpenWorldHint = false
}

type MethodOptions_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Title           *string
	Description     *string
	ReadOnlyHint    *bool
	DestructiveHint *bool
	IdempotentHint  *bool
	OpenWorldHint   *bool
}

func (b0 MethodOptions_builder) Build() *MethodOptions {
	m0 := &MethodOptions{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Title != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 6)
		x.xxx_hidden_Title = b.Title
	}
	if b.Description != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 1, 6)
		x.xxx_hidden_Description = b.Description
	}
	if b.ReadOnlyHint != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 2, 6)
		x.xxx_hidden_ReadOnlyHint = *b.ReadOnlyHint
	}
	if b.DestructiveHint != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 3, 6)
		x.xxx_hidden_DestructiveHint = *b.DestructiveHint
	}
	if b.IdempotentHint != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 4, 6)
		x.xxx_hidden_IdempotentHint = *b.IdempotentHint
	}
	if b.OpenWorldHint != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 5, 6)
		x.xxx_hidden_OpenWorldHint = *b.OpenWorldHint
	}
	return m0
}

type ServiceOptions struct {
	state                  protoimpl.MessageState `protogen:"opaque.v1"`
	xxx_hidden_Enabled     bool                   `protobuf:"varint,1,opt,name=enabled"`
	XXX_raceDetectHookData protoimpl.RaceDetectHookData
	XXX_presence           [1]uint32
	unknownFields          protoimpl.UnknownFields
	sizeCache              protoimpl.SizeCache
}

func (x *ServiceOptions) Reset() {
	*x = ServiceOptions{}
	mi := &file_mcpgw_v1_mcpgw_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ServiceOptions) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ServiceOptions) ProtoMessage() {}

func (x *ServiceOptions) ProtoReflect() protoreflect.Message {
	mi := &file_mcpgw_v1_mcpgw_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (x *ServiceOptions) GetEnabled() bool {
	if x != nil {
		return x.xxx_hidden_Enabled
	}
	return false
}

func (x *ServiceOptions) SetEnabled(v bool) {
	x.xxx_hidden_Enabled = v
	protoimpl.X.SetPresent(&(x.XXX_presence[0]), 0, 1)
}

func (x *ServiceOptions) HasEnabled() bool {
	if x == nil {
		return false
	}
	return protoimpl.X.Present(&(x.XXX_presence[0]), 0)
}

func (x *ServiceOptions) ClearEnabled() {
	protoimpl.X.ClearPresent(&(x.XXX_presence[0]), 0)
	x.xxx_hidden_Enabled = false
}

type ServiceOptions_builder struct {
	_ [0]func() // Prevents comparability and use of unkeyed literals for the builder.

	Enabled *bool
}

func (b0 ServiceOptions_builder) Build() *ServiceOptions {
	m0 := &ServiceOptions{}
	b, x := &b0, m0
	_, _ = b, x
	if b.Enabled != nil {
		protoimpl.X.SetPresentNonAtomic(&(x.XXX_presence[0]), 0, 1)
		x.xxx_hidden_Enabled = *b.Enabled
	}
	return m0
}

var file_mcpgw_v1_mcpgw_proto_extTypes = []protoimpl.ExtensionInfo{
	{
		ExtendedType:  (*descriptorpb.ServiceOptions)(nil),
		ExtensionType: (*ServiceOptions)(nil),
		Field:         8650,
		Name:          "mcpgw.v1.service",
		Tag:           "bytes,8650,opt,name=service",
		Filename:      "mcpgw/v1/mcpgw.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MethodOptions)(nil),
		ExtensionType: (*MethodOptions)(nil),
		Field:         8651,
		Name:          "mcpgw.v1.method",
		Tag:           "bytes,8651,opt,name=method",
		Filename:      "mcpgw/v1/mcpgw.proto",
	},
	{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*FieldOptions)(nil),
		Field:         8652,
		Name:          "mcpgw.v1.field",
		Tag:           "bytes,8652,opt,name=field",
		Filename:      "mcpgw/v1/mcpgw.proto",
	},
	{
		ExtendedType:  (*descriptorpb.MessageOptions)(nil),
		ExtensionType: (*MessageOptions)(nil),
		Field:         8653,
		Name:          "mcpgw.v1.message",
		Tag:           "bytes,8653,opt,name=message",
		Filename:      "mcpgw/v1/mcpgw.proto",
	},
}

// Extension fields to descriptorpb.ServiceOptions.
var (
	// optional mcpgw.v1.ServiceOptions service = 8650;
	E_Service = &file_mcpgw_v1_mcpgw_proto_extTypes[0]
)

// Extension fields to descriptorpb.MethodOptions.
var (
	// optional mcpgw.v1.MethodOptions method = 8651;
	E_Method = &file_mcpgw_v1_mcpgw_proto_extTypes[1]
)

// Extension fields to descriptorpb.FieldOptions.
var (
	// optional mcpgw.v1.FieldOptions field = 8652;
	E_Field = &file_mcpgw_v1_mcpgw_proto_extTypes[2]
)

// Extension fields to descriptorpb.MessageOptions.
var (
	// optional mcpgw.v1.MessageOptions message = 8653;
	E_Message = &file_mcpgw_v1_mcpgw_proto_extTypes[3]
)

var File_mcpgw_v1_mcpgw_proto protoreflect.FileDescriptor

const file_mcpgw_v1_mcpgw_proto_rawDesc = "" +
	"\n" +
	"\x14mcpgw/v1/mcpgw.proto\x12\bmcpgw.v1\x1a google/protobuf/descriptor.proto\x1a!google/protobuf/go_features.proto\"\x10\n" +
	"\x0eMessageOptions\"0\n" +
	"\fFieldOptions\x12 \n" +
	"\vdescription\x18\x01 \x01(\tR\vdescription\"\xe9\x01\n" +
	"\rMethodOptions\x12\x14\n" +
	"\x05title\x18\x01 \x01(\tR\x05title\x12 \n" +
	"\vdescription\x18\x02 \x01(\tR\vdescription\x12$\n" +
	"\x0eread_only_hint\x18\x03 \x01(\bR\freadOnlyHint\x12)\n" +
	"\x10destructive_hint\x18\x04 \x01(\bR\x0fdestructiveHint\x12'\n" +
	"\x0fidempotent_hint\x18\x05 \x01(\bR\x0eidempotentHint\x12&\n" +
	"\x0fopen_world_hint\x18\x06 \x01(\bR\ropenWorldHint\"*\n" +
	"\x0eServiceOptions\x12\x18\n" +
	"\aenabled\x18\x01 \x01(\bR\aenabled:T\n" +
	"\aservice\x12\x1f.google.protobuf.ServiceOptions\x18\xcaC \x01(\v2\x18.mcpgw.v1.ServiceOptionsR\aservice:P\n" +
	"\x06method\x12\x1e.google.protobuf.MethodOptions\x18\xcbC \x01(\v2\x17.mcpgw.v1.MethodOptionsR\x06method:L\n" +
	"\x05field\x12\x1d.google.protobuf.FieldOptions\x18\xccC \x01(\v2\x16.mcpgw.v1.FieldOptionsR\x05field:T\n" +
	"\amessage\x12\x1f.google.protobuf.MessageOptions\x18\xcdC \x01(\v2\x18.mcpgw.v1.MessageOptionsR\amessageB\x91\x01\n" +
	"\fcom.mcpgw.v1B\n" +
	"McpgwProtoP\x01Z,github.com/ductone/protoc-gen-mcpgw/mcpgw/v1\xa2\x02\x03MXX\xaa\x02\bMcpgw.V1\xca\x02\bMcpgw\\V1\xe2\x02\x14Mcpgw\\V1\\GPBMetadata\xea\x02\tMcpgw::V1\x92\x03\x05\xd2>\x02\x10\x03b\beditionsp\xe8\a"

var file_mcpgw_v1_mcpgw_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_mcpgw_v1_mcpgw_proto_goTypes = []any{
	(*MessageOptions)(nil),              // 0: mcpgw.v1.MessageOptions
	(*FieldOptions)(nil),                // 1: mcpgw.v1.FieldOptions
	(*MethodOptions)(nil),               // 2: mcpgw.v1.MethodOptions
	(*ServiceOptions)(nil),              // 3: mcpgw.v1.ServiceOptions
	(*descriptorpb.ServiceOptions)(nil), // 4: google.protobuf.ServiceOptions
	(*descriptorpb.MethodOptions)(nil),  // 5: google.protobuf.MethodOptions
	(*descriptorpb.FieldOptions)(nil),   // 6: google.protobuf.FieldOptions
	(*descriptorpb.MessageOptions)(nil), // 7: google.protobuf.MessageOptions
}
var file_mcpgw_v1_mcpgw_proto_depIdxs = []int32{
	4, // 0: mcpgw.v1.service:extendee -> google.protobuf.ServiceOptions
	5, // 1: mcpgw.v1.method:extendee -> google.protobuf.MethodOptions
	6, // 2: mcpgw.v1.field:extendee -> google.protobuf.FieldOptions
	7, // 3: mcpgw.v1.message:extendee -> google.protobuf.MessageOptions
	3, // 4: mcpgw.v1.service:type_name -> mcpgw.v1.ServiceOptions
	2, // 5: mcpgw.v1.method:type_name -> mcpgw.v1.MethodOptions
	1, // 6: mcpgw.v1.field:type_name -> mcpgw.v1.FieldOptions
	0, // 7: mcpgw.v1.message:type_name -> mcpgw.v1.MessageOptions
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	4, // [4:8] is the sub-list for extension type_name
	0, // [0:4] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_mcpgw_v1_mcpgw_proto_init() }
func file_mcpgw_v1_mcpgw_proto_init() {
	if File_mcpgw_v1_mcpgw_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_mcpgw_v1_mcpgw_proto_rawDesc), len(file_mcpgw_v1_mcpgw_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 4,
			NumServices:   0,
		},
		GoTypes:           file_mcpgw_v1_mcpgw_proto_goTypes,
		DependencyIndexes: file_mcpgw_v1_mcpgw_proto_depIdxs,
		MessageInfos:      file_mcpgw_v1_mcpgw_proto_msgTypes,
		ExtensionInfos:    file_mcpgw_v1_mcpgw_proto_extTypes,
	}.Build()
	File_mcpgw_v1_mcpgw_proto = out.File
	file_mcpgw_v1_mcpgw_proto_goTypes = nil
	file_mcpgw_v1_mcpgw_proto_depIdxs = nil
}
