package bt

import (
	"pnas/prpc"
	"pnas/ptype"
)

type UserTorrents interface {
	HasTorrent(userId ptype.UserID, infoHash *InfoHash) bool
	GetTorrent(infoHash *InfoHash) (*Torrent, error)

	SetCallback(userId ptype.UserID, sid ptype.SessionID, callback UserOnBtStatusCallback)

	Download(*DownloadParams) (*prpc.DownloadRespone, error)
	RemoveTorrent(*RemoveTorrentParams) (*prpc.RemoveTorrentRes, error)
	GetMagnetUri(*GetMagnetUriParams) (*prpc.GetMagnetUriRsp, error)

	GetTorrents(userId ptype.UserID) []*Torrent

	GetBtClient() *BtClient

	Close()
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
	UserId ptype.UserID
	Req    *prpc.GetMagnetUriReq
}

type UserOnBtStatusCallback func(*prpc.TorrentStatus)
