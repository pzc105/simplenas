// Generated by the protocol buffer compiler.  DO NOT EDIT!
// source: danmaku.proto

#ifndef GOOGLE_PROTOBUF_INCLUDED_danmaku_2eproto_2epb_2eh
#define GOOGLE_PROTOBUF_INCLUDED_danmaku_2eproto_2epb_2eh

#include <limits>
#include <string>
#include <type_traits>

#include "google/protobuf/port_def.inc"
#if PROTOBUF_VERSION < 4024000
#error "This file was generated by a newer version of protoc which is"
#error "incompatible with your Protocol Buffer headers. Please update"
#error "your headers."
#endif  // PROTOBUF_VERSION

#if 4024003 < PROTOBUF_MIN_PROTOC_VERSION
#error "This file was generated by an older version of protoc which is"
#error "incompatible with your Protocol Buffer headers. Please"
#error "regenerate this file with a newer version of protoc."
#endif  // PROTOBUF_MIN_PROTOC_VERSION
#include "google/protobuf/port_undef.inc"
#include "google/protobuf/io/coded_stream.h"
#include "google/protobuf/arena.h"
#include "google/protobuf/arenastring.h"
#include "google/protobuf/generated_message_tctable_decl.h"
#include "google/protobuf/generated_message_util.h"
#include "google/protobuf/metadata_lite.h"
#include "google/protobuf/generated_message_reflection.h"
#include "google/protobuf/message.h"
#include "google/protobuf/repeated_field.h"  // IWYU pragma: export
#include "google/protobuf/extension_set.h"  // IWYU pragma: export
#include "google/protobuf/unknown_field_set.h"
// @@protoc_insertion_point(includes)

// Must be included last.
#include "google/protobuf/port_def.inc"

#define PROTOBUF_INTERNAL_EXPORT_danmaku_2eproto

namespace google {
namespace protobuf {
namespace internal {
class AnyMetadata;
}  // namespace internal
}  // namespace protobuf
}  // namespace google

// Internal implementation detail -- do not use these members.
struct TableStruct_danmaku_2eproto {
  static const ::uint32_t offsets[];
};
extern const ::google::protobuf::internal::DescriptorTable
    descriptor_table_danmaku_2eproto;
namespace prpc {
class Danmaku;
struct DanmakuDefaultTypeInternal;
extern DanmakuDefaultTypeInternal _Danmaku_default_instance_;
}  // namespace prpc
namespace google {
namespace protobuf {
}  // namespace protobuf
}  // namespace google

namespace prpc {

// ===================================================================


// -------------------------------------------------------------------

class Danmaku final :
    public ::google::protobuf::Message /* @@protoc_insertion_point(class_definition:prpc.Danmaku) */ {
 public:
  inline Danmaku() : Danmaku(nullptr) {}
  ~Danmaku() override;
  template<typename = void>
  explicit PROTOBUF_CONSTEXPR Danmaku(::google::protobuf::internal::ConstantInitialized);

  Danmaku(const Danmaku& from);
  Danmaku(Danmaku&& from) noexcept
    : Danmaku() {
    *this = ::std::move(from);
  }

  inline Danmaku& operator=(const Danmaku& from) {
    CopyFrom(from);
    return *this;
  }
  inline Danmaku& operator=(Danmaku&& from) noexcept {
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

  inline const ::google::protobuf::UnknownFieldSet& unknown_fields() const {
    return _internal_metadata_.unknown_fields<::google::protobuf::UnknownFieldSet>(::google::protobuf::UnknownFieldSet::default_instance);
  }
  inline ::google::protobuf::UnknownFieldSet* mutable_unknown_fields() {
    return _internal_metadata_.mutable_unknown_fields<::google::protobuf::UnknownFieldSet>();
  }

  static const ::google::protobuf::Descriptor* descriptor() {
    return GetDescriptor();
  }
  static const ::google::protobuf::Descriptor* GetDescriptor() {
    return default_instance().GetMetadata().descriptor;
  }
  static const ::google::protobuf::Reflection* GetReflection() {
    return default_instance().GetMetadata().reflection;
  }
  static const Danmaku& default_instance() {
    return *internal_default_instance();
  }
  static inline const Danmaku* internal_default_instance() {
    return reinterpret_cast<const Danmaku*>(
               &_Danmaku_default_instance_);
  }
  static constexpr int kIndexInFileMessages =
    0;

  friend void swap(Danmaku& a, Danmaku& b) {
    a.Swap(&b);
  }
  inline void Swap(Danmaku* other) {
    if (other == this) return;
  #ifdef PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() != nullptr &&
        GetOwningArena() == other->GetOwningArena()) {
   #else  // PROTOBUF_FORCE_COPY_IN_SWAP
    if (GetOwningArena() == other->GetOwningArena()) {
  #endif  // !PROTOBUF_FORCE_COPY_IN_SWAP
      InternalSwap(other);
    } else {
      ::google::protobuf::internal::GenericSwap(this, other);
    }
  }
  void UnsafeArenaSwap(Danmaku* other) {
    if (other == this) return;
    ABSL_DCHECK(GetOwningArena() == other->GetOwningArena());
    InternalSwap(other);
  }

  // implements Message ----------------------------------------------

  Danmaku* New(::google::protobuf::Arena* arena = nullptr) const final {
    return CreateMaybeMessage<Danmaku>(arena);
  }
  using ::google::protobuf::Message::CopyFrom;
  void CopyFrom(const Danmaku& from);
  using ::google::protobuf::Message::MergeFrom;
  void MergeFrom( const Danmaku& from) {
    Danmaku::MergeImpl(*this, from);
  }
  private:
  static void MergeImpl(::google::protobuf::Message& to_msg, const ::google::protobuf::Message& from_msg);
  public:
  PROTOBUF_ATTRIBUTE_REINITIALIZES void Clear() final;
  bool IsInitialized() const final;

  ::size_t ByteSizeLong() const final;
  const char* _InternalParse(const char* ptr, ::google::protobuf::internal::ParseContext* ctx) final;
  ::uint8_t* _InternalSerialize(
      ::uint8_t* target, ::google::protobuf::io::EpsCopyOutputStream* stream) const final;
  int GetCachedSize() const final { return _impl_._cached_size_.Get(); }

  private:
  void SharedCtor(::google::protobuf::Arena* arena);
  void SharedDtor();
  void SetCachedSize(int size) const final;
  void InternalSwap(Danmaku* other);

  private:
  friend class ::google::protobuf::internal::AnyMetadata;
  static ::absl::string_view FullMessageName() {
    return "prpc.Danmaku";
  }
  protected:
  explicit Danmaku(::google::protobuf::Arena* arena);
  public:

  static const ClassData _class_data_;
  const ::google::protobuf::Message::ClassData*GetClassData() const final;

  ::google::protobuf::Metadata GetMetadata() const final;

  // nested types ----------------------------------------------------

  // accessors -------------------------------------------------------

  enum : int {
    kIdFieldNumber = 1,
    kUserNameFieldNumber = 3,
    kTextFieldNumber = 5,
    kUserIdFieldNumber = 2,
    kSTimeFieldNumber = 4,
    kTypeFieldNumber = 6,
    kColorFieldNumber = 7,
    kDTimeFieldNumber = 8,
  };
  // string id = 1;
  void clear_id() ;
  const std::string& id() const;
  template <typename Arg_ = const std::string&, typename... Args_>
  void set_id(Arg_&& arg, Args_... args);
  std::string* mutable_id();
  PROTOBUF_NODISCARD std::string* release_id();
  void set_allocated_id(std::string* ptr);

  private:
  const std::string& _internal_id() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_id(
      const std::string& value);
  std::string* _internal_mutable_id();

  public:
  // string user_name = 3;
  void clear_user_name() ;
  const std::string& user_name() const;
  template <typename Arg_ = const std::string&, typename... Args_>
  void set_user_name(Arg_&& arg, Args_... args);
  std::string* mutable_user_name();
  PROTOBUF_NODISCARD std::string* release_user_name();
  void set_allocated_user_name(std::string* ptr);

  private:
  const std::string& _internal_user_name() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_user_name(
      const std::string& value);
  std::string* _internal_mutable_user_name();

  public:
  // string text = 5;
  void clear_text() ;
  const std::string& text() const;
  template <typename Arg_ = const std::string&, typename... Args_>
  void set_text(Arg_&& arg, Args_... args);
  std::string* mutable_text();
  PROTOBUF_NODISCARD std::string* release_text();
  void set_allocated_text(std::string* ptr);

  private:
  const std::string& _internal_text() const;
  inline PROTOBUF_ALWAYS_INLINE void _internal_set_text(
      const std::string& value);
  std::string* _internal_mutable_text();

  public:
  // int64 user_id = 2;
  void clear_user_id() ;
  ::int64_t user_id() const;
  void set_user_id(::int64_t value);

  private:
  ::int64_t _internal_user_id() const;
  void _internal_set_user_id(::int64_t value);

  public:
  // int64 s_time = 4;
  void clear_s_time() ;
  ::int64_t s_time() const;
  void set_s_time(::int64_t value);

  private:
  ::int64_t _internal_s_time() const;
  void _internal_set_s_time(::int64_t value);

  public:
  // int32 type = 6;
  void clear_type() ;
  ::int32_t type() const;
  void set_type(::int32_t value);

  private:
  ::int32_t _internal_type() const;
  void _internal_set_type(::int32_t value);

  public:
  // int32 color = 7;
  void clear_color() ;
  ::int32_t color() const;
  void set_color(::int32_t value);

  private:
  ::int32_t _internal_color() const;
  void _internal_set_color(::int32_t value);

  public:
  // double d_time = 8;
  void clear_d_time() ;
  double d_time() const;
  void set_d_time(double value);

  private:
  double _internal_d_time() const;
  void _internal_set_d_time(double value);

  public:
  // @@protoc_insertion_point(class_scope:prpc.Danmaku)
 private:
  class _Internal;

  friend class ::google::protobuf::internal::TcParser;
  static const ::google::protobuf::internal::TcParseTable<3, 8, 0, 44, 2> _table_;
  template <typename T> friend class ::google::protobuf::Arena::InternalHelper;
  typedef void InternalArenaConstructable_;
  typedef void DestructorSkippable_;
  struct Impl_ {
    ::google::protobuf::internal::ArenaStringPtr id_;
    ::google::protobuf::internal::ArenaStringPtr user_name_;
    ::google::protobuf::internal::ArenaStringPtr text_;
    ::int64_t user_id_;
    ::int64_t s_time_;
    ::int32_t type_;
    ::int32_t color_;
    double d_time_;
    mutable ::google::protobuf::internal::CachedSize _cached_size_;
    PROTOBUF_TSAN_DECLARE_MEMBER
  };
  union { Impl_ _impl_; };
  friend struct ::TableStruct_danmaku_2eproto;
};

// ===================================================================




// ===================================================================


#ifdef __GNUC__
#pragma GCC diagnostic push
#pragma GCC diagnostic ignored "-Wstrict-aliasing"
#endif  // __GNUC__
// -------------------------------------------------------------------

// Danmaku

// string id = 1;
inline void Danmaku::clear_id() {
  _impl_.id_.ClearToEmpty();
}
inline const std::string& Danmaku::id() const {
  // @@protoc_insertion_point(field_get:prpc.Danmaku.id)
  return _internal_id();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void Danmaku::set_id(Arg_&& arg,
                                                     Args_... args) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.id_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:prpc.Danmaku.id)
}
inline std::string* Danmaku::mutable_id() {
  std::string* _s = _internal_mutable_id();
  // @@protoc_insertion_point(field_mutable:prpc.Danmaku.id)
  return _s;
}
inline const std::string& Danmaku::_internal_id() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.id_.Get();
}
inline void Danmaku::_internal_set_id(const std::string& value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.id_.Set(value, GetArenaForAllocation());
}
inline std::string* Danmaku::_internal_mutable_id() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  return _impl_.id_.Mutable( GetArenaForAllocation());
}
inline std::string* Danmaku::release_id() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  // @@protoc_insertion_point(field_release:prpc.Danmaku.id)
  return _impl_.id_.Release();
}
inline void Danmaku::set_allocated_id(std::string* value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  _impl_.id_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.id_.IsDefault()) {
          _impl_.id_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:prpc.Danmaku.id)
}

// int64 user_id = 2;
inline void Danmaku::clear_user_id() {
  _impl_.user_id_ = ::int64_t{0};
}
inline ::int64_t Danmaku::user_id() const {
  // @@protoc_insertion_point(field_get:prpc.Danmaku.user_id)
  return _internal_user_id();
}
inline void Danmaku::set_user_id(::int64_t value) {
  _internal_set_user_id(value);
  // @@protoc_insertion_point(field_set:prpc.Danmaku.user_id)
}
inline ::int64_t Danmaku::_internal_user_id() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.user_id_;
}
inline void Danmaku::_internal_set_user_id(::int64_t value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.user_id_ = value;
}

// string user_name = 3;
inline void Danmaku::clear_user_name() {
  _impl_.user_name_.ClearToEmpty();
}
inline const std::string& Danmaku::user_name() const {
  // @@protoc_insertion_point(field_get:prpc.Danmaku.user_name)
  return _internal_user_name();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void Danmaku::set_user_name(Arg_&& arg,
                                                     Args_... args) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.user_name_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:prpc.Danmaku.user_name)
}
inline std::string* Danmaku::mutable_user_name() {
  std::string* _s = _internal_mutable_user_name();
  // @@protoc_insertion_point(field_mutable:prpc.Danmaku.user_name)
  return _s;
}
inline const std::string& Danmaku::_internal_user_name() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.user_name_.Get();
}
inline void Danmaku::_internal_set_user_name(const std::string& value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.user_name_.Set(value, GetArenaForAllocation());
}
inline std::string* Danmaku::_internal_mutable_user_name() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  return _impl_.user_name_.Mutable( GetArenaForAllocation());
}
inline std::string* Danmaku::release_user_name() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  // @@protoc_insertion_point(field_release:prpc.Danmaku.user_name)
  return _impl_.user_name_.Release();
}
inline void Danmaku::set_allocated_user_name(std::string* value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  _impl_.user_name_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.user_name_.IsDefault()) {
          _impl_.user_name_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:prpc.Danmaku.user_name)
}

// int64 s_time = 4;
inline void Danmaku::clear_s_time() {
  _impl_.s_time_ = ::int64_t{0};
}
inline ::int64_t Danmaku::s_time() const {
  // @@protoc_insertion_point(field_get:prpc.Danmaku.s_time)
  return _internal_s_time();
}
inline void Danmaku::set_s_time(::int64_t value) {
  _internal_set_s_time(value);
  // @@protoc_insertion_point(field_set:prpc.Danmaku.s_time)
}
inline ::int64_t Danmaku::_internal_s_time() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.s_time_;
}
inline void Danmaku::_internal_set_s_time(::int64_t value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.s_time_ = value;
}

// string text = 5;
inline void Danmaku::clear_text() {
  _impl_.text_.ClearToEmpty();
}
inline const std::string& Danmaku::text() const {
  // @@protoc_insertion_point(field_get:prpc.Danmaku.text)
  return _internal_text();
}
template <typename Arg_, typename... Args_>
inline PROTOBUF_ALWAYS_INLINE void Danmaku::set_text(Arg_&& arg,
                                                     Args_... args) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.text_.Set(static_cast<Arg_&&>(arg), args..., GetArenaForAllocation());
  // @@protoc_insertion_point(field_set:prpc.Danmaku.text)
}
inline std::string* Danmaku::mutable_text() {
  std::string* _s = _internal_mutable_text();
  // @@protoc_insertion_point(field_mutable:prpc.Danmaku.text)
  return _s;
}
inline const std::string& Danmaku::_internal_text() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.text_.Get();
}
inline void Danmaku::_internal_set_text(const std::string& value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.text_.Set(value, GetArenaForAllocation());
}
inline std::string* Danmaku::_internal_mutable_text() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  return _impl_.text_.Mutable( GetArenaForAllocation());
}
inline std::string* Danmaku::release_text() {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  // @@protoc_insertion_point(field_release:prpc.Danmaku.text)
  return _impl_.text_.Release();
}
inline void Danmaku::set_allocated_text(std::string* value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  _impl_.text_.SetAllocated(value, GetArenaForAllocation());
  #ifdef PROTOBUF_FORCE_COPY_DEFAULT_STRING
        if (_impl_.text_.IsDefault()) {
          _impl_.text_.Set("", GetArenaForAllocation());
        }
  #endif  // PROTOBUF_FORCE_COPY_DEFAULT_STRING
  // @@protoc_insertion_point(field_set_allocated:prpc.Danmaku.text)
}

// int32 type = 6;
inline void Danmaku::clear_type() {
  _impl_.type_ = 0;
}
inline ::int32_t Danmaku::type() const {
  // @@protoc_insertion_point(field_get:prpc.Danmaku.type)
  return _internal_type();
}
inline void Danmaku::set_type(::int32_t value) {
  _internal_set_type(value);
  // @@protoc_insertion_point(field_set:prpc.Danmaku.type)
}
inline ::int32_t Danmaku::_internal_type() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.type_;
}
inline void Danmaku::_internal_set_type(::int32_t value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.type_ = value;
}

// int32 color = 7;
inline void Danmaku::clear_color() {
  _impl_.color_ = 0;
}
inline ::int32_t Danmaku::color() const {
  // @@protoc_insertion_point(field_get:prpc.Danmaku.color)
  return _internal_color();
}
inline void Danmaku::set_color(::int32_t value) {
  _internal_set_color(value);
  // @@protoc_insertion_point(field_set:prpc.Danmaku.color)
}
inline ::int32_t Danmaku::_internal_color() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.color_;
}
inline void Danmaku::_internal_set_color(::int32_t value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.color_ = value;
}

// double d_time = 8;
inline void Danmaku::clear_d_time() {
  _impl_.d_time_ = 0;
}
inline double Danmaku::d_time() const {
  // @@protoc_insertion_point(field_get:prpc.Danmaku.d_time)
  return _internal_d_time();
}
inline void Danmaku::set_d_time(double value) {
  _internal_set_d_time(value);
  // @@protoc_insertion_point(field_set:prpc.Danmaku.d_time)
}
inline double Danmaku::_internal_d_time() const {
  PROTOBUF_TSAN_READ(&_impl_._tsan_detect_race);
  return _impl_.d_time_;
}
inline void Danmaku::_internal_set_d_time(double value) {
  PROTOBUF_TSAN_WRITE(&_impl_._tsan_detect_race);
  ;
  _impl_.d_time_ = value;
}

#ifdef __GNUC__
#pragma GCC diagnostic pop
#endif  // __GNUC__

// @@protoc_insertion_point(namespace_scope)
}  // namespace prpc


// @@protoc_insertion_point(global_scope)

#include "google/protobuf/port_undef.inc"

#endif  // GOOGLE_PROTOBUF_INCLUDED_danmaku_2eproto_2epb_2eh
