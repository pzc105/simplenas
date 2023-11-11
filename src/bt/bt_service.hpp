#pragma once

#include "libtorrent/entry.hpp"
#include "libtorrent/bencode.hpp"
#include "libtorrent/session.hpp"
#include "libtorrent/torrent_info.hpp"
#include "libtorrent/magnet_uri.hpp"
#include "libtorrent/alert_types.hpp"
#include "prpc/bt.grpc.pb.h"
#include <queue>
#include <mutex>

namespace prpc
{
  class bt_service;
  class pusher_manager;

  struct pusher_base
  {
    struct tag
    {
      pusher_base *_owner;
    };

    pusher_base(pusher_manager *owner) : _owner(owner), _new_tag{this}, _read_tag{this}, _write_tag{this} {}
    virtual void completed(tag *t, bool ok) = 0;
    virtual void done() = 0;

    pusher_manager *_owner;
    tag _new_tag;
    tag _read_tag;
    tag _write_tag;
  };

  template <typename request_type, typename respone_type>
  class pusher : public pusher_base
  {
  public:
    using push_class = pusher;

    pusher(pusher_manager *owner, bt_service *ser) : pusher_base(owner), _stream(&_context), _ser(ser) {}
    void read();
    void write(respone_type &&r);

    void completed(tag *t, bool ok) override;

    std::mutex _mtx;
    uint32_t _read_count{0};
    uint32_t _write_count{0};
    std::queue<respone_type> _write_queue;

    request_type _req;
    grpc::ServerContext _context;
    grpc::ServerAsyncReaderWriter<respone_type, request_type> _stream;

    bt_service *_ser;
  };

  class bt_status_pusher : public pusher<BtStatusRequest, BtStatusRespone>
  {
  public:
    bt_status_pusher(pusher_manager *owner, bt_service *ser);
    void push(std::vector<lt::torrent_status> const &sts);
    void completed(tag *t, bool ok) override;
    void done() override;
  };

  class bt_info_pusher : public pusher<TorrentInfoReq, TorrentInfoRes>
  {
  public:
    bt_info_pusher(pusher_manager *owner, bt_service *ser);
    void push(lt::add_torrent_params const &params);
    void completed(tag *t, bool ok) override;
    void done() override;
  };

  class bt_filecompleted_pusher : public pusher<FileCompletedReq, FileCompletedRes>
  {
  public:
    bt_filecompleted_pusher(pusher_manager *owner, bt_service *ser);
    void push(lt::file_completed_alert const &params);
    void completed(tag *t, bool ok) override;
    void done() override;
  };

  using status_pusher_ptr = std::unique_ptr<bt_status_pusher>;
  using filecompleted_pusher_ptr = std::unique_ptr<bt_filecompleted_pusher>;

  class pusher_manager
  {
    friend class bt_status_pusher;
    friend class bt_info_pusher;
    friend class bt_filecompleted_pusher;

  public:
    pusher_manager(bt_service *ser);
    void start();
    void push_bt_status(std::vector<lt::torrent_status> const &sts);
    void push_bt_filecompleted(lt::file_completed_alert const &params);

  private:
    void accepted_status_pusher(pusher_base *pusher);
    void remove_status_pusher(pusher_base *pusher);

    void accepted_filecompleted_pusher(pusher_base *pusher);
    void remove_filecompleted_pusher(pusher_base *pusher);

  private:
    bt_service *_ser;

    std::vector<status_pusher_ptr> _st_pushers;
    status_pusher_ptr _st_pusher;

    std::vector<filecompleted_pusher_ptr> _filecompleted_pushers;
    filecompleted_pusher_ptr _filecompleted_pusher;
  };

  using _bt_service = BtService::WithAsyncMethod_OnFileCompleted<
      BtService::WithAsyncMethod_OnBtStatus<BtService::Service>>;

  class bt_service : public _bt_service
  {
  public:
    bt_service();
    ~bt_service();
    void own_completion_queue(std::unique_ptr<grpc::ServerCompletionQueue> &&cq);
    void run();
    void push_bt();
    grpc::ServerCompletionQueue *get_cq() const { return _cq.get(); }

  private:
    ::grpc::Status InitedSession(::grpc::ServerContext* context, const ::prpc::InitedSessionReq* request, ::prpc::InitedSessionRsp* response) override;
    ::grpc::Status InitSession(::grpc::ServerContext *context, const ::prpc::InitSessionReq *request, ::prpc::InitSessionRsp *response) override;
    ::grpc::Status Parse(::grpc::ServerContext *context, const ::prpc::DownloadRequest *request, ::prpc::DownloadRespone *response) override;
    ::grpc::Status Download(::grpc::ServerContext *context, const ::prpc::DownloadRequest *request, ::prpc::DownloadRespone *response) override;
    ::grpc::Status RemoveTorrent(::grpc::ServerContext *context, const ::prpc::RemoveTorrentReq *request, ::prpc::RemoveTorrentRes *response) override;
    ::grpc::Status GetMagnetUri(::grpc::ServerContext *context, const ::prpc::GetMagnetUriReq *request, ::prpc::GetMagnetUriRsp *response) override;
    ::grpc::Status GetResumeData(::grpc::ServerContext *context, const ::prpc::GetResumeDataReq *request, ::prpc::GetResumeDataRsp *response) override;
    ::grpc::Status GetTorrentInfo(::grpc::ServerContext *context, const ::prpc::GetTorrentInfoReq *request, ::prpc::GetTorrentInfoRsp *response) override;
    ::grpc::Status GetBtStatus(::grpc::ServerContext *context, const ::prpc::GetBtStatusReq *request, ::prpc::GetBtStatusRsp *response) override;
    ::grpc::Status GetSessionParams(::grpc::ServerContext *context, const ::prpc::GetSessionParamsReq *request, ::prpc::GetSessionParamsRsp *response) override;

  private:
    std::unique_ptr<lt::session> _ses;
    lt::time_point _last_push_time;

    std::unique_ptr<grpc::ServerCompletionQueue> _cq;

    pusher_manager _pusher_manager;
  };

  template <typename request_type, typename respone_type>
  void pusher<request_type, respone_type>::read()
  {
    std::lock_guard<std::mutex> lk(_mtx);
    _stream.Read(&_req, &_read_tag);
    _read_count++;
  }

  template <typename request_type, typename respone_type>
  void pusher<request_type, respone_type>::write(respone_type &&r)
  {
    std::lock_guard<std::mutex> lk(_mtx);
    if (_write_count > 0)
    {
      _write_queue.emplace(std::forward<respone_type>(r));
    }
    else
    {
      _stream.Write(r, &_write_tag);
    }
    _write_count++;
  }

  template <typename request_type, typename respone_type>
  void pusher<request_type, respone_type>::completed(pusher::tag *t, bool ok)
  {
    bool r = false;

    if (t == &_new_tag && ok)
    {
      read();
    }

    {
      std::lock_guard<std::mutex> lk(_mtx);
      if (t == &_write_tag)
      {
        _write_count -= 1;
        if (!ok)
        {
          _write_count = 0;
          _write_queue = std::queue<respone_type>();
        }
        else if (_write_count > 0)
        {
          respone_type rs = std::move(_write_queue.front());
          _write_queue.pop();
          _stream.Write(rs, &_write_tag);
        }
      }
      else if (t == &_read_tag)
      {
        _read_count--;
        if (ok)
        {
          read();
        }
      }
      else if (t == &_new_tag)
      {
      }
      else
      {
        assert(false);
      }

      if (!ok && _write_count == 0 && _read_count == 0)
      {
        r = true;
      }
    }

    if (r)
    {
      done();
    }
  }
}