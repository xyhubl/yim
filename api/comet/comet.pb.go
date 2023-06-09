// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.1
// 	protoc        v3.21.5
// source: comet/comet.proto

package comet

import (
	protocol "github.com/xyhubl/yim/api/protocol"
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

type PushMsgReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Keys    []string        `protobuf:"bytes,1,rep,name=Keys,proto3" json:"Keys,omitempty"`
	ProtoOp int32           `protobuf:"varint,2,opt,name=ProtoOp,proto3" json:"ProtoOp,omitempty"`
	Proto   *protocol.Proto `protobuf:"bytes,3,opt,name=proto,proto3" json:"proto,omitempty"`
}

func (x *PushMsgReq) Reset() {
	*x = PushMsgReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgReq) ProtoMessage() {}

func (x *PushMsgReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgReq.ProtoReflect.Descriptor instead.
func (*PushMsgReq) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{0}
}

func (x *PushMsgReq) GetKeys() []string {
	if x != nil {
		return x.Keys
	}
	return nil
}

func (x *PushMsgReq) GetProtoOp() int32 {
	if x != nil {
		return x.ProtoOp
	}
	return 0
}

func (x *PushMsgReq) GetProto() *protocol.Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

type PushMsgReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *PushMsgReply) Reset() {
	*x = PushMsgReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PushMsgReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PushMsgReply) ProtoMessage() {}

func (x *PushMsgReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PushMsgReply.ProtoReflect.Descriptor instead.
func (*PushMsgReply) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{1}
}

type BroadcastRoomReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	RoomID string          `protobuf:"bytes,1,opt,name=roomID,proto3" json:"roomID,omitempty"`
	Proto  *protocol.Proto `protobuf:"bytes,2,opt,name=proto,proto3" json:"proto,omitempty"`
}

func (x *BroadcastRoomReq) Reset() {
	*x = BroadcastRoomReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastRoomReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastRoomReq) ProtoMessage() {}

func (x *BroadcastRoomReq) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastRoomReq.ProtoReflect.Descriptor instead.
func (*BroadcastRoomReq) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{2}
}

func (x *BroadcastRoomReq) GetRoomID() string {
	if x != nil {
		return x.RoomID
	}
	return ""
}

func (x *BroadcastRoomReq) GetProto() *protocol.Proto {
	if x != nil {
		return x.Proto
	}
	return nil
}

type BroadcastRoomReply struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *BroadcastRoomReply) Reset() {
	*x = BroadcastRoomReply{}
	if protoimpl.UnsafeEnabled {
		mi := &file_comet_comet_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BroadcastRoomReply) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BroadcastRoomReply) ProtoMessage() {}

func (x *BroadcastRoomReply) ProtoReflect() protoreflect.Message {
	mi := &file_comet_comet_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BroadcastRoomReply.ProtoReflect.Descriptor instead.
func (*BroadcastRoomReply) Descriptor() ([]byte, []int) {
	return file_comet_comet_proto_rawDescGZIP(), []int{3}
}

var File_comet_comet_proto protoreflect.FileDescriptor

var file_comet_comet_proto_rawDesc = []byte{
	0x0a, 0x11, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x79, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x1a, 0x17,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x65, 0x0a, 0x0a, 0x50, 0x75, 0x73, 0x68, 0x4d,
	0x73, 0x67, 0x52, 0x65, 0x71, 0x12, 0x12, 0x0a, 0x04, 0x4b, 0x65, 0x79, 0x73, 0x18, 0x01, 0x20,
	0x03, 0x28, 0x09, 0x52, 0x04, 0x4b, 0x65, 0x79, 0x73, 0x12, 0x18, 0x0a, 0x07, 0x50, 0x72, 0x6f,
	0x74, 0x6f, 0x4f, 0x70, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x50, 0x72, 0x6f, 0x74,
	0x6f, 0x4f, 0x70, 0x12, 0x29, 0x0a, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x18, 0x03, 0x20, 0x01,
	0x28, 0x0b, 0x32, 0x13, 0x2e, 0x79, 0x69, 0x6d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f,
	0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x0e,
	0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x22, 0x55,
	0x0a, 0x10, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x52,
	0x65, 0x71, 0x12, 0x16, 0x0a, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x06, 0x72, 0x6f, 0x6f, 0x6d, 0x49, 0x44, 0x12, 0x29, 0x0a, 0x05, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x79, 0x69, 0x6d, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x63, 0x6f, 0x6c, 0x2e, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x52, 0x05,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x14, 0x0a, 0x12, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61,
	0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x32, 0x8f, 0x01, 0x0a, 0x05,
	0x43, 0x6f, 0x6d, 0x65, 0x74, 0x12, 0x39, 0x0a, 0x07, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67,
	0x12, 0x15, 0x2e, 0x79, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x50, 0x75, 0x73,
	0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x71, 0x1a, 0x17, 0x2e, 0x79, 0x69, 0x6d, 0x2e, 0x63, 0x6f,
	0x6d, 0x65, 0x74, 0x2e, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x73, 0x67, 0x52, 0x65, 0x70, 0x6c, 0x79,
	0x12, 0x4b, 0x0a, 0x0d, 0x42, 0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f, 0x6f,
	0x6d, 0x12, 0x1b, 0x2e, 0x79, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72,
	0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x71, 0x1a, 0x1d,
	0x2e, 0x79, 0x69, 0x6d, 0x2e, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x2e, 0x42, 0x72, 0x6f, 0x61, 0x64,
	0x63, 0x61, 0x73, 0x74, 0x52, 0x6f, 0x6f, 0x6d, 0x52, 0x65, 0x70, 0x6c, 0x79, 0x42, 0x27, 0x5a,
	0x25, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x79, 0x68, 0x75,
	0x62, 0x6c, 0x2f, 0x79, 0x69, 0x6d, 0x2f, 0x61, 0x70, 0x69, 0x2f, 0x63, 0x6f, 0x6d, 0x65, 0x74,
	0x3b, 0x63, 0x6f, 0x6d, 0x65, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_comet_comet_proto_rawDescOnce sync.Once
	file_comet_comet_proto_rawDescData = file_comet_comet_proto_rawDesc
)

func file_comet_comet_proto_rawDescGZIP() []byte {
	file_comet_comet_proto_rawDescOnce.Do(func() {
		file_comet_comet_proto_rawDescData = protoimpl.X.CompressGZIP(file_comet_comet_proto_rawDescData)
	})
	return file_comet_comet_proto_rawDescData
}

var file_comet_comet_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_comet_comet_proto_goTypes = []interface{}{
	(*PushMsgReq)(nil),         // 0: yim.comet.PushMsgReq
	(*PushMsgReply)(nil),       // 1: yim.comet.PushMsgReply
	(*BroadcastRoomReq)(nil),   // 2: yim.comet.BroadcastRoomReq
	(*BroadcastRoomReply)(nil), // 3: yim.comet.BroadcastRoomReply
	(*protocol.Proto)(nil),     // 4: yim.protocol.Proto
}
var file_comet_comet_proto_depIdxs = []int32{
	4, // 0: yim.comet.PushMsgReq.proto:type_name -> yim.protocol.Proto
	4, // 1: yim.comet.BroadcastRoomReq.proto:type_name -> yim.protocol.Proto
	0, // 2: yim.comet.Comet.PushMsg:input_type -> yim.comet.PushMsgReq
	2, // 3: yim.comet.Comet.BroadcastRoom:input_type -> yim.comet.BroadcastRoomReq
	1, // 4: yim.comet.Comet.PushMsg:output_type -> yim.comet.PushMsgReply
	3, // 5: yim.comet.Comet.BroadcastRoom:output_type -> yim.comet.BroadcastRoomReply
	4, // [4:6] is the sub-list for method output_type
	2, // [2:4] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_comet_comet_proto_init() }
func file_comet_comet_proto_init() {
	if File_comet_comet_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_comet_comet_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgReq); i {
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
		file_comet_comet_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PushMsgReply); i {
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
		file_comet_comet_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastRoomReq); i {
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
		file_comet_comet_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BroadcastRoomReply); i {
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
			RawDescriptor: file_comet_comet_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_comet_comet_proto_goTypes,
		DependencyIndexes: file_comet_comet_proto_depIdxs,
		MessageInfos:      file_comet_comet_proto_msgTypes,
	}.Build()
	File_comet_comet_proto = out.File
	file_comet_comet_proto_rawDesc = nil
	file_comet_comet_proto_goTypes = nil
	file_comet_comet_proto_depIdxs = nil
}
