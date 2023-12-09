package bt

import (
	"errors"
	"pnas/prpc"
	"pnas/ptype"
)

var (
	ErrDownloaded = errors.New("downloaded")
)

type UserTorrents interface {
	NewTorrentByMagnet(magnetUri string) (*Torrent, error)
	HasTorrent(userId ptype.UserID, infoHash *InfoHash) bool
	GetTorrent(infoHash *InfoHash) (*Torrent, error)
	GetTorrents(userId ptype.UserID) []*Torrent

	SetTaskCallback(params *SetTaskCallbackParams)
	SetSessionCallback(userId ptype.UserID, sid ptype.SessionID, callback UserOnBtStatusCallback)

	Download(*DownloadParams) (*prpc.DownloadRespone, error)
	RemoveTorrent(*RemoveTorrentParams) (*prpc.RemoveTorrentRes, error)
	GetMagnetUri(*GetMagnetUriParams) (*prpc.GetMagnetUriRsp, error)
	GetPeerInfo(*GetPeerInfoParams) (*prpc.GetPeerInfoRsp, error)
	NewDownloadTask(*DownloadTaskParams) (*prpc.DownloadRespone, error)

	Close()
}

type UserOnBtStatusCallback func(error, *prpc.TorrentStatus)

type SetTaskCallbackParams struct {
	UserId    ptype.UserID
	TaskId    ptype.TaskId
	TorrentId ptype.TorrentID
	Callback  UserOnBtStatusCallback
}

type DownloadParams struct {
	UserId ptype.UserID
	Req    *prpc.DownloadRequest
}

type RemoveTorrentParams struct {
	UserId ptype.UserID
	Req    *prpc.RemoveTorrentReq
}

type GetMagnetUriParams struct {
	Req *prpc.GetMagnetUriReq
}

type DownloadTaskParams struct {
	UserId   ptype.UserID
	TaskId   ptype.TaskId
	Req      *prpc.DownloadRequest
	Callback UserOnBtStatusCallback
}

type GetPeerInfoParams struct {
	UserId ptype.UserID
	Req    *prpc.GetPeerInfoReq
}
