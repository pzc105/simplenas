// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.29.0
// 	protoc        v4.24.3
// source: category.proto

package prpc

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

type CategoryItem_Type int32

const (
	CategoryItem_Unknown   CategoryItem_Type = 0
	CategoryItem_Home      CategoryItem_Type = 1
	CategoryItem_Directory CategoryItem_Type = 2
	CategoryItem_Video     CategoryItem_Type = 3
	CategoryItem_Other     CategoryItem_Type = 4
	CategoryItem_Audio     CategoryItem_Type = 5
	CategoryItem_MagnetUri CategoryItem_Type = 6
)

// Enum value maps for CategoryItem_Type.
var (
	CategoryItem_Type_name = map[int32]string{
		0: "Unknown",
		1: "Home",
		2: "Directory",
		3: "Video",
		4: "Other",
		5: "Audio",
		6: "MagnetUri",
	}
	CategoryItem_Type_value = map[string]int32{
		"Unknown":   0,
		"Home":      1,
		"Directory": 2,
		"Video":     3,
		"Other":     4,
		"Audio":     5,
		"MagnetUri": 6,
	}
)

func (x CategoryItem_Type) Enum() *CategoryItem_Type {
	p := new(CategoryItem_Type)
	*p = x
	return p
}

func (x CategoryItem_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (CategoryItem_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_category_proto_enumTypes[0].Descriptor()
}

func (CategoryItem_Type) Type() protoreflect.EnumType {
	return &file_category_proto_enumTypes[0]
}

func (x CategoryItem_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use CategoryItem_Type.Descriptor instead.
func (CategoryItem_Type) EnumDescriptor() ([]byte, []int) {
	return file_category_proto_rawDescGZIP(), []int{0, 0}
}

type CategoryItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Id           int64             `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	TypeId       CategoryItem_Type `protobuf:"varint,2,opt,name=type_id,json=typeId,proto3,enum=prpc.CategoryItem_Type" json:"type_id,omitempty"`
	Creator      int64             `protobuf:"varint,3,opt,name=creator,proto3" json:"creator,omitempty"`
	Name         string            `protobuf:"bytes,4,opt,name=name,proto3" json:"name,omitempty"`
	ResourcePath string            `protobuf:"bytes,5,opt,name=resource_path,json=resourcePath,proto3" json:"resource_path,omitempty"`
	PosterPath   string            `protobuf:"bytes,6,opt,name=poster_path,json=posterPath,proto3" json:"poster_path,omitempty"`
	Introduce    string            `protobuf:"bytes,7,opt,name=introduce,proto3" json:"introduce,omitempty"`
	Other        string            `protobuf:"bytes,8,opt,name=other,proto3" json:"other,omitempty"`
	ParentId     int64             `protobuf:"varint,9,opt,name=parent_id,json=parentId,proto3" json:"parent_id,omitempty"`
	SubItemIds   []int64           `protobuf:"varint,10,rep,packed,name=sub_item_ids,json=subItemIds,proto3" json:"sub_item_ids,omitempty"`
}

func (x *CategoryItem) Reset() {
	*x = CategoryItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_category_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CategoryItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CategoryItem) ProtoMessage() {}

func (x *CategoryItem) ProtoReflect() protoreflect.Message {
	mi := &file_category_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CategoryItem.ProtoReflect.Descriptor instead.
func (*CategoryItem) Descriptor() ([]byte, []int) {
	return file_category_proto_rawDescGZIP(), []int{0}
}

func (x *CategoryItem) GetId() int64 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *CategoryItem) GetTypeId() CategoryItem_Type {
	if x != nil {
		return x.TypeId
	}
	return CategoryItem_Unknown
}

func (x *CategoryItem) GetCreator() int64 {
	if x != nil {
		return x.Creator
	}
	return 0
}

func (x *CategoryItem) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *CategoryItem) GetResourcePath() string {
	if x != nil {
		return x.ResourcePath
	}
	return ""
}

func (x *CategoryItem) GetPosterPath() string {
	if x != nil {
		return x.PosterPath
	}
	return ""
}

func (x *CategoryItem) GetIntroduce() string {
	if x != nil {
		return x.Introduce
	}
	return ""
}

func (x *CategoryItem) GetOther() string {
	if x != nil {
		return x.Other
	}
	return ""
}

func (x *CategoryItem) GetParentId() int64 {
	if x != nil {
		return x.ParentId
	}
	return 0
}

func (x *CategoryItem) GetSubItemIds() []int64 {
	if x != nil {
		return x.SubItemIds
	}
	return nil
}

type SharedItem struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ItemId  int64  `protobuf:"varint,1,opt,name=item_id,json=itemId,proto3" json:"item_id,omitempty"`
	ShareId string `protobuf:"bytes,2,opt,name=share_id,json=shareId,proto3" json:"share_id,omitempty"`
}

func (x *SharedItem) Reset() {
	*x = SharedItem{}
	if protoimpl.UnsafeEnabled {
		mi := &file_category_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *SharedItem) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*SharedItem) ProtoMessage() {}

func (x *SharedItem) ProtoReflect() protoreflect.Message {
	mi := &file_category_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use SharedItem.ProtoReflect.Descriptor instead.
func (*SharedItem) Descriptor() ([]byte, []int) {
	return file_category_proto_rawDescGZIP(), []int{1}
}

func (x *SharedItem) GetItemId() int64 {
	if x != nil {
		return x.ItemId
	}
	return 0
}

func (x *SharedItem) GetShareId() string {
	if x != nil {
		return x.ShareId
	}
	return ""
}

var File_category_proto protoreflect.FileDescriptor

var file_category_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x63, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x04, 0x70, 0x72, 0x70, 0x63, 0x22, 0x95, 0x03, 0x0a, 0x0c, 0x43, 0x61, 0x74, 0x65, 0x67,
	0x6f, 0x72, 0x79, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x03, 0x52, 0x02, 0x69, 0x64, 0x12, 0x30, 0x0a, 0x07, 0x74, 0x79, 0x70, 0x65, 0x5f,
	0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x17, 0x2e, 0x70, 0x72, 0x70, 0x63, 0x2e,
	0x43, 0x61, 0x74, 0x65, 0x67, 0x6f, 0x72, 0x79, 0x49, 0x74, 0x65, 0x6d, 0x2e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x06, 0x74, 0x79, 0x70, 0x65, 0x49, 0x64, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x72, 0x65,
	0x61, 0x74, 0x6f, 0x72, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x07, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x0d, 0x72, 0x65, 0x73, 0x6f, 0x75,
	0x72, 0x63, 0x65, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x05, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c,
	0x72, 0x65, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1f, 0x0a, 0x0b,
	0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x5f, 0x70, 0x61, 0x74, 0x68, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0a, 0x70, 0x6f, 0x73, 0x74, 0x65, 0x72, 0x50, 0x61, 0x74, 0x68, 0x12, 0x1c, 0x0a,
	0x09, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x09, 0x69, 0x6e, 0x74, 0x72, 0x6f, 0x64, 0x75, 0x63, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6f,
	0x74, 0x68, 0x65, 0x72, 0x18, 0x08, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x6f, 0x74, 0x68, 0x65,
	0x72, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x5f, 0x69, 0x64, 0x18, 0x09,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x08, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x49, 0x64, 0x12, 0x20,
	0x0a, 0x0c, 0x73, 0x75, 0x62, 0x5f, 0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x0a,
	0x20, 0x03, 0x28, 0x03, 0x52, 0x0a, 0x73, 0x75, 0x62, 0x49, 0x74, 0x65, 0x6d, 0x49, 0x64, 0x73,
	0x22, 0x5c, 0x0a, 0x04, 0x54, 0x79, 0x70, 0x65, 0x12, 0x0b, 0x0a, 0x07, 0x55, 0x6e, 0x6b, 0x6e,
	0x6f, 0x77, 0x6e, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x48, 0x6f, 0x6d, 0x65, 0x10, 0x01, 0x12,
	0x0d, 0x0a, 0x09, 0x44, 0x69, 0x72, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x79, 0x10, 0x02, 0x12, 0x09,
	0x0a, 0x05, 0x56, 0x69, 0x64, 0x65, 0x6f, 0x10, 0x03, 0x12, 0x09, 0x0a, 0x05, 0x4f, 0x74, 0x68,
	0x65, 0x72, 0x10, 0x04, 0x12, 0x09, 0x0a, 0x05, 0x41, 0x75, 0x64, 0x69, 0x6f, 0x10, 0x05, 0x12,
	0x0d, 0x0a, 0x09, 0x4d, 0x61, 0x67, 0x6e, 0x65, 0x74, 0x55, 0x72, 0x69, 0x10, 0x06, 0x22, 0x40,
	0x0a, 0x0a, 0x53, 0x68, 0x61, 0x72, 0x65, 0x64, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x17, 0x0a, 0x07,
	0x69, 0x74, 0x65, 0x6d, 0x5f, 0x69, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x06, 0x69,
	0x74, 0x65, 0x6d, 0x49, 0x64, 0x12, 0x19, 0x0a, 0x08, 0x73, 0x68, 0x61, 0x72, 0x65, 0x5f, 0x69,
	0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x68, 0x61, 0x72, 0x65, 0x49, 0x64,
	0x42, 0x08, 0x5a, 0x06, 0x2e, 0x2f, 0x70, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_category_proto_rawDescOnce sync.Once
	file_category_proto_rawDescData = file_category_proto_rawDesc
)

func file_category_proto_rawDescGZIP() []byte {
	file_category_proto_rawDescOnce.Do(func() {
		file_category_proto_rawDescData = protoimpl.X.CompressGZIP(file_category_proto_rawDescData)
	})
	return file_category_proto_rawDescData
}

var file_category_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_category_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_category_proto_goTypes = []interface{}{
	(CategoryItem_Type)(0), // 0: prpc.CategoryItem.Type
	(*CategoryItem)(nil),   // 1: prpc.CategoryItem
	(*SharedItem)(nil),     // 2: prpc.SharedItem
}
var file_category_proto_depIdxs = []int32{
	0, // 0: prpc.CategoryItem.type_id:type_name -> prpc.CategoryItem.Type
	1, // [1:1] is the sub-list for method output_type
	1, // [1:1] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_category_proto_init() }
func file_category_proto_init() {
	if File_category_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_category_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CategoryItem); i {
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
		file_category_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*SharedItem); i {
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
			RawDescriptor: file_category_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_category_proto_goTypes,
		DependencyIndexes: file_category_proto_depIdxs,
		EnumInfos:         file_category_proto_enumTypes,
		MessageInfos:      file_category_proto_msgTypes,
	}.Build()
	File_category_proto = out.File
	file_category_proto_rawDesc = nil
	file_category_proto_goTypes = nil
	file_category_proto_depIdxs = nil
}
