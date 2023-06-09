// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: category.proto

#ifndef GOOGLE_PROTOBUF_INCLUDED_category_2eproto_2epb_2eh
#define GOOGLE_PROTOBUF_INCLUDED_category_2eproto_2epb_2eh

#include <limits>
#include <string>
#include <type_traits>

#include "google/protobuf/port_def.inc"
#if PROTOBUF_VERSION < 4022000
#error "This file was generated by a newer version of protoc which is"
#error "incompatible with your Protocol Buffer headers. Please update"
#error "your headers."
#endif  // PROTOBUF_VERSION

#if 4022002 < PROTOBUF_MIN_PROTOC_VERSION
#error "This file was generated by an older version of protoc which is"
#error "incompatible with your Protocol Buffer headers. Please"
#error "regenerate this file with a newer version of protoc."
#endif  // PROTOBUF_MIN_PROTOC_VERSION
#include "google/protobuf/port_undef.inc"
#include "google/protobuf/io/coded_stream.h"
#include "google/protobuf/arena.h"
#include "google/protobuf/arenastring.h"
#include "google/protobuf/generated_message_util.h"
#include "google/protobuf/metadata_lite.h"
#include "google/protobuf/generated_message_reflection.h"
#include "google/protobuf/message.h"
#include "google/protobuf/repeated_field.h"  // IWYU pragma: export
#include "google/protobuf/extension_set.h"  // IWYU pragma: export
#include "google/protobuf/generated_enum_reflection.h"
#include "google/protobuf/unknown_field_set.h"
// @@protoc_insertion_point(includes)

// Must be included last.
#include "google/protobuf/port_def.inc"

#define PROTOBUF_INTERNAL_EXPORT_category_2eproto

PROTOBUF_NAMESPACE_OPEN
namespace internal {
class AnyMetadata;
}  // namespace internal
PROTOBUF_NAMESPACE_CLOSE

// Internal implementation detail -- do not use these members.
struct TableStruct_category_2eproto {
  static const ::uint32_t offsets[];
};
extern const ::PROTOBUF_NAMESPACE_ID::internal::DescriptorTable
    descriptor_table_category_2eproto;
namespace prpc {
class CategoryItem;
struct CategoryItemDefaultTypeInternal;
extern CategoryItemDefaultTypeInternal _CategoryItem_default_instance_;
class SharedItem;
struct SharedItemDefaultTypeInternal;
extern SharedItemDefaultTypeInternal _SharedItem_default_instance_;
}  // namespace prpc
PROTOBUF_NAMESPACE_OPEN
template <>
::prpc::CategoryItem* Arena::CreateMaybeMessage<::prpc::CategoryItem>(Arena*);
template <>
::prpc::SharedItem* Arena::CreateMaybeMessage<::prpc::SharedItem>(Arena*);
PROTOBUF_NAMESPACE_CLOSE

namespace prpc {
enum CategoryItem_Type : int {
  CategoryItem_Type_Unknown = 0,
  CategoryItem_Type_Home = 1,
  CategoryItem_Type_Directory = 2,
  CategoryItem_Type_Video = 3,
  CategoryItem_Type_OtherFile = 4,
  CategoryItem_Type_CategoryItem_Type_INT_MIN_SENTINEL_DO_NOT_USE_ =
      std::numeric_limits<::int32_t>::min(),
  CategoryItem_Type_CategoryItem_Type_INT_MAX_SENTINEL_DO_NOT_USE_ =
      std::numeric_limits<::int32_t>::max(),
};

bool CategoryItem_Type_IsValid(int value);
constexpr CategoryItem_Type CategoryItem_Type_Type_MIN = static_cast<CategoryItem_Type>(0);
constexpr CategoryItem_Type CategoryItem_Type_Type_MAX = static_cast<CategoryItem_Type>(4);
constexpr int CategoryItem_Type_Type_ARRAYSIZE = 4 + 1;
const ::PROTOBUF_NAMESPACE_ID::EnumDescriptor*
CategoryItem_Type_descriptor();
template <typename T>
const std::string& CategoryItem_Type_Name(T value) {
  static_assert(std::is_same<T, CategoryItem_Type>::value ||
                    std::is_integral<T>::value,
                "Incorrect type passed to Type_Name().");
  return CategoryItem_Type_Name(static_cast<CategoryItem_Type>(value));
}
template <>
inline const std::string& CategoryItem_Type_Name(CategoryItem_Type value) {
  return ::PROTOBUF_NAMESPACE_ID::internal::NameOfDenseEnum<CategoryItem_Type_descriptor,
                                                 0, 4>(
      static_cast<int>(value));
}
inline bool CategoryItem_Type_Parse(absl::string_view name, CategoryItem_Type* value) {
  return ::PROTOBUF_NAMESPACE_ID::internal::ParseNamedEnum<CategoryItem_Type>(
      CategoryItem_Type_descriptor(), name, value);
}

// ===================================================================


// -------------------------------------------------------------------

class CategoryItem final :
    public ::PROTOBUF_NAMESPACE_ID::Message /* @@protoc_insertion_point(class_definition:prpc.CategoryItem) */ {
 public:
  inline CategoryItem() : CategoryItem(nullptr) {}
  ~CategoryItem() override;
  explicit PROTOBUF_CONSTEXPR CategoryItem(::PROTOBUF_NAMESPACE_ID::internal::ConstantInitialized);

  CategoryItem(const CategoryItem& from);
  CategoryItem(CategoryItem&& from) noexcept
    : CategoryItem() {
    *this = ::std::move(from);
  }

  inline CategoryItem& operator=(const CategoryItem& from) {
    CopyFrom(from);
    return *this;
  }
  inline CategoryItem& operator=(CategoryItem&& from) noexcept {
    if (this == &from) return *this;
    if (GetOwningArena() == from.GetOwningArena()
  #ifdef PROTOBUF_FORCE_COPY_IN_MOVE
        && GetOwningArena() != nullptr
  #endif  // !PROTOBUF_FORCE_COPY_IN_MOVE
    ) {
      InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }

  static const ::PROTOBUF_NAMESPACE_ID::Descriptor* descriptor() {
    return GetDescriptor();
  }
  static const ::PROTOBUF_NAMESPACE_ID::Descriptor* GetDescriptor() {
    return default_instance().GetMetadata().descriptor;
  }
  static const ::PROTOBUF_NAMESPACE_ID::Reflection* GetReflection() {
    return default_instance().GetMetadata().reflection;
  }
  static const CategoryItem& default_instance() {
    return *internal_default_instance();
  }
  static inline const CategoryItem* internal_default_instance() {
    return reinterpret_cast<const CategoryItem*>(
               &_CategoryItem_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    0;

  friend void swap(CategoryItem& a, CategoryItem& b) {
    a.Swap(&b);
  }
  inline void Swap(CategoryItem* other) {
    if (other == this) return;
  #ifdef PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() != nullptr &&
        GetOwningArena() == other->GetOwningArena()) {
   #else  // PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() == other->GetOwningArena()) {
  #endif  // !PROTOBUF_FORCE_COPY_IN_SWAP
      InternalSwap(other);
    } else {
      ::PROTOBUF_NAMESPACE_ID::internal::GenericSwap(this, other);
    }
  }
  void UnsafeArenaSwap(CategoryItem* other) {
    if (other == this) return;
    ABSL_DCHECK(GetOwningArena() == other->GetOwningArena());
    InternalSwap(other);
  }

  // implements Message ----------------------------------------------

  CategoryItem* New(::PROTOBUF_NAMESPACE_ID::Arena* arena = nullptr) const final {
    return CreateMaybeMessage<CategoryItem>(arena);
  }
  using ::PROTOBUF_NAMESPACE_ID::Message::CopyFrom;
  void CopyFrom(const CategoryItem& from);
  using ::PROTOBUF_NAMESPACE_ID::Message::MergeFrom;
  void MergeFrom( const CategoryItem& from) {
    CategoryItem::MergeImpl(*this, from);
  }
  private:
  static void MergeImpl(::PROTOBUF_NAMESPACE_ID::Message& to_msg, const ::PROTOBUF_NAMESPACE_ID::Message& from_msg);
  public:
  PROTOBUF_ATTRIBUTE_REINITIALIZES void Clear() final;
  bool IsInitialized() const final;

  ::size_t ByteSizeLong() const final;
  const char* _InternalParse(const char* ptr, ::PROTOBUF_NAMESPACE_ID::internal::ParseContext* ctx) final;
  ::uint8_t* _InternalSerialize(
      ::uint8_t* target, ::PROTOBUF_NAMESPACE_ID::io::EpsCopyOutputStream* stream) const final;
  int GetCachedSize() const final { return _impl_._cached_size_.Get(); }

  private:
  void SharedCtor(::PROTOBUF_NAMESPACE_ID::Arena* arena);
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(CategoryItem* other);

  private:
  friend class ::PROTOBUF_NAMESPACE_ID::internal::AnyMetadata;
  static ::absl::string_view FullMessageName() {
    return "prpc.CategoryItem";
  }
  protected:
  explicit CategoryItem(::PROTOBUF_NAMESPACE_ID::Arena* arena);
  public:

  static const ClassData _class_data_;
  const ::PROTOBUF_NAMESPACE_ID::Message::ClassData*GetClassData() const final;

  ::PROTOBUF_NAMESPACE_ID::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  using Type = CategoryItem_Type;
  static constexpr Type Unknown = CategoryItem_Type_Unknown;
  static constexpr Type Home = CategoryItem_Type_Home;
  static constexpr Type Directory = CategoryItem_Type_Directory;
  static constexpr Type Video = CategoryItem_Type_Video;
  static constexpr Type OtherFile = CategoryItem_Type_OtherFile;
  static inline bool Type_IsValid(int value) {
    return CategoryItem_Type_IsValid(value);
  }
  static constexpr Type Type_MIN = CategoryItem_Type_Type_MIN;
  static constexpr Type Type_MAX = CategoryItem_Type_Type_MAX;
  static constexpr int Type_ARRAYSIZE = CategoryItem_Type_Type_ARRAYSIZE;
  static inline const ::PROTOBUF_NAMESPACE_ID::EnumDescriptor* Type_descriptor() {
    return CategoryItem_Type_descriptor();
  }
  template <typename T>
  static inline const std::string& Type_Name(T value) {
    return CategoryItem_Type_Name(value);
  }
  static inline bool Type_Parse(absl::string_view name, Type* value) {
    return CategoryItem_Type_Parse(name, value);
  }

  // accessors -------------------------------------------------------

  enum : int {
    kSubItemIdsFieldNumber = 9,
    kNameFieldNumber = 4,
    kResourcePathFieldNumber = 5,
    kPosterPathFieldNumber = 6,
    kIntroduceFieldNumber = 7,
    kIdFieldNumber = 1,
    kCreatorFieldNumber = 3,
    kParentIdFieldNumber = 8,
    kTypeIdFieldNumber = 2,
  };
  // repeated int64 sub_item_ids = 9;
  int sub_item_ids_size() const;
  private:
  int _internal_sub_item_ids_size() const;

  public:
  void clear_sub_item_ids() ;
  ::int64_t sub_item_ids(int index) const;
  void set_sub_item_ids(int index, ::int64_t value);
  void add_sub_item_ids(::int64_t value);
  const ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t>& sub_item_ids() const;
  ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t>* mutable_sub_item_ids();

  private:
  ::int64_t _internal_sub_item_ids(int index) const;
  void _internal_add_sub_item_ids(::int64_t value);
  const ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t>& _internal_sub_item_ids() const;
  ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t>* _internal_mutable_sub_item_ids();

  public:
  // string name = 4;
  void clear_name() ;
  const std::string& name() const;




  template <typename Arg_ = const std::string&, typename... Args_>
  void set_name(Arg_&& arg, Args_... args);
  std::string* mutable_name();
  PROTOBUF_NODISCARD std::string* release_name();
  void set_allocated_name(std::string* ptr);

  private:
  const std::string& _internal_name() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_name(
      const std::string& value);
  std::string* _internal_mutable_name();

  public:
  // string resource_path = 5;
  void clear_resource_path() ;
  const std::string& resource_path() const;




  template <typename Arg_ = const std::string&, typename... Args_>
  void set_resource_path(Arg_&& arg, Args_... args);
  std::string* mutable_resource_path();
  PROTOBUF_NODISCARD std::string* release_resource_path();
  void set_allocated_resource_path(std::string* ptr);

  private:
  const std::string& _internal_resource_path() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_resource_path(
      const std::string& value);
  std::string* _internal_mutable_resource_path();

  public:
  // string poster_path = 6;
  void clear_poster_path() ;
  const std::string& poster_path() const;




  template <typename Arg_ = const std::string&, typename... Args_>
  void set_poster_path(Arg_&& arg, Args_... args);
  std::string* mutable_poster_path();
  PROTOBUF_NODISCARD std::string* release_poster_path();
  void set_allocated_poster_path(std::string* ptr);

  private:
  const std::string& _internal_poster_path() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_poster_path(
      const std::string& value);
  std::string* _internal_mutable_poster_path();

  public:
  // string introduce = 7;
  void clear_introduce() ;
  const std::string& introduce() const;




  template <typename Arg_ = const std::string&, typename... Args_>
  void set_introduce(Arg_&& arg, Args_... args);
  std::string* mutable_introduce();
  PROTOBUF_NODISCARD std::string* release_introduce();
  void set_allocated_introduce(std::string* ptr);

  private:
  const std::string& _internal_introduce() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_introduce(
      const std::string& value);
  std::string* _internal_mutable_introduce();

  public:
  // int64 id = 1;
  void clear_id() ;
  ::int64_t id() const;
  void set_id(::int64_t value);

  private:
  ::int64_t _internal_id() const;
  void _internal_set_id(::int64_t value);

  public:
  // int64 creator = 3;
  void clear_creator() ;
  ::int64_t creator() const;
  void set_creator(::int64_t value);

  private:
  ::int64_t _internal_creator() const;
  void _internal_set_creator(::int64_t value);

  public:
  // int64 parent_id = 8;
  void clear_parent_id() ;
  ::int64_t parent_id() const;
  void set_parent_id(::int64_t value);

  private:
  ::int64_t _internal_parent_id() const;
  void _internal_set_parent_id(::int64_t value);

  public:
  // .prpc.CategoryItem.Type type_id = 2;
  void clear_type_id() ;
  ::prpc::CategoryItem_Type type_id() const;
  void set_type_id(::prpc::CategoryItem_Type value);

  private:
  ::prpc::CategoryItem_Type _internal_type_id() const;
  void _internal_set_type_id(::prpc::CategoryItem_Type value);

  public:
  // @@protoc_insertion_point(class_scope:prpc.CategoryItem)
 private:
  class _Internal;

  template <typename T> friend class ::PROTOBUF_NAMESPACE_ID::Arena::InternalHelper;
  typedef void InternalArenaConstructable_;
  typedef void DestructorSkippable_;
  struct Impl_ {
    ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t> sub_item_ids_;
    mutable ::PROTOBUF_NAMESPACE_ID::internal::CachedSize _sub_item_ids_cached_byte_size_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr name_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr resource_path_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr poster_path_;
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr introduce_;
    ::int64_t id_;
    ::int64_t creator_;
    ::int64_t parent_id_;
    int type_id_;
    mutable ::PROTOBUF_NAMESPACE_ID::internal::CachedSize _cached_size_;
  };
  union { Impl_ _impl_; };
  friend struct ::TableStruct_category_2eproto;
};// -------------------------------------------------------------------

class SharedItem final :
    public ::PROTOBUF_NAMESPACE_ID::Message /* @@protoc_insertion_point(class_definition:prpc.SharedItem) */ {
 public:
  inline SharedItem() : SharedItem(nullptr) {}
  ~SharedItem() override;
  explicit PROTOBUF_CONSTEXPR SharedItem(::PROTOBUF_NAMESPACE_ID::internal::ConstantInitialized);

  SharedItem(const SharedItem& from);
  SharedItem(SharedItem&& from) noexcept
    : SharedItem() {
    *this = ::std::move(from);
  }

  inline SharedItem& operator=(const SharedItem& from) {
    CopyFrom(from);
    return *this;
  }
  inline SharedItem& operator=(SharedItem&& from) noexcept {
    if (this == &from) return *this;
    if (GetOwningArena() == from.GetOwningArena()
  #ifdef PROTOBUF_FORCE_COPY_IN_MOVE
        && GetOwningArena() != nullptr
  #endif  // !PROTOBUF_FORCE_COPY_IN_MOVE
    ) {
      InternalSwap(&from);
    } else {
      CopyFrom(from);
    }
    return *this;
  }

  static const ::PROTOBUF_NAMESPACE_ID::Descriptor* descriptor() {
    return GetDescriptor();
  }
  static const ::PROTOBUF_NAMESPACE_ID::Descriptor* GetDescriptor() {
    return default_instance().GetMetadata().descriptor;
  }
  static const ::PROTOBUF_NAMESPACE_ID::Reflection* GetReflection() {
    return default_instance().GetMetadata().reflection;
  }
  static const SharedItem& default_instance() {
    return *internal_default_instance();
  }
  static inline const SharedItem* internal_default_instance() {
    return reinterpret_cast<const SharedItem*>(
               &_SharedItem_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    1;

  friend void swap(SharedItem& a, SharedItem& b) {
    a.Swap(&b);
  }
  inline void Swap(SharedItem* other) {
    if (other == this) return;
  #ifdef PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() != nullptr &&
        GetOwningArena() == other->GetOwningArena()) {
   #else  // PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() == other->GetOwningArena()) {
  #endif  // !PROTOBUF_FORCE_COPY_IN_SWAP
      InternalSwap(other);
    } else {
      ::PROTOBUF_NAMESPACE_ID::internal::GenericSwap(this, other);
    }
  }
  void UnsafeArenaSwap(SharedItem* other) {
    if (other == this) return;
    ABSL_DCHECK(GetOwningArena() == other->GetOwningArena());
    InternalSwap(other);
  }

  // implements Message ----------------------------------------------

  SharedItem* New(::PROTOBUF_NAMESPACE_ID::Arena* arena = nullptr) const final {
    return CreateMaybeMessage<SharedItem>(arena);
  }
  using ::PROTOBUF_NAMESPACE_ID::Message::CopyFrom;
  void CopyFrom(const SharedItem& from);
  using ::PROTOBUF_NAMESPACE_ID::Message::MergeFrom;
  void MergeFrom( const SharedItem& from) {
    SharedItem::MergeImpl(*this, from);
  }
  private:
  static void MergeImpl(::PROTOBUF_NAMESPACE_ID::Message& to_msg, const ::PROTOBUF_NAMESPACE_ID::Message& from_msg);
  public:
  PROTOBUF_ATTRIBUTE_REINITIALIZES void Clear() final;
  bool IsInitialized() const final;

  ::size_t ByteSizeLong() const final;
  const char* _InternalParse(const char* ptr, ::PROTOBUF_NAMESPACE_ID::internal::ParseContext* ctx) final;
  ::uint8_t* _InternalSerialize(
      ::uint8_t* target, ::PROTOBUF_NAMESPACE_ID::io::EpsCopyOutputStream* stream) const final;
  int GetCachedSize() const final { return _impl_._cached_size_.Get(); }

  private:
  void SharedCtor(::PROTOBUF_NAMESPACE_ID::Arena* arena);
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(SharedItem* other);

  private:
  friend class ::PROTOBUF_NAMESPACE_ID::internal::AnyMetadata;
  static ::absl::string_view FullMessageName() {
    return "prpc.SharedItem";
  }
  protected:
  explicit SharedItem(::PROTOBUF_NAMESPACE_ID::Arena* arena);
  public:

  static const ClassData _class_data_;
  const ::PROTOBUF_NAMESPACE_ID::Message::ClassData*GetClassData() const final;

  ::PROTOBUF_NAMESPACE_ID::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  enum : int {
    kShareIdFieldNumber = 2,
    kItemIdFieldNumber = 1,
  };
  // string share_id = 2;
  void clear_share_id() ;
  const std::string& share_id() const;




  template <typename Arg_ = const std::string&, typename... Args_>
  void set_share_id(Arg_&& arg, Args_... args);
  std::string* mutable_share_id();
  PROTOBUF_NODISCARD std::string* release_share_id();
  void set_allocated_share_id(std::string* ptr);

  private:
  const std::string& _internal_share_id() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_share_id(
      const std::string& value);
  std::string* _internal_mutable_share_id();

  public:
  // int64 item_id = 1;
  void clear_item_id() ;
  ::int64_t item_id() const;
  void set_item_id(::int64_t value);

  private:
  ::int64_t _internal_item_id() const;
  void _internal_set_item_id(::int64_t value);

  public:
  // @@protoc_insertion_point(class_scope:prpc.SharedItem)
 private:
  class _Internal;

  template <typename T> friend class ::PROTOBUF_NAMESPACE_ID::Arena::InternalHelper;
  typedef void InternalArenaConstructable_;
  typedef void DestructorSkippable_;
  struct Impl_ {
    ::PROTOBUF_NAMESPACE_ID::internal::ArenaStringPtr share_id_;
    ::int64_t item_id_;
    mutable ::PROTOBUF_NAMESPACE_ID::internal::CachedSize _cached_size_;
  };
  union { Impl_ _impl_; };
  friend struct ::TableStruct_category_2eproto;
};

// ===================================================================




// ===================================================================


#ifdef __GNUC__
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// -------------------------------------------------------------------

// CategoryItem

// int64 id = 1;
inline void CategoryItem::clear_id() {
  _impl_.id_ = ::int64_t{0};
}
inline ::int64_t CategoryItem::id() const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.id)
  return _internal_id();
}
inline void CategoryItem::set_id(::int64_t value) {
  _internal_set_id(value);
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.id)
}
inline ::int64_t CategoryItem::_internal_id() const {
  return _impl_.id_;
}
inline void CategoryItem::_internal_set_id(::int64_t value) {
  ;
  _impl_.id_ = value;
}

// .prpc.CategoryItem.Type type_id = 2;
inline void CategoryItem::clear_type_id() {
  _impl_.type_id_ = 0;
}
inline ::prpc::CategoryItem_Type CategoryItem::type_id() const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.type_id)
  return _internal_type_id();
}
inline void CategoryItem::set_type_id(::prpc::CategoryItem_Type value) {
   _internal_set_type_id(value);
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.type_id)
}
inline ::prpc::CategoryItem_Type CategoryItem::_internal_type_id() const {
  return static_cast<::prpc::CategoryItem_Type>(_impl_.type_id_);
}
inline void CategoryItem::_internal_set_type_id(::prpc::CategoryItem_Type value) {
  ;
  _impl_.type_id_ = value;
}

// int64 creator = 3;
inline void CategoryItem::clear_creator() {
  _impl_.creator_ = ::int64_t{0};
}
inline ::int64_t CategoryItem::creator() const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.creator)
  return _internal_creator();
}
inline void CategoryItem::set_creator(::int64_t value) {
  _internal_set_creator(value);
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.creator)
}
inline ::int64_t CategoryItem::_internal_creator() const {
  return _impl_.creator_;
}
inline void CategoryItem::_internal_set_creator(::int64_t value) {
  ;
  _impl_.creator_ = value;
}

// string name = 4;
inline void CategoryItem::clear_name() {
  _impl_.name_.ClearToEmpty();
}
inline const std::string& CategoryItem::name() const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.name)
  return _internal_name();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void CategoryItem::set_name(Arg_&& arg,
                                                     Args_... args) {
  ;
  _impl_.name_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.name)
}
inline std::string* CategoryItem::mutable_name() {
  std::string* _s = _internal_mutable_name();
  // @@protoc_insertion_point(field_mutable:prpc.CategoryItem.name)
  return _s;
}
inline const std::string& CategoryItem::_internal_name() const {
  return _impl_.name_.Get();
}
inline void CategoryItem::_internal_set_name(const std::string& value) {
  ;


  _impl_.name_.Set(value, GetArenaForAllocation());
}
inline std::string* CategoryItem::_internal_mutable_name() {
  ;
  return _impl_.name_.Mutable( GetArenaForAllocation());
}
inline std::string* CategoryItem::release_name() {
  // @@protoc_insertion_point(field_release:prpc.CategoryItem.name)
  return _impl_.name_.Release();
}
inline void CategoryItem::set_allocated_name(std::string* value) {
  _impl_.name_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.name_.IsDefault()) {
          _impl_.name_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:prpc.CategoryItem.name)
}

// string resource_path = 5;
inline void CategoryItem::clear_resource_path() {
  _impl_.resource_path_.ClearToEmpty();
}
inline const std::string& CategoryItem::resource_path() const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.resource_path)
  return _internal_resource_path();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void CategoryItem::set_resource_path(Arg_&& arg,
                                                     Args_... args) {
  ;
  _impl_.resource_path_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.resource_path)
}
inline std::string* CategoryItem::mutable_resource_path() {
  std::string* _s = _internal_mutable_resource_path();
  // @@protoc_insertion_point(field_mutable:prpc.CategoryItem.resource_path)
  return _s;
}
inline const std::string& CategoryItem::_internal_resource_path() const {
  return _impl_.resource_path_.Get();
}
inline void CategoryItem::_internal_set_resource_path(const std::string& value) {
  ;


  _impl_.resource_path_.Set(value, GetArenaForAllocation());
}
inline std::string* CategoryItem::_internal_mutable_resource_path() {
  ;
  return _impl_.resource_path_.Mutable( GetArenaForAllocation());
}
inline std::string* CategoryItem::release_resource_path() {
  // @@protoc_insertion_point(field_release:prpc.CategoryItem.resource_path)
  return _impl_.resource_path_.Release();
}
inline void CategoryItem::set_allocated_resource_path(std::string* value) {
  _impl_.resource_path_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.resource_path_.IsDefault()) {
          _impl_.resource_path_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:prpc.CategoryItem.resource_path)
}

// string poster_path = 6;
inline void CategoryItem::clear_poster_path() {
  _impl_.poster_path_.ClearToEmpty();
}
inline const std::string& CategoryItem::poster_path() const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.poster_path)
  return _internal_poster_path();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void CategoryItem::set_poster_path(Arg_&& arg,
                                                     Args_... args) {
  ;
  _impl_.poster_path_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.poster_path)
}
inline std::string* CategoryItem::mutable_poster_path() {
  std::string* _s = _internal_mutable_poster_path();
  // @@protoc_insertion_point(field_mutable:prpc.CategoryItem.poster_path)
  return _s;
}
inline const std::string& CategoryItem::_internal_poster_path() const {
  return _impl_.poster_path_.Get();
}
inline void CategoryItem::_internal_set_poster_path(const std::string& value) {
  ;


  _impl_.poster_path_.Set(value, GetArenaForAllocation());
}
inline std::string* CategoryItem::_internal_mutable_poster_path() {
  ;
  return _impl_.poster_path_.Mutable( GetArenaForAllocation());
}
inline std::string* CategoryItem::release_poster_path() {
  // @@protoc_insertion_point(field_release:prpc.CategoryItem.poster_path)
  return _impl_.poster_path_.Release();
}
inline void CategoryItem::set_allocated_poster_path(std::string* value) {
  _impl_.poster_path_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.poster_path_.IsDefault()) {
          _impl_.poster_path_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:prpc.CategoryItem.poster_path)
}

// string introduce = 7;
inline void CategoryItem::clear_introduce() {
  _impl_.introduce_.ClearToEmpty();
}
inline const std::string& CategoryItem::introduce() const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.introduce)
  return _internal_introduce();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void CategoryItem::set_introduce(Arg_&& arg,
                                                     Args_... args) {
  ;
  _impl_.introduce_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.introduce)
}
inline std::string* CategoryItem::mutable_introduce() {
  std::string* _s = _internal_mutable_introduce();
  // @@protoc_insertion_point(field_mutable:prpc.CategoryItem.introduce)
  return _s;
}
inline const std::string& CategoryItem::_internal_introduce() const {
  return _impl_.introduce_.Get();
}
inline void CategoryItem::_internal_set_introduce(const std::string& value) {
  ;


  _impl_.introduce_.Set(value, GetArenaForAllocation());
}
inline std::string* CategoryItem::_internal_mutable_introduce() {
  ;
  return _impl_.introduce_.Mutable( GetArenaForAllocation());
}
inline std::string* CategoryItem::release_introduce() {
  // @@protoc_insertion_point(field_release:prpc.CategoryItem.introduce)
  return _impl_.introduce_.Release();
}
inline void CategoryItem::set_allocated_introduce(std::string* value) {
  _impl_.introduce_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.introduce_.IsDefault()) {
          _impl_.introduce_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:prpc.CategoryItem.introduce)
}

// int64 parent_id = 8;
inline void CategoryItem::clear_parent_id() {
  _impl_.parent_id_ = ::int64_t{0};
}
inline ::int64_t CategoryItem::parent_id() const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.parent_id)
  return _internal_parent_id();
}
inline void CategoryItem::set_parent_id(::int64_t value) {
  _internal_set_parent_id(value);
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.parent_id)
}
inline ::int64_t CategoryItem::_internal_parent_id() const {
  return _impl_.parent_id_;
}
inline void CategoryItem::_internal_set_parent_id(::int64_t value) {
  ;
  _impl_.parent_id_ = value;
}

// repeated int64 sub_item_ids = 9;
inline int CategoryItem::_internal_sub_item_ids_size() const {
  return _impl_.sub_item_ids_.size();
}
inline int CategoryItem::sub_item_ids_size() const {
  return _internal_sub_item_ids_size();
}
inline void CategoryItem::clear_sub_item_ids() {
  _impl_.sub_item_ids_.Clear();
}
inline ::int64_t CategoryItem::sub_item_ids(int index) const {
  // @@protoc_insertion_point(field_get:prpc.CategoryItem.sub_item_ids)
  return _internal_sub_item_ids(index);
}
inline void CategoryItem::set_sub_item_ids(int index, ::int64_t value) {
  _impl_.sub_item_ids_.Set(index, value);
  // @@protoc_insertion_point(field_set:prpc.CategoryItem.sub_item_ids)
}
inline void CategoryItem::add_sub_item_ids(::int64_t value) {
  _internal_add_sub_item_ids(value);
  // @@protoc_insertion_point(field_add:prpc.CategoryItem.sub_item_ids)
}
inline const ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t>& CategoryItem::sub_item_ids() const {
  // @@protoc_insertion_point(field_list:prpc.CategoryItem.sub_item_ids)
  return _internal_sub_item_ids();
}
inline ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t>* CategoryItem::mutable_sub_item_ids() {
  // @@protoc_insertion_point(field_mutable_list:prpc.CategoryItem.sub_item_ids)
  return _internal_mutable_sub_item_ids();
}

inline ::int64_t CategoryItem::_internal_sub_item_ids(int index) const {
  return _impl_.sub_item_ids_.Get(index);
}
inline void CategoryItem::_internal_add_sub_item_ids(::int64_t value) { _impl_.sub_item_ids_.Add(value); }
inline const ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t>& CategoryItem::_internal_sub_item_ids() const {
  return _impl_.sub_item_ids_;
}
inline ::PROTOBUF_NAMESPACE_ID::RepeatedField<::int64_t>* CategoryItem::_internal_mutable_sub_item_ids() {
  return &_impl_.sub_item_ids_;
}

// -------------------------------------------------------------------

// SharedItem

// int64 item_id = 1;
inline void SharedItem::clear_item_id() {
  _impl_.item_id_ = ::int64_t{0};
}
inline ::int64_t SharedItem::item_id() const {
  // @@protoc_insertion_point(field_get:prpc.SharedItem.item_id)
  return _internal_item_id();
}
inline void SharedItem::set_item_id(::int64_t value) {
  _internal_set_item_id(value);
  // @@protoc_insertion_point(field_set:prpc.SharedItem.item_id)
}
inline ::int64_t SharedItem::_internal_item_id() const {
  return _impl_.item_id_;
}
inline void SharedItem::_internal_set_item_id(::int64_t value) {
  ;
  _impl_.item_id_ = value;
}

// string share_id = 2;
inline void SharedItem::clear_share_id() {
  _impl_.share_id_.ClearToEmpty();
}
inline const std::string& SharedItem::share_id() const {
  // @@protoc_insertion_point(field_get:prpc.SharedItem.share_id)
  return _internal_share_id();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void SharedItem::set_share_id(Arg_&& arg,
                                                     Args_... args) {
  ;
  _impl_.share_id_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:prpc.SharedItem.share_id)
}
inline std::string* SharedItem::mutable_share_id() {
  std::string* _s = _internal_mutable_share_id();
  // @@protoc_insertion_point(field_mutable:prpc.SharedItem.share_id)
  return _s;
}
inline const std::string& SharedItem::_internal_share_id() const {
  return _impl_.share_id_.Get();
}
inline void SharedItem::_internal_set_share_id(const std::string& value) {
  ;


  _impl_.share_id_.Set(value, GetArenaForAllocation());
}
inline std::string* SharedItem::_internal_mutable_share_id() {
  ;
  return _impl_.share_id_.Mutable( GetArenaForAllocation());
}
inline std::string* SharedItem::release_share_id() {
  // @@protoc_insertion_point(field_release:prpc.SharedItem.share_id)
  return _impl_.share_id_.Release();
}
inline void SharedItem::set_allocated_share_id(std::string* value) {
  _impl_.share_id_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.share_id_.IsDefault()) {
          _impl_.share_id_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:prpc.SharedItem.share_id)
}

#ifdef __GNUC__
#pragma GCC diagnostic pop
#endif  // __GNUC__

// @@protoc_insertion_point(namespace_scope)
}  // namespace prpc


PROTOBUF_NAMESPACE_OPEN

template <>
struct is_proto_enum<::prpc::CategoryItem_Type> : std::true_type {};
template <>
inline const EnumDescriptor* GetEnumDescriptor<::prpc::CategoryItem_Type>() {
  return ::prpc::CategoryItem_Type_descriptor();
}

PROTOBUF_NAMESPACE_CLOSE

// @@protoc_insertion_point(global_scope)

#include "google/protobuf/port_undef.inc"

#endif  // GOOGLE_PROTOBUF_INCLUDED_category_2eproto_2epb_2eh
