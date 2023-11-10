package bt

import (
	"context"
	"encoding/hex"
	"os"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/setting"
	"pnas/video"
	"sync"
	"time"

	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	AdminId = 1
)

type UpdateTorrentParams struct {
	Base       *TorrentBase
	State      prpc.BtStateEnum
	FileNames  []File
	ResumeData []byte
}

type FileCompleted struct {
	InfoHash  InfoHash
	FileIndex int32
}

type userData struct {
	mtx      sync.Mutex
	userId   ptype.UserID
	torrents map[ptype.TorrentID]bool
}

func (ud *userData) hasTorrent(id ptype.TorrentID) bool {
	ud.mtx.Lock()
	_, ok := ud.torrents[id]
	ud.mtx.Unlock()
	return ok
}

func (ud *userData) addTorrent(id ptype.TorrentID) error {
	if ud.hasTorrent(id) {
		return errors.New(("duplicated"))
	}
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	sql := "insert into user_torrent (user_id, torrent_id) values(?, ?)"
	_, err := db.Exec(sql, ud.userId, id)
	return err
}

func (ud *userData) removeTorrent(id ptype.TorrentID) error {
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	_, ok := ud.torrents[id]
	if !ok {
		return errors.New("not found torrent")
	}
	sql := "delete from user_torrent where user_id=? and torrent_id=?"
	_, err := db.Exec(sql, ud.userId, id)
	return err
}

type UserTorrentsImpl struct {
	UserTorrents
	mtx      sync.Mutex
	torrents map[InfoHash]*Torrent
	users    map[ptype.UserID]*userData

	bt BtClient
}

func (ut *UserTorrentsImpl) Init() {
	ut.torrents = make(map[InfoHash]*Torrent)

	ut.bt.Init(WithOnStatus(ut.handleBtStatus),
		WithOnConnect(ut.handleBtClientConnected),
		WithOnFileCompleted(ut.handleBtFileCompleted))
}

func (ut *UserTorrentsImpl) Close() {
	ut.bt.Close()
}

func (ut *UserTorrentsImpl) handleBtClientConnected() {
	log.Info("connected to bt service")
	resumeData := LoadDownloadingTorrent()
	for _, resume := range resumeData {
		var req prpc.DownloadRequest
		req.Type = prpc.DownloadRequest_Resume
		req.Content = resume
		req.SavePath = setting.GS().Bt.SavePath
		_, err := ut.bt.Download(context.Background(), &req)
		if err != nil {
			log.Warnf("[bt] failed to download, err: %v", err)
		}
	}
}

func (ut *UserTorrentsImpl) handleBtStatus(sr *prpc.StatusRespone) {

}

func (ut *UserTorrentsImpl) handleBtFileCompleted(fs *prpc.FileCompletedRes) {
	lfc := &FileCompleted{
		InfoHash:  TranInfoHash(fs.InfoHash),
		FileIndex: fs.FileIndex,
	}
	go ut.BtFileStateComplete(lfc)
}

func updateBtFileType(st *Torrent, fileIndex int, absFileName string) {
	if video.IsVideo(absFileName) {
		st.UpdateFileType(fileIndex, FileVideoType)
		meta, _ := video.GetMetadata(absFileName)
		st.UpdateVideoFileMeta(fileIndex, meta)
	}
	if video.IsSubTitle(absFileName) {
		st.UpdateFileType(fileIndex, FileSubtitleType)
	}
	if video.IsAudio(absFileName) {
		st.UpdateFileType(fileIndex, FileAudioType)
	}
}

func (ut *UserTorrentsImpl) UpdateTorrent(params *UpdateTorrentParams) {
	srcState := prpc.BtStateEnum_unknown
	ut.mtx.Lock()
	st, ok := ut.torrents[params.Base.InfoHash]
	if !ok {
		st = NewTorrent(params.Base)
		ut.torrents[params.Base.InfoHash] = st
	}
	ut.mtx.Unlock()
	srcState = st.GetState()

	lastUpdateTime := st.GetUpdateTime()

	if srcState == prpc.BtStateEnum_seeding && time.Since(lastUpdateTime) < time.Hour*1 {
		return
	}

	st.Update(params.Base, &params.State, params.FileNames, params.ResumeData)

	baseInfo := st.GetBaseInfo()
	sql := "update torrent set name=?, state=?, total_size=?, piece_length=?, num_pieces=?, resume_data=? where info_hash=? and version=?"
	_, err := db.Exec(sql, baseInfo.Name, params.State,
		baseInfo.TotalSize, baseInfo.PieceLength, baseInfo.NumPieces,
		st.GetResumeData(), baseInfo.InfoHash.Hash, baseInfo.InfoHash.Version)
	if err != nil {
		log.Warnf("[bt] failed to update torrent %s err: %v", hex.EncodeToString([]byte(params.Base.InfoHash.Hash)), err)
	}

	if params.State != prpc.BtStateEnum_seeding {
		return
	}

	files := st.GetFiles()
	var fileIndexes []int32
	for i := range files {
		fileName := files[i].Name

		if files[i].FileType == FileUnknownType {
			absFileName := setting.GS().Bt.SavePath + "/" + fileName
			updateBtFileType(st, i, absFileName)
			ft, _ := st.GetFileType(i)
			log.Debugf("[bt] torrent:%s file: %s type: %d", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, ft)
			if log.EnabledDebug() {
				_, err := os.Stat(absFileName)
				if err != nil {
					log.Debugf("[bt] torrent: %s file: %s error: %v", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, err)
				}
				meta, _ := video.GetMetadata(absFileName)
				if meta != nil {
					log.Debugf("[bt] torrent: %s file: %s format: %s", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, meta.Format.FormatName)
					for _, s := range meta.Streams {
						log.Debugf("[bt] torrent: %s file: %s s: %d codeType: %s", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, s.Index, s.CodecType)
					}
				}
			}
		}

		v, err := video.GetVideoByFileName(fileName)
		if err == nil && v.HlsCreated {
			continue
		}

		if ft, _ := st.GetFileType(i); (ft | FileVideoType) != 0 {
			fileIndexes = append(fileIndexes, int32(i))
		}
	}
	if len(fileIndexes) == 0 {
		return
	}
}

func (ut *UserTorrentsImpl) BtFileStateComplete(fs *FileCompleted) {
	ut.mtx.Lock()
	t, ok := ut.torrents[fs.InfoHash]
	ut.mtx.Unlock()
	if !ok {
		return
	}
	tfs, err := t.GetFileState(int(fs.FileIndex))
	if err != nil {
		log.Warnf("[bt] failed to update torrent: %s fileindex: %d state err: %v", hex.EncodeToString([]byte(fs.InfoHash.Hash)), fs.FileIndex, err)
		return
	}
	if tfs == prpc.BtFile_completed {
		return
	}

	t.UpdateFileState(int(fs.FileIndex), prpc.BtFile_completed)

	baseInfo := t.GetBaseInfo()
	files := t.GetFiles()

	log.Infof("[bt] torrent: %s file: %s completed", hex.EncodeToString([]byte(fs.InfoHash.Hash)), files[fs.FileIndex].Name)

	fileName := files[fs.FileIndex].Name
	absFileName := baseInfo.SavePath + "/" + fileName
	updateBtFileType(t, int(fs.FileIndex), absFileName)
	ft, _ := t.GetFileType(int(fs.FileIndex))
	log.Debugf("[bt] torrent:%s file: %s type: %d", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, ft)

	if log.EnabledDebug() {
		_, err := os.Stat(absFileName)
		if err != nil {
			log.Debugf("[bt] torrent: %s file: %s error: %v", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, err)
		}
		meta, _ := video.GetMetadata(absFileName)
		if meta != nil {
			log.Debugf("[bt] torrent: %s file: %s format: %s", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, meta.Format.FormatName)
			for _, s := range meta.Streams {
				log.Debugf("[bt] torrent: %s file: %s s: %d codeType: %s", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, s.Index, s.CodecType)
			}
		}
	}
}

func (ut *UserTorrentsImpl) HasTorrent(userId ptype.UserID, infoHash InfoHash) bool {
	if userId == AdminId {
		return true
	}
	ut.mtx.Lock()
	t, ok1 := ut.torrents[infoHash]
	u, ok2 := ut.users[userId]
	ut.mtx.Unlock()
	if !ok1 || !ok2 {
		return false
	}
	return u.hasTorrent(t.base.Id)
}

func (ut *UserTorrentsImpl) GetTorrent(infoHash InfoHash) (*Torrent, error) {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	t, ok := ut.torrents[infoHash]
	if !ok {
		return nil, errors.New("not found bt")
	}
	return t, nil
}

func (ut *UserTorrentsImpl) Download(params *DownloadParams) (*prpc.DownloadRespone, error) {
	req := params.Req
	res, err := ut.bt.Parse(context.Background(), req)
	if err != nil {
		return nil, err
	}
	resumeData, err := LoadTorrentResumeData(TranInfoHash(res.InfoHash))
	if err == nil {
		req.Type = prpc.DownloadRequest_Resume
		req.Content = resumeData
	}

	req.SavePath = setting.GS().Bt.SavePath
	res, err = ut.bt.Download(context.Background(), req)
	if err == nil {
		return res, nil
	} else {
		return res, status.Error(codes.InvalidArgument, "")
	}
}

func (ut *UserTorrentsImpl) RemoveTorrent(params *RemoveTorrentParams) (*prpc.RemoveTorrentRes, error) {
	return nil, nil
}
