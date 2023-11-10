package bt

import (
	"pnas/prpc"
	"pnas/ptype"
)

type UserTorrents interface {
	HasTorrent(userId ptype.UserID, infoHash InfoHash) bool
	GetTorrent(infoHash InfoHash) (*Torrent, error)

	Download(*DownloadParams) (*prpc.DownloadRespone, error)
	RemoveTorrent(*RemoveTorrentParams) (*prpc.RemoveTorrentRes, error)
	GetMagnetUri(*GetMagnetUriParams) (*prpc.GetMagnetUriRsp, error)

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
