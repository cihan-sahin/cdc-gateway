package cdcv1

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
	unsafe "unsafe"
)

const (
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Policy struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Table         string                 `protobuf:"bytes,1,opt,name=table,proto3" json:"table,omitempty"`
	Mode          string                 `protobuf:"bytes,2,opt,name=mode,proto3" json:"mode,omitempty"`
	WindowMs      int32                  `protobuf:"varint,3,opt,name=window_ms,json=windowMs,proto3" json:"window_ms,omitempty"`
	MaxBatchSize  int32                  `protobuf:"varint,4,opt,name=max_batch_size,json=maxBatchSize,proto3" json:"max_batch_size,omitempty"`
	MergeStrategy string                 `protobuf:"bytes,5,opt,name=merge_strategy,json=mergeStrategy,proto3" json:"merge_strategy,omitempty"`
	TargetTopic   string                 `protobuf:"bytes,6,opt,name=target_topic,json=targetTopic,proto3" json:"target_topic,omitempty"`
	Enabled       bool                   `protobuf:"varint,7,opt,name=enabled,proto3" json:"enabled,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *Policy) Reset() {
	*x = Policy{}
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Policy) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Policy) ProtoMessage() {}

func (x *Policy) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*Policy) Descriptor() ([]byte, []int) {
	return file_proto_cdcgateway_v1_policy_proto_rawDescGZIP(), []int{0}
}

func (x *Policy) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

func (x *Policy) GetMode() string {
	if x != nil {
		return x.Mode
	}
	return ""
}

func (x *Policy) GetWindowMs() int32 {
	if x != nil {
		return x.WindowMs
	}
	return 0
}

func (x *Policy) GetMaxBatchSize() int32 {
	if x != nil {
		return x.MaxBatchSize
	}
	return 0
}

func (x *Policy) GetMergeStrategy() string {
	if x != nil {
		return x.MergeStrategy
	}
	return ""
}

func (x *Policy) GetTargetTopic() string {
	if x != nil {
		return x.TargetTopic
	}
	return ""
}

func (x *Policy) GetEnabled() bool {
	if x != nil {
		return x.Enabled
	}
	return false
}

type ListPoliciesRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListPoliciesRequest) Reset() {
	*x = ListPoliciesRequest{}
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListPoliciesRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPoliciesRequest) ProtoMessage() {}

func (x *ListPoliciesRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ListPoliciesRequest.ProtoReflect.Descriptor instead.
func (*ListPoliciesRequest) Descriptor() ([]byte, []int) {
	return file_proto_cdcgateway_v1_policy_proto_rawDescGZIP(), []int{1}
}

type ListPoliciesResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Policies      []*Policy              `protobuf:"bytes,1,rep,name=policies,proto3" json:"policies,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *ListPoliciesResponse) Reset() {
	*x = ListPoliciesResponse{}
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ListPoliciesResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ListPoliciesResponse) ProtoMessage() {}

func (x *ListPoliciesResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*ListPoliciesResponse) Descriptor() ([]byte, []int) {
	return file_proto_cdcgateway_v1_policy_proto_rawDescGZIP(), []int{2}
}

func (x *ListPoliciesResponse) GetPolicies() []*Policy {
	if x != nil {
		return x.Policies
	}
	return nil
}

type UpsertPolicyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Policy        *Policy                `protobuf:"bytes,1,opt,name=policy,proto3" json:"policy,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpsertPolicyRequest) Reset() {
	*x = UpsertPolicyRequest{}
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[3]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpsertPolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertPolicyRequest) ProtoMessage() {}

func (x *UpsertPolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[3]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*UpsertPolicyRequest) Descriptor() ([]byte, []int) {
	return file_proto_cdcgateway_v1_policy_proto_rawDescGZIP(), []int{3}
}

func (x *UpsertPolicyRequest) GetPolicy() *Policy {
	if x != nil {
		return x.Policy
	}
	return nil
}

type UpsertPolicyResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *UpsertPolicyResponse) Reset() {
	*x = UpsertPolicyResponse{}
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[4]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *UpsertPolicyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpsertPolicyResponse) ProtoMessage() {}

func (x *UpsertPolicyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[4]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*UpsertPolicyResponse) Descriptor() ([]byte, []int) {
	return file_proto_cdcgateway_v1_policy_proto_rawDescGZIP(), []int{4}
}

type DeletePolicyRequest struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	Table         string                 `protobuf:"bytes,1,opt,name=table,proto3" json:"table,omitempty"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeletePolicyRequest) Reset() {
	*x = DeletePolicyRequest{}
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[5]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeletePolicyRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePolicyRequest) ProtoMessage() {}

func (x *DeletePolicyRequest) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[5]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DeletePolicyRequest) Descriptor() ([]byte, []int) {
	return file_proto_cdcgateway_v1_policy_proto_rawDescGZIP(), []int{5}
}

func (x *DeletePolicyRequest) GetTable() string {
	if x != nil {
		return x.Table
	}
	return ""
}

type DeletePolicyResponse struct {
	state         protoimpl.MessageState `protogen:"open.v1"`
	unknownFields protoimpl.UnknownFields
	sizeCache     protoimpl.SizeCache
}

func (x *DeletePolicyResponse) Reset() {
	*x = DeletePolicyResponse{}
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[6]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *DeletePolicyResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeletePolicyResponse) ProtoMessage() {}

func (x *DeletePolicyResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_cdcgateway_v1_policy_proto_msgTypes[6]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

func (*DeletePolicyResponse) Descriptor() ([]byte, []int) {
	return file_proto_cdcgateway_v1_policy_proto_rawDescGZIP(), []int{6}
}

var File_proto_cdcgateway_v1_policy_proto protoreflect.FileDescriptor

const file_proto_cdcgateway_v1_policy_proto_rawDesc = "" +
	"\n" +
	" proto/cdcgateway/v1/policy.proto\x12\rcdcgateway.v1\"\xd9\x01\n" +
	"\x06Policy\x12\x14\n" +
	"\x05table\x18\x01 \x01(\tR\x05table\x12\x12\n" +
	"\x04mode\x18\x02 \x01(\tR\x04mode\x12\x1b\n" +
	"\twindow_ms\x18\x03 \x01(\x05R\bwindowMs\x12$\n" +
	"\x0emax_batch_size\x18\x04 \x01(\x05R\fmaxBatchSize\x12%\n" +
	"\x0emerge_strategy\x18\x05 \x01(\tR\rmergeStrategy\x12!\n" +
	"\ftarget_topic\x18\x06 \x01(\tR\vtargetTopic\x12\x18\n" +
	"\aenabled\x18\a \x01(\bR\aenabled\"\x15\n" +
	"\x13ListPoliciesRequest\"I\n" +
	"\x14ListPoliciesResponse\x121\n" +
	"\bpolicies\x18\x01 \x03(\v2\x15.cdcgateway.v1.PolicyR\bpolicies\"D\n" +
	"\x13UpsertPolicyRequest\x12-\n" +
	"\x06policy\x18\x01 \x01(\v2\x15.cdcgateway.v1.PolicyR\x06policy\"\x16\n" +
	"\x14UpsertPolicyResponse\"+\n" +
	"\x13DeletePolicyRequest\x12\x14\n" +
	"\x05table\x18\x01 \x01(\tR\x05table\"\x16\n" +
	"\x14DeletePolicyResponse2\x9a\x02\n" +
	"\rPolicyService\x12W\n" +
	"\fListPolicies\x12\".cdcgateway.v1.ListPoliciesRequest\x1a#.cdcgateway.v1.ListPoliciesResponse\x12W\n" +
	"\fUpsertPolicy\x12\".cdcgateway.v1.UpsertPolicyRequest\x1a#.cdcgateway.v1.UpsertPolicyResponse\x12W\n" +
	"\fDeletePolicy\x12\".cdcgateway.v1.DeletePolicyRequest\x1a#.cdcgateway.v1.DeletePolicyResponseB<Z:github.com/cihan-sahin/cdc-gateway/gen/cdcgateway/v1;cdcv1b\x06proto3"

var (
	file_proto_cdcgateway_v1_policy_proto_rawDescOnce sync.Once
	file_proto_cdcgateway_v1_policy_proto_rawDescData []byte
)

func file_proto_cdcgateway_v1_policy_proto_rawDescGZIP() []byte {
	file_proto_cdcgateway_v1_policy_proto_rawDescOnce.Do(func() {
		file_proto_cdcgateway_v1_policy_proto_rawDescData = protoimpl.X.CompressGZIP(unsafe.Slice(unsafe.StringData(file_proto_cdcgateway_v1_policy_proto_rawDesc), len(file_proto_cdcgateway_v1_policy_proto_rawDesc)))
	})
	return file_proto_cdcgateway_v1_policy_proto_rawDescData
}

var file_proto_cdcgateway_v1_policy_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_proto_cdcgateway_v1_policy_proto_goTypes = []any{
	(*Policy)(nil),               // 0: cdcgateway.v1.Policy
	(*ListPoliciesRequest)(nil),  // 1: cdcgateway.v1.ListPoliciesRequest
	(*ListPoliciesResponse)(nil), // 2: cdcgateway.v1.ListPoliciesResponse
	(*UpsertPolicyRequest)(nil),  // 3: cdcgateway.v1.UpsertPolicyRequest
	(*UpsertPolicyResponse)(nil), // 4: cdcgateway.v1.UpsertPolicyResponse
	(*DeletePolicyRequest)(nil),  // 5: cdcgateway.v1.DeletePolicyRequest
	(*DeletePolicyResponse)(nil), // 6: cdcgateway.v1.DeletePolicyResponse
}
var file_proto_cdcgateway_v1_policy_proto_depIdxs = []int32{
	0, // 0: cdcgateway.v1.ListPoliciesResponse.policies:type_name -> cdcgateway.v1.Policy
	0, // 1: cdcgateway.v1.UpsertPolicyRequest.policy:type_name -> cdcgateway.v1.Policy
	1, // 2: cdcgateway.v1.PolicyService.ListPolicies:input_type -> cdcgateway.v1.ListPoliciesRequest
	3, // 3: cdcgateway.v1.PolicyService.UpsertPolicy:input_type -> cdcgateway.v1.UpsertPolicyRequest
	5, // 4: cdcgateway.v1.PolicyService.DeletePolicy:input_type -> cdcgateway.v1.DeletePolicyRequest
	2, // 5: cdcgateway.v1.PolicyService.ListPolicies:output_type -> cdcgateway.v1.ListPoliciesResponse
	4, // 6: cdcgateway.v1.PolicyService.UpsertPolicy:output_type -> cdcgateway.v1.UpsertPolicyResponse
	6, // 7: cdcgateway.v1.PolicyService.DeletePolicy:output_type -> cdcgateway.v1.DeletePolicyResponse
	5, // [5:8] is the sub-list for method output_type
	2, // [2:5] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_proto_cdcgateway_v1_policy_proto_init() }
func file_proto_cdcgateway_v1_policy_proto_init() {
	if File_proto_cdcgateway_v1_policy_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: unsafe.Slice(unsafe.StringData(file_proto_cdcgateway_v1_policy_proto_rawDesc), len(file_proto_cdcgateway_v1_policy_proto_rawDesc)),
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_proto_cdcgateway_v1_policy_proto_goTypes,
		DependencyIndexes: file_proto_cdcgateway_v1_policy_proto_depIdxs,
		MessageInfos:      file_proto_cdcgateway_v1_policy_proto_msgTypes,
	}.Build()
	File_proto_cdcgateway_v1_policy_proto = out.File
	file_proto_cdcgateway_v1_policy_proto_goTypes = nil
	file_proto_cdcgateway_v1_policy_proto_depIdxs = nil
}
