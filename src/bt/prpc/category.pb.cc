// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: category.proto

#include "category.pb.h"

#include <algorithm>
#include "google/protobuf/io/coded_stream.h"
#include "google/protobuf/extension_set.h"
#include "google/protobuf/wire_format_lite.h"
#include "google/protobuf/descriptor.h"
#include "google/protobuf/generated_message_reflection.h"
#include "google/protobuf/reflection_ops.h"
#include "google/protobuf/wire_format.h"
// @@protoc_insertion_point(includes)

// Must be included last.
#include "google/protobuf/port_def.inc"
PROTOBUF_PRAGMA_INIT_SEG
namespace _pb = ::PROTOBUF_NAMESPACE_ID;
namespace _pbi = ::PROTOBUF_NAMESPACE_ID::internal;
namespace prpc {
PROTOBUF_CONSTEXPR CategoryItem::CategoryItem(
    ::_pbi::ConstantInitialized): _impl_{
    /*decltype(_impl_.sub_item_ids_)*/ {}
  ,/* _impl_._sub_item_ids_cached_byte_size_ = */ { 0 }

  , /*decltype(_impl_.name_)*/ {
    &::_pbi::fixed_address_empty_string, ::_pbi::ConstantInitialized {}
  }

  , /*decltype(_impl_.resource_path_)*/ {
    &::_pbi::fixed_address_empty_string, ::_pbi::ConstantInitialized {}
  }

  , /*decltype(_impl_.poster_path_)*/ {
    &::_pbi::fixed_address_empty_string, ::_pbi::ConstantInitialized {}
  }

  , /*decltype(_impl_.introduce_)*/ {
    &::_pbi::fixed_address_empty_string, ::_pbi::ConstantInitialized {}
  }

  , /*decltype(_impl_.id_)*/ ::int64_t{0}

  , /*decltype(_impl_.creator_)*/ ::int64_t{0}

  , /*decltype(_impl_.parent_id_)*/ ::int64_t{0}

  , /*decltype(_impl_.type_id_)*/ 0

  , /*decltype(_impl_._cached_size_)*/{}} {}
struct CategoryItemDefaultTypeInternal {
  PROTOBUF_CONSTEXPR CategoryItemDefaultTypeInternal() : _instance(::_pbi::ConstantInitialized{}) {}
  ~CategoryItemDefaultTypeInternal() {}
  union {
    CategoryItem _instance;
  };
};

PROTOBUF_ATTRIBUTE_NO_DESTROY PROTOBUF_CONSTINIT
    PROTOBUF_ATTRIBUTE_INIT_PRIORITY1 CategoryItemDefaultTypeInternal _CategoryItem_default_instance_;
PROTOBUF_CONSTEXPR SharedItem::SharedItem(
    ::_pbi::ConstantInitialized): _impl_{
    /*decltype(_impl_.share_id_)*/ {
    &::_pbi::fixed_address_empty_string, ::_pbi::ConstantInitialized {}
  }

  , /*decltype(_impl_.item_id_)*/ ::int64_t{0}

  , /*decltype(_impl_._cached_size_)*/{}} {}
struct SharedItemDefaultTypeInternal {
  PROTOBUF_CONSTEXPR SharedItemDefaultTypeInternal() : _instance(::_pbi::ConstantInitialized{}) {}
  ~SharedItemDefaultTypeInternal() {}
  union {
    SharedItem _instance;
  };
};

PROTOBUF_ATTRIBUTE_NO_DESTROY PROTOBUF_CONSTINIT
    PROTOBUF_ATTRIBUTE_INIT_PRIORITY1 SharedItemDefaultTypeInternal _SharedItem_default_instance_;
}  // namespace prpc
static ::_pb::Metadata file_level_metadata_category_2eproto[2];
static const ::_pb::EnumDescriptor* file_level_enum_descriptors_category_2eproto[1];
static constexpr const ::_pb::ServiceDescriptor**
    file_level_service_descriptors_category_2eproto = nullptr;
const ::uint32_t TableStruct_category_2eproto::offsets[] PROTOBUF_SECTION_VARIABLE(
    protodesc_cold) = {
    ~0u,  // no _has_bits_
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _internal_metadata_),
    ~0u,  // no _extensions_
    ~0u,  // no _oneof_case_
    ~0u,  // no _weak_field_map_
    ~0u,  // no _inlined_string_donated_
    ~0u,  // no _split_
    ~0u,  // no sizeof(Split)
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.id_),
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.type_id_),
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.creator_),
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.name_),
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.resource_path_),
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.poster_path_),
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.introduce_),
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.parent_id_),
    PROTOBUF_FIELD_OFFSET(::prpc::CategoryItem, _impl_.sub_item_ids_),
    ~0u,  // no _has_bits_
    PROTOBUF_FIELD_OFFSET(::prpc::SharedItem, _internal_metadata_),
    ~0u,  // no _extensions_
    ~0u,  // no _oneof_case_
    ~0u,  // no _weak_field_map_
    ~0u,  // no _inlined_string_donated_
    ~0u,  // no _split_
    ~0u,  // no sizeof(Split)
    PROTOBUF_FIELD_OFFSET(::prpc::SharedItem, _impl_.item_id_),
    PROTOBUF_FIELD_OFFSET(::prpc::SharedItem, _impl_.share_id_),
};

static const ::_pbi::MigrationSchema
    schemas[] PROTOBUF_SECTION_VARIABLE(protodesc_cold) = {
        { 0, -1, -1, sizeof(::prpc::CategoryItem)},
        { 17, -1, -1, sizeof(::prpc::SharedItem)},
};

static const ::_pb::Message* const file_default_instances[] = {
    &::prpc::_CategoryItem_default_instance_._instance,
    &::prpc::_SharedItem_default_instance_._instance,
};
const char descriptor_table_protodef_category_2eproto[] PROTOBUF_SECTION_VARIABLE(protodesc_cold) = {
    "\n\016category.proto\022\004prpc\"\223\002\n\014CategoryItem\022"
    "\n\n\002id\030\001 \001(\003\022(\n\007type_id\030\002 \001(\0162\027.prpc.Cate"
    "goryItem.Type\022\017\n\007creator\030\003 \001(\003\022\014\n\004name\030\004"
    " \001(\t\022\025\n\rresource_path\030\005 \001(\t\022\023\n\013poster_pa"
    "th\030\006 \001(\t\022\021\n\tintroduce\030\007 \001(\t\022\021\n\tparent_id"
    "\030\010 \001(\003\022\024\n\014sub_item_ids\030\t \003(\003\"F\n\004Type\022\013\n\007"
    "Unknown\020\000\022\010\n\004Home\020\001\022\r\n\tDirectory\020\002\022\t\n\005Vi"
    "deo\020\003\022\r\n\tOtherFile\020\004\"/\n\nSharedItem\022\017\n\007it"
    "em_id\030\001 \001(\003\022\020\n\010share_id\030\002 \001(\tB\010Z\006./prpcb"
    "\006proto3"
};
static ::absl::once_flag descriptor_table_category_2eproto_once;
const ::_pbi::DescriptorTable descriptor_table_category_2eproto = {
    false,
    false,
    367,
    descriptor_table_protodef_category_2eproto,
    "category.proto",
    &descriptor_table_category_2eproto_once,
    nullptr,
    0,
    2,
    schemas,
    file_default_instances,
    TableStruct_category_2eproto::offsets,
    file_level_metadata_category_2eproto,
    file_level_enum_descriptors_category_2eproto,
    file_level_service_descriptors_category_2eproto,
};

// This function exists to be marked as weak.
// It can significantly speed up compilation by breaking up LLVM's SCC
// in the .pb.cc translation units. Large translation units see a
// reduction of more than 35% of walltime for optimized builds. Without
// the weak attribute all the messages in the file, including all the
// vtables and everything they use become part of the same SCC through
// a cycle like:
// GetMetadata -> descriptor table -> default instances ->
//   vtables -> GetMetadata
// By adding a weak function here we break the connection from the
// individual vtables back into the descriptor table.
PROTOBUF_ATTRIBUTE_WEAK const ::_pbi::DescriptorTable* descriptor_table_category_2eproto_getter() {
  return &descriptor_table_category_2eproto;
}
// Force running AddDescriptors() at dynamic initialization time.
PROTOBUF_ATTRIBUTE_INIT_PRIORITY2
static ::_pbi::AddDescriptorsRunner dynamic_init_dummy_category_2eproto(&descriptor_table_category_2eproto);
namespace prpc {
const ::PROTOBUF_NAMESPACE_ID::EnumDescriptor* CategoryItem_Type_descriptor() {
  ::PROTOBUF_NAMESPACE_ID::internal::AssignDescriptors(&descriptor_table_category_2eproto);
  return file_level_enum_descriptors_category_2eproto[0];
}
bool CategoryItem_Type_IsValid(int value) {
  switch (value) {
    case 0:
    case 1:
    case 2:
    case 3:
    case 4:
      return true;
    default:
      return false;
  }
}
#if (__cplusplus < 201703) && \
  (!defined(_MSC_VER) || (_MSC_VER >= 1900 && _MSC_VER < 1912))

constexpr CategoryItem_Type CategoryItem::Unknown;
constexpr CategoryItem_Type CategoryItem::Home;
constexpr CategoryItem_Type CategoryItem::Directory;
constexpr CategoryItem_Type CategoryItem::Video;
constexpr CategoryItem_Type CategoryItem::OtherFile;
constexpr CategoryItem_Type CategoryItem::Type_MIN;
constexpr CategoryItem_Type CategoryItem::Type_MAX;
constexpr int CategoryItem::Type_ARRAYSIZE;

#endif  // (__cplusplus < 201703) &&
        // (!defined(_MSC_VER) || (_MSC_VER >= 1900 && _MSC_VER < 1912))
// ===================================================================

class CategoryItem::_Internal {
 public:
};

CategoryItem::CategoryItem(::PROTOBUF_NAMESPACE_ID::Arena* arena)
  : ::PROTOBUF_NAMESPACE_ID::Message(arena) {
  SharedCtor(arena);
  // @@protoc_insertion_point(arena_constructor:prpc.CategoryItem)
}
CategoryItem::CategoryItem(const CategoryItem& from)
  : ::PROTOBUF_NAMESPACE_ID::Message() {
  CategoryItem* const _this = this; (void)_this;
  new (&_impl_) Impl_{
      decltype(_impl_.sub_item_ids_) { from._impl_.sub_item_ids_ }
    ,/* _impl_._sub_item_ids_cached_byte_size_ = */ { 0 }

    , decltype(_impl_.name_) {}

    , decltype(_impl_.resource_path_) {}

    , decltype(_impl_.poster_path_) {}

    , decltype(_impl_.introduce_) {}

    , decltype(_impl_.id_) {}

    , decltype(_impl_.creator_) {}

    , decltype(_impl_.parent_id_) {}

    , decltype(_impl_.type_id_) {}

    , /*decltype(_impl_._cached_size_)*/{}};

  _internal_metadata_.MergeFrom<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>(from._internal_metadata_);
  _impl_.name_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.name_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (!from._internal_name().empty()) {
    _this->_impl_.name_.Set(from._internal_name(), _this->GetArenaForAllocation());
  }
  _impl_.resource_path_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.resource_path_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (!from._internal_resource_path().empty()) {
    _this->_impl_.resource_path_.Set(from._internal_resource_path(), _this->GetArenaForAllocation());
  }
  _impl_.poster_path_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.poster_path_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (!from._internal_poster_path().empty()) {
    _this->_impl_.poster_path_.Set(from._internal_poster_path(), _this->GetArenaForAllocation());
  }
  _impl_.introduce_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.introduce_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (!from._internal_introduce().empty()) {
    _this->_impl_.introduce_.Set(from._internal_introduce(), _this->GetArenaForAllocation());
  }
  ::memcpy(&_impl_.id_, &from._impl_.id_,
    static_cast<::size_t>(reinterpret_cast<char*>(&_impl_.type_id_) -
    reinterpret_cast<char*>(&_impl_.id_)) + sizeof(_impl_.type_id_));
  // @@protoc_insertion_point(copy_constructor:prpc.CategoryItem)
}

inline void CategoryItem::SharedCtor(::_pb::Arena* arena) {
  (void)arena;
  new (&_impl_) Impl_{
      decltype(_impl_.sub_item_ids_) { arena }
    ,/* _impl_._sub_item_ids_cached_byte_size_ = */ { 0 }

    , decltype(_impl_.name_) {}

    , decltype(_impl_.resource_path_) {}

    , decltype(_impl_.poster_path_) {}

    , decltype(_impl_.introduce_) {}

    , decltype(_impl_.id_) { ::int64_t{0} }

    , decltype(_impl_.creator_) { ::int64_t{0} }

    , decltype(_impl_.parent_id_) { ::int64_t{0} }

    , decltype(_impl_.type_id_) { 0 }

    , /*decltype(_impl_._cached_size_)*/{}
  };
  _impl_.name_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.name_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  _impl_.resource_path_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.resource_path_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  _impl_.poster_path_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.poster_path_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  _impl_.introduce_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.introduce_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
}

CategoryItem::~CategoryItem() {
  // @@protoc_insertion_point(destructor:prpc.CategoryItem)
  if (auto *arena = _internal_metadata_.DeleteReturnArena<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>()) {
  (void)arena;
    return;
  }
  SharedDtor();
}

inline void CategoryItem::SharedDtor() {
  ABSL_DCHECK(GetArenaForAllocation() == nullptr);
  _impl_.sub_item_ids_.~RepeatedField();
  _impl_.name_.Destroy();
  _impl_.resource_path_.Destroy();
  _impl_.poster_path_.Destroy();
  _impl_.introduce_.Destroy();
}

void CategoryItem::SetCachedSize(int size) const {
  _impl_._cached_size_.Set(size);
}

void CategoryItem::Clear() {
// @@protoc_insertion_point(message_clear_start:prpc.CategoryItem)
  ::uint32_t cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  _impl_.sub_item_ids_.Clear();
  _impl_.name_.ClearToEmpty();
  _impl_.resource_path_.ClearToEmpty();
  _impl_.poster_path_.ClearToEmpty();
  _impl_.introduce_.ClearToEmpty();
  ::memset(&_impl_.id_, 0, static_cast<::size_t>(
      reinterpret_cast<char*>(&_impl_.type_id_) -
      reinterpret_cast<char*>(&_impl_.id_)) + sizeof(_impl_.type_id_));
  _internal_metadata_.Clear<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>();
}

const char* CategoryItem::_InternalParse(const char* ptr, ::_pbi::ParseContext* ctx) {
#define CHK_(x) if (PROTOBUF_PREDICT_FALSE(!(x))) goto failure
  while (!ctx->Done(&ptr)) {
    ::uint32_t tag;
    ptr = ::_pbi::ReadTag(ptr, &tag);
    switch (tag >> 3) {
      // int64 id = 1;
      case 1:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 8)) {
          _impl_.id_ = ::PROTOBUF_NAMESPACE_ID::internal::ReadVarint64(&ptr);
          CHK_(ptr);
        } else {
          goto handle_unusual;
        }
        continue;
      // .prpc.CategoryItem.Type type_id = 2;
      case 2:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 16)) {
          ::uint32_t val = ::PROTOBUF_NAMESPACE_ID::internal::ReadVarint32(&ptr);
          CHK_(ptr);
          _internal_set_type_id(static_cast<::prpc::CategoryItem_Type>(val));
        } else {
          goto handle_unusual;
        }
        continue;
      // int64 creator = 3;
      case 3:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 24)) {
          _impl_.creator_ = ::PROTOBUF_NAMESPACE_ID::internal::ReadVarint64(&ptr);
          CHK_(ptr);
        } else {
          goto handle_unusual;
        }
        continue;
      // string name = 4;
      case 4:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 34)) {
          auto str = _internal_mutable_name();
          ptr = ::_pbi::InlineGreedyStringParser(str, ptr, ctx);
          CHK_(ptr);
          CHK_(::_pbi::VerifyUTF8(str, "prpc.CategoryItem.name"));
        } else {
          goto handle_unusual;
        }
        continue;
      // string resource_path = 5;
      case 5:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 42)) {
          auto str = _internal_mutable_resource_path();
          ptr = ::_pbi::InlineGreedyStringParser(str, ptr, ctx);
          CHK_(ptr);
          CHK_(::_pbi::VerifyUTF8(str, "prpc.CategoryItem.resource_path"));
        } else {
          goto handle_unusual;
        }
        continue;
      // string poster_path = 6;
      case 6:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 50)) {
          auto str = _internal_mutable_poster_path();
          ptr = ::_pbi::InlineGreedyStringParser(str, ptr, ctx);
          CHK_(ptr);
          CHK_(::_pbi::VerifyUTF8(str, "prpc.CategoryItem.poster_path"));
        } else {
          goto handle_unusual;
        }
        continue;
      // string introduce = 7;
      case 7:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 58)) {
          auto str = _internal_mutable_introduce();
          ptr = ::_pbi::InlineGreedyStringParser(str, ptr, ctx);
          CHK_(ptr);
          CHK_(::_pbi::VerifyUTF8(str, "prpc.CategoryItem.introduce"));
        } else {
          goto handle_unusual;
        }
        continue;
      // int64 parent_id = 8;
      case 8:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 64)) {
          _impl_.parent_id_ = ::PROTOBUF_NAMESPACE_ID::internal::ReadVarint64(&ptr);
          CHK_(ptr);
        } else {
          goto handle_unusual;
        }
        continue;
      // repeated int64 sub_item_ids = 9;
      case 9:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 74)) {
          ptr = ::PROTOBUF_NAMESPACE_ID::internal::PackedInt64Parser(_internal_mutable_sub_item_ids(), ptr, ctx);
          CHK_(ptr);
        } else if (static_cast<::uint8_t>(tag) == 72) {
          _internal_add_sub_item_ids(::PROTOBUF_NAMESPACE_ID::internal::ReadVarint64(&ptr));
          CHK_(ptr);
        } else {
          goto handle_unusual;
        }
        continue;
      default:
        goto handle_unusual;
    }  // switch
  handle_unusual:
    if ((tag == 0) || ((tag & 7) == 4)) {
      CHK_(ptr);
      ctx->SetLastTag(tag);
      goto message_done;
    }
    ptr = UnknownFieldParse(
        tag,
        _internal_metadata_.mutable_unknown_fields<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>(),
        ptr, ctx);
    CHK_(ptr != nullptr);
  }  // while
message_done:
  return ptr;
failure:
  ptr = nullptr;
  goto message_done;
#undef CHK_
}

::uint8_t* CategoryItem::_InternalSerialize(
    ::uint8_t* target, ::PROTOBUF_NAMESPACE_ID::io::EpsCopyOutputStream* stream) const {
  // @@protoc_insertion_point(serialize_to_array_start:prpc.CategoryItem)
  ::uint32_t cached_has_bits = 0;
  (void) cached_has_bits;

  // int64 id = 1;
  if (this->_internal_id() != 0) {
    target = stream->EnsureSpace(target);
    target = ::_pbi::WireFormatLite::WriteInt64ToArray(
        1, this->_internal_id(), target);
  }

  // .prpc.CategoryItem.Type type_id = 2;
  if (this->_internal_type_id() != 0) {
    target = stream->EnsureSpace(target);
    target = ::_pbi::WireFormatLite::WriteEnumToArray(
        2, this->_internal_type_id(), target);
  }

  // int64 creator = 3;
  if (this->_internal_creator() != 0) {
    target = stream->EnsureSpace(target);
    target = ::_pbi::WireFormatLite::WriteInt64ToArray(
        3, this->_internal_creator(), target);
  }

  // string name = 4;
  if (!this->_internal_name().empty()) {
    const std::string& _s = this->_internal_name();
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
        _s.data(), static_cast<int>(_s.length()), ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE, "prpc.CategoryItem.name");
    target = stream->WriteStringMaybeAliased(4, _s, target);
  }

  // string resource_path = 5;
  if (!this->_internal_resource_path().empty()) {
    const std::string& _s = this->_internal_resource_path();
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
        _s.data(), static_cast<int>(_s.length()), ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE, "prpc.CategoryItem.resource_path");
    target = stream->WriteStringMaybeAliased(5, _s, target);
  }

  // string poster_path = 6;
  if (!this->_internal_poster_path().empty()) {
    const std::string& _s = this->_internal_poster_path();
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
        _s.data(), static_cast<int>(_s.length()), ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE, "prpc.CategoryItem.poster_path");
    target = stream->WriteStringMaybeAliased(6, _s, target);
  }

  // string introduce = 7;
  if (!this->_internal_introduce().empty()) {
    const std::string& _s = this->_internal_introduce();
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
        _s.data(), static_cast<int>(_s.length()), ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE, "prpc.CategoryItem.introduce");
    target = stream->WriteStringMaybeAliased(7, _s, target);
  }

  // int64 parent_id = 8;
  if (this->_internal_parent_id() != 0) {
    target = stream->EnsureSpace(target);
    target = ::_pbi::WireFormatLite::WriteInt64ToArray(
        8, this->_internal_parent_id(), target);
  }

  // repeated int64 sub_item_ids = 9;
  {
    int byte_size = _impl_._sub_item_ids_cached_byte_size_.Get();
    if (byte_size > 0) {
      target = stream->WriteInt64Packed(9, _internal_sub_item_ids(),
                                                 byte_size, target);
    }
  }

  if (PROTOBUF_PREDICT_FALSE(_internal_metadata_.have_unknown_fields())) {
    target = ::_pbi::WireFormat::InternalSerializeUnknownFieldsToArray(
        _internal_metadata_.unknown_fields<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>(::PROTOBUF_NAMESPACE_ID::UnknownFieldSet::default_instance), target, stream);
  }
  // @@protoc_insertion_point(serialize_to_array_end:prpc.CategoryItem)
  return target;
}

::size_t CategoryItem::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:prpc.CategoryItem)
  ::size_t total_size = 0;

  ::uint32_t cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  // repeated int64 sub_item_ids = 9;
  {
    std::size_t data_size = ::_pbi::WireFormatLite::Int64Size(this->_impl_.sub_item_ids_)
    ;
    _impl_._sub_item_ids_cached_byte_size_.Set(::_pbi::ToCachedSize(data_size));
    std::size_t tag_size = data_size == 0
        ? 0
        : 1 + ::_pbi::WireFormatLite::Int32Size(
                            static_cast<int32_t>(data_size))
    ;
    total_size += tag_size + data_size;
  }

  // string name = 4;
  if (!this->_internal_name().empty()) {
    total_size += 1 + ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::StringSize(
                                    this->_internal_name());
  }

  // string resource_path = 5;
  if (!this->_internal_resource_path().empty()) {
    total_size += 1 + ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::StringSize(
                                    this->_internal_resource_path());
  }

  // string poster_path = 6;
  if (!this->_internal_poster_path().empty()) {
    total_size += 1 + ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::StringSize(
                                    this->_internal_poster_path());
  }

  // string introduce = 7;
  if (!this->_internal_introduce().empty()) {
    total_size += 1 + ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::StringSize(
                                    this->_internal_introduce());
  }

  // int64 id = 1;
  if (this->_internal_id() != 0) {
    total_size += ::_pbi::WireFormatLite::Int64SizePlusOne(
        this->_internal_id());
  }

  // int64 creator = 3;
  if (this->_internal_creator() != 0) {
    total_size += ::_pbi::WireFormatLite::Int64SizePlusOne(
        this->_internal_creator());
  }

  // int64 parent_id = 8;
  if (this->_internal_parent_id() != 0) {
    total_size += ::_pbi::WireFormatLite::Int64SizePlusOne(
        this->_internal_parent_id());
  }

  // .prpc.CategoryItem.Type type_id = 2;
  if (this->_internal_type_id() != 0) {
    total_size += 1 +
                  ::_pbi::WireFormatLite::EnumSize(this->_internal_type_id());
  }

  return MaybeComputeUnknownFieldsSize(total_size, &_impl_._cached_size_);
}

const ::PROTOBUF_NAMESPACE_ID::Message::ClassData CategoryItem::_class_data_ = {
    ::PROTOBUF_NAMESPACE_ID::Message::CopyWithSourceCheck,
    CategoryItem::MergeImpl
};
const ::PROTOBUF_NAMESPACE_ID::Message::ClassData*CategoryItem::GetClassData() const { return &_class_data_; }


void CategoryItem::MergeImpl(::PROTOBUF_NAMESPACE_ID::Message& to_msg, const ::PROTOBUF_NAMESPACE_ID::Message& from_msg) {
  auto* const _this = static_cast<CategoryItem*>(&to_msg);
  auto& from = static_cast<const CategoryItem&>(from_msg);
  // @@protoc_insertion_point(class_specific_merge_from_start:prpc.CategoryItem)
  ABSL_DCHECK_NE(&from, _this);
  ::uint32_t cached_has_bits = 0;
  (void) cached_has_bits;

  _this->_impl_.sub_item_ids_.MergeFrom(from._impl_.sub_item_ids_);
  if (!from._internal_name().empty()) {
    _this->_internal_set_name(from._internal_name());
  }
  if (!from._internal_resource_path().empty()) {
    _this->_internal_set_resource_path(from._internal_resource_path());
  }
  if (!from._internal_poster_path().empty()) {
    _this->_internal_set_poster_path(from._internal_poster_path());
  }
  if (!from._internal_introduce().empty()) {
    _this->_internal_set_introduce(from._internal_introduce());
  }
  if (from._internal_id() != 0) {
    _this->_internal_set_id(from._internal_id());
  }
  if (from._internal_creator() != 0) {
    _this->_internal_set_creator(from._internal_creator());
  }
  if (from._internal_parent_id() != 0) {
    _this->_internal_set_parent_id(from._internal_parent_id());
  }
  if (from._internal_type_id() != 0) {
    _this->_internal_set_type_id(from._internal_type_id());
  }
  _this->_internal_metadata_.MergeFrom<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>(from._internal_metadata_);
}

void CategoryItem::CopyFrom(const CategoryItem& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:prpc.CategoryItem)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool CategoryItem::IsInitialized() const {
  return true;
}

void CategoryItem::InternalSwap(CategoryItem* other) {
  using std::swap;
  auto* lhs_arena = GetArenaForAllocation();
  auto* rhs_arena = other->GetArenaForAllocation();
  _internal_metadata_.InternalSwap(&other->_internal_metadata_);
  _impl_.sub_item_ids_.InternalSwap(&other->_impl_.sub_item_ids_);
  ::_pbi::ArenaStringPtr::InternalSwap(&_impl_.name_, lhs_arena,
                                       &other->_impl_.name_, rhs_arena);
  ::_pbi::ArenaStringPtr::InternalSwap(&_impl_.resource_path_, lhs_arena,
                                       &other->_impl_.resource_path_, rhs_arena);
  ::_pbi::ArenaStringPtr::InternalSwap(&_impl_.poster_path_, lhs_arena,
                                       &other->_impl_.poster_path_, rhs_arena);
  ::_pbi::ArenaStringPtr::InternalSwap(&_impl_.introduce_, lhs_arena,
                                       &other->_impl_.introduce_, rhs_arena);
  ::PROTOBUF_NAMESPACE_ID::internal::memswap<
      PROTOBUF_FIELD_OFFSET(CategoryItem, _impl_.type_id_)
      + sizeof(CategoryItem::_impl_.type_id_)
      - PROTOBUF_FIELD_OFFSET(CategoryItem, _impl_.id_)>(
          reinterpret_cast<char*>(&_impl_.id_),
          reinterpret_cast<char*>(&other->_impl_.id_));
}

::PROTOBUF_NAMESPACE_ID::Metadata CategoryItem::GetMetadata() const {
  return ::_pbi::AssignDescriptors(
      &descriptor_table_category_2eproto_getter, &descriptor_table_category_2eproto_once,
      file_level_metadata_category_2eproto[0]);
}
// ===================================================================

class SharedItem::_Internal {
 public:
};

SharedItem::SharedItem(::PROTOBUF_NAMESPACE_ID::Arena* arena)
  : ::PROTOBUF_NAMESPACE_ID::Message(arena) {
  SharedCtor(arena);
  // @@protoc_insertion_point(arena_constructor:prpc.SharedItem)
}
SharedItem::SharedItem(const SharedItem& from)
  : ::PROTOBUF_NAMESPACE_ID::Message() {
  SharedItem* const _this = this; (void)_this;
  new (&_impl_) Impl_{
      decltype(_impl_.share_id_) {}

    , decltype(_impl_.item_id_) {}

    , /*decltype(_impl_._cached_size_)*/{}};

  _internal_metadata_.MergeFrom<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>(from._internal_metadata_);
  _impl_.share_id_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.share_id_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  if (!from._internal_share_id().empty()) {
    _this->_impl_.share_id_.Set(from._internal_share_id(), _this->GetArenaForAllocation());
  }
  _this->_impl_.item_id_ = from._impl_.item_id_;
  // @@protoc_insertion_point(copy_constructor:prpc.SharedItem)
}

inline void SharedItem::SharedCtor(::_pb::Arena* arena) {
  (void)arena;
  new (&_impl_) Impl_{
      decltype(_impl_.share_id_) {}

    , decltype(_impl_.item_id_) { ::int64_t{0} }

    , /*decltype(_impl_._cached_size_)*/{}
  };
  _impl_.share_id_.InitDefault();
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        _impl_.share_id_.Set("", GetArenaForAllocation());
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
}

SharedItem::~SharedItem() {
  // @@protoc_insertion_point(destructor:prpc.SharedItem)
  if (auto *arena = _internal_metadata_.DeleteReturnArena<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>()) {
  (void)arena;
    return;
  }
  SharedDtor();
}

inline void SharedItem::SharedDtor() {
  ABSL_DCHECK(GetArenaForAllocation() == nullptr);
  _impl_.share_id_.Destroy();
}

void SharedItem::SetCachedSize(int size) const {
  _impl_._cached_size_.Set(size);
}

void SharedItem::Clear() {
// @@protoc_insertion_point(message_clear_start:prpc.SharedItem)
  ::uint32_t cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  _impl_.share_id_.ClearToEmpty();
  _impl_.item_id_ = ::int64_t{0};
  _internal_metadata_.Clear<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>();
}

const char* SharedItem::_InternalParse(const char* ptr, ::_pbi::ParseContext* ctx) {
#define CHK_(x) if (PROTOBUF_PREDICT_FALSE(!(x))) goto failure
  while (!ctx->Done(&ptr)) {
    ::uint32_t tag;
    ptr = ::_pbi::ReadTag(ptr, &tag);
    switch (tag >> 3) {
      // int64 item_id = 1;
      case 1:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 8)) {
          _impl_.item_id_ = ::PROTOBUF_NAMESPACE_ID::internal::ReadVarint64(&ptr);
          CHK_(ptr);
        } else {
          goto handle_unusual;
        }
        continue;
      // string share_id = 2;
      case 2:
        if (PROTOBUF_PREDICT_TRUE(static_cast<::uint8_t>(tag) == 18)) {
          auto str = _internal_mutable_share_id();
          ptr = ::_pbi::InlineGreedyStringParser(str, ptr, ctx);
          CHK_(ptr);
          CHK_(::_pbi::VerifyUTF8(str, "prpc.SharedItem.share_id"));
        } else {
          goto handle_unusual;
        }
        continue;
      default:
        goto handle_unusual;
    }  // switch
  handle_unusual:
    if ((tag == 0) || ((tag & 7) == 4)) {
      CHK_(ptr);
      ctx->SetLastTag(tag);
      goto message_done;
    }
    ptr = UnknownFieldParse(
        tag,
        _internal_metadata_.mutable_unknown_fields<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>(),
        ptr, ctx);
    CHK_(ptr != nullptr);
  }  // while
message_done:
  return ptr;
failure:
  ptr = nullptr;
  goto message_done;
#undef CHK_
}

::uint8_t* SharedItem::_InternalSerialize(
    ::uint8_t* target, ::PROTOBUF_NAMESPACE_ID::io::EpsCopyOutputStream* stream) const {
  // @@protoc_insertion_point(serialize_to_array_start:prpc.SharedItem)
  ::uint32_t cached_has_bits = 0;
  (void) cached_has_bits;

  // int64 item_id = 1;
  if (this->_internal_item_id() != 0) {
    target = stream->EnsureSpace(target);
    target = ::_pbi::WireFormatLite::WriteInt64ToArray(
        1, this->_internal_item_id(), target);
  }

  // string share_id = 2;
  if (!this->_internal_share_id().empty()) {
    const std::string& _s = this->_internal_share_id();
    ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::VerifyUtf8String(
        _s.data(), static_cast<int>(_s.length()), ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::SERIALIZE, "prpc.SharedItem.share_id");
    target = stream->WriteStringMaybeAliased(2, _s, target);
  }

  if (PROTOBUF_PREDICT_FALSE(_internal_metadata_.have_unknown_fields())) {
    target = ::_pbi::WireFormat::InternalSerializeUnknownFieldsToArray(
        _internal_metadata_.unknown_fields<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>(::PROTOBUF_NAMESPACE_ID::UnknownFieldSet::default_instance), target, stream);
  }
  // @@protoc_insertion_point(serialize_to_array_end:prpc.SharedItem)
  return target;
}

::size_t SharedItem::ByteSizeLong() const {
// @@protoc_insertion_point(message_byte_size_start:prpc.SharedItem)
  ::size_t total_size = 0;

  ::uint32_t cached_has_bits = 0;
  // Prevent compiler warnings about cached_has_bits being unused
  (void) cached_has_bits;

  // string share_id = 2;
  if (!this->_internal_share_id().empty()) {
    total_size += 1 + ::PROTOBUF_NAMESPACE_ID::internal::WireFormatLite::StringSize(
                                    this->_internal_share_id());
  }

  // int64 item_id = 1;
  if (this->_internal_item_id() != 0) {
    total_size += ::_pbi::WireFormatLite::Int64SizePlusOne(
        this->_internal_item_id());
  }

  return MaybeComputeUnknownFieldsSize(total_size, &_impl_._cached_size_);
}

const ::PROTOBUF_NAMESPACE_ID::Message::ClassData SharedItem::_class_data_ = {
    ::PROTOBUF_NAMESPACE_ID::Message::CopyWithSourceCheck,
    SharedItem::MergeImpl
};
const ::PROTOBUF_NAMESPACE_ID::Message::ClassData*SharedItem::GetClassData() const { return &_class_data_; }


void SharedItem::MergeImpl(::PROTOBUF_NAMESPACE_ID::Message& to_msg, const ::PROTOBUF_NAMESPACE_ID::Message& from_msg) {
  auto* const _this = static_cast<SharedItem*>(&to_msg);
  auto& from = static_cast<const SharedItem&>(from_msg);
  // @@protoc_insertion_point(class_specific_merge_from_start:prpc.SharedItem)
  ABSL_DCHECK_NE(&from, _this);
  ::uint32_t cached_has_bits = 0;
  (void) cached_has_bits;

  if (!from._internal_share_id().empty()) {
    _this->_internal_set_share_id(from._internal_share_id());
  }
  if (from._internal_item_id() != 0) {
    _this->_internal_set_item_id(from._internal_item_id());
  }
  _this->_internal_metadata_.MergeFrom<::PROTOBUF_NAMESPACE_ID::UnknownFieldSet>(from._internal_metadata_);
}

void SharedItem::CopyFrom(const SharedItem& from) {
// @@protoc_insertion_point(class_specific_copy_from_start:prpc.SharedItem)
  if (&from == this) return;
  Clear();
  MergeFrom(from);
}

bool SharedItem::IsInitialized() const {
  return true;
}

void SharedItem::InternalSwap(SharedItem* other) {
  using std::swap;
  auto* lhs_arena = GetArenaForAllocation();
  auto* rhs_arena = other->GetArenaForAllocation();
  _internal_metadata_.InternalSwap(&other->_internal_metadata_);
  ::_pbi::ArenaStringPtr::InternalSwap(&_impl_.share_id_, lhs_arena,
                                       &other->_impl_.share_id_, rhs_arena);

  swap(_impl_.item_id_, other->_impl_.item_id_);
}

::PROTOBUF_NAMESPACE_ID::Metadata SharedItem::GetMetadata() const {
  return ::_pbi::AssignDescriptors(
      &descriptor_table_category_2eproto_getter, &descriptor_table_category_2eproto_once,
      file_level_metadata_category_2eproto[1]);
}
// @@protoc_insertion_point(namespace_scope)
}  // namespace prpc
PROTOBUF_NAMESPACE_OPEN
template<> PROTOBUF_NOINLINE ::prpc::CategoryItem*
Arena::CreateMaybeMessage< ::prpc::CategoryItem >(Arena* arena) {
  return Arena::CreateMessageInternal< ::prpc::CategoryItem >(arena);
}
template<> PROTOBUF_NOINLINE ::prpc::SharedItem*
Arena::CreateMaybeMessage< ::prpc::SharedItem >(Arena* arena) {
  return Arena::CreateMessageInternal< ::prpc::SharedItem >(arena);
}
PROTOBUF_NAMESPACE_CLOSE
// @@protoc_insertion_point(global_scope)
#include "google/protobuf/port_undef.inc"
