package user

import (
	"encoding/hex"
	"fmt"
	"pnas/bt"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/setting"
	"pnas/video"
	"strings"
	"sync"
	"time"

	"github.com/pkg/errors"
)

type UpdateTorrentParams struct {
	Base       *bt.TorrentBase
	State      prpc.BtStateEnum
	FileNames  []bt.File
	ResumeData []byte
}

type FileCompleted struct {
	InfoHash  bt.InfoHash
	FileIndex int32
}

type UserTorrents interface {
	UpdateTorrent(params *UpdateTorrentParams)
	BtFileStateComplete(fs *FileCompleted)
	AddTorrent(userId ID, t *bt.TorrentBase) error
	AddUserTorrent(userId ID, infoHash bt.InfoHash)
	HasTorrent(userId ID, infoHash bt.InfoHash) bool
	GetTorrent(infoHash bt.InfoHash) (*bt.Torrent, error)
	RealHasTorrent(userId ID, infoHash bt.InfoHash) bool
	RemoveUserTorrent(userId ID, infoHash bt.InfoHash) error
	RemoveTorrent(userId ID, infoHash bt.InfoHash) error
}

func LoadTorrent(infoHash bt.InfoHash) ([]byte, error) {
	sql := `select resume_data from torrent where info_hash=? and version=?`
	var resumeData []byte
	err := db.QueryRow(sql, infoHash.Hash, infoHash.Version).Scan(&resumeData)
	return resumeData, err
}

func LoadDownloadingTorrent() [][]byte {
	sql := `select resume_data from user_torrent u 
					left join torrent t on u.torrent_id = t.id`

	rows, err := db.Query(sql)
	if err != nil {
		log.Warnf("[bt] failed to load downloading torrent err: %v", err)
		return [][]byte{}
	}
	defer rows.Close()

	var resumData [][]byte
	for rows.Next() {
		var resume []byte
		err = rows.Scan(&resume)
		if err != nil {
			log.Warnf("[bt] failed to load downloading torrent err: %v", err)
			continue
		}
		resumData = append(resumData, resume)
	}
	return resumData
}

type UserTorrentsImpl struct {
	UserTorrents
	mtx          sync.Mutex
	torrents     map[bt.InfoHash]*bt.Torrent
	userTorrents map[bt.InfoHash]map[ID]bool
}

func (ut *UserTorrentsImpl) Init() {
	ut.torrents = make(map[bt.InfoHash]*bt.Torrent)
	ut.userTorrents = make(map[bt.InfoHash]map[ID]bool)
}

func updateBtFileType(st *bt.Torrent, fileIndex int, absFileName string) {
	if video.IsVideo(absFileName) {
		st.UpdateFileType(fileIndex, bt.FileVideoType)
		meta, _ := video.GetMetadata(absFileName)
		st.UpdateVideoFileMeta(fileIndex, meta)
	}
	if video.IsSubTitle(absFileName) {
		st.UpdateFileType(fileIndex, bt.FileSubtitleType)
	}
	if video.IsAudio(absFileName) {
		st.UpdateFileType(fileIndex, bt.FileAudioType)
	}
}

func (ut *UserTorrentsImpl) UpdateTorrent(params *UpdateTorrentParams) {
	srcState := prpc.BtStateEnum_unknown
	ut.mtx.Lock()
	st, ok := ut.torrents[params.Base.InfoHash]
	if !ok {
		st = bt.NewTorrent(params.Base)
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

		if files[i].FileType == bt.FileUnknownType {
			absFileName := setting.GS().Bt.SavePath + "/" + fileName
			updateBtFileType(st, i, absFileName)
			ft, _ := st.GetFileType(i)
			log.Debugf("[bt] torrent:%s file: %s type: %d", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, ft)
			if log.EnabledDebug() {
				meta, _ := video.GetMetadata(absFileName)
				log.Debugf("[bt] torrent:%s file: %s format: %s", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, meta.Format.FormatName)
				if meta != nil {
					for _, s := range meta.Streams {
						log.Debugf("[bt] torrent:%s file: %s s: %d codeType: %s", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, s.Index, s.CodecType)
					}
				}
			}
		}

		v, err := video.GetVideoByFileName(fileName)
		if err == nil && v.HlsCreated {
			continue
		}

		if ft, _ := st.GetFileType(i); (ft | bt.FileVideoType) != 0 {
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
		meta, _ := video.GetMetadata(absFileName)
		if meta != nil {
			log.Debugf("[bt] torrent:%s file: %s format: %s", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, meta.Format.FormatName)
			for _, s := range meta.Streams {
				log.Debugf("[bt] torrent:%s file: %s s: %d codeType: %s", hex.EncodeToString([]byte(baseInfo.InfoHash.Hash)), absFileName, s.Index, s.CodecType)
			}
		}
	}
}

func (ut *UserTorrentsImpl) AddTorrent(userId ID, t *bt.TorrentBase) error {
	ut.mtx.Lock()
	_, ok := ut.torrents[t.InfoHash]
	if !ok {
		ut.torrents[t.InfoHash] = bt.NewTorrent(t)
	}
	_, ok = ut.userTorrents[t.InfoHash]
	if !ok {
		ut.userTorrents[t.InfoHash] = make(map[ID]bool)
	}
	_, ok = ut.userTorrents[t.InfoHash][userId]
	if ok {
		ut.mtx.Unlock()
		return errors.New("repeat download")
	}
	ut.mtx.Unlock()

	_, err := db.Query("call new_torrent(?, ?, ?, @torrent_id);", t.InfoHash.Version, t.InfoHash.Hash, userId)
	if err != nil {
		log.Warnf("[user] %d failed to add torrent: %s err: %v", userId, hex.EncodeToString([]byte(t.InfoHash.Hash)), err)
		return err
	}

	ut.mtx.Lock()
	ut.userTorrents[t.InfoHash][userId] = true
	ut.mtx.Unlock()
	return nil
}

func (ut *UserTorrentsImpl) AddUserTorrent(userId ID, infoHash bt.InfoHash) {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	uts, ok := ut.userTorrents[infoHash]
	if !ok {
		uts = make(map[ID]bool)
		ut.userTorrents[infoHash] = uts
	}
	uts[userId] = true
}

func (ut *UserTorrentsImpl) RealHasTorrent(userId ID, infoHash bt.InfoHash) bool {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	i, ok := ut.userTorrents[infoHash]
	if !ok {
		return false
	}
	_, ok = i[userId]
	return ok
}

func (ut *UserTorrentsImpl) HasTorrent(userId ID, infoHash bt.InfoHash) bool {
	if userId == AdminId {
		return true
	}
	return ut.RealHasTorrent(userId, infoHash)
}

func (ut *UserTorrentsImpl) GetTorrent(infoHash bt.InfoHash) (*bt.Torrent, error) {
	ut.mtx.Lock()
	defer ut.mtx.Unlock()
	t, ok := ut.torrents[infoHash]
	if !ok {
		return nil, errors.New("not found bt")
	}
	return t, nil
}

func (ut *UserTorrentsImpl) RemoveUserTorrent(userId ID, infoHash bt.InfoHash) error {
	if !ut.RealHasTorrent(userId, infoHash) {
		return errors.New("not found")
	}
	sql := "select id from torrent where info_hash=? and version=?"
	rows, err := db.Query(sql, infoHash.Hash, infoHash.Version)
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var tid int64
		err := rows.Scan(&tid)
		if err != nil {
			return err
		}
		sql := "delete from user_torrent where user_id=? and torrent_id=?"
		_, err = db.Exec(sql, userId, tid)
		if err != nil {
			return err
		}
	}
	ut.mtx.Lock()
	delete(ut.userTorrents[infoHash], userId)
	ut.mtx.Unlock()
	return nil
}

func (ut *UserTorrentsImpl) RemoveTorrent(userId ID, infoHash bt.InfoHash) error {
	if !ut.HasTorrent(userId, infoHash) {
		return errors.New("not found")
	}
	var uids []string
	ut.mtx.Lock()
	us, ok := ut.userTorrents[infoHash]
	if ok {
		for u, k := range us {
			if k {
				uids = append(uids, fmt.Sprintf("user_id=%d", u))
			}
		}
	}
	ut.mtx.Unlock()
	if len(uids) == 0 {
		return nil
	}

	sql := "select id from torrent where info_hash=? and version=?"
	rows, err := db.Query(sql, infoHash.Hash, infoHash.Version)
	if err != nil {
		return err
	}
	defer rows.Close()
	userIdsCond := strings.Join(uids, " or ")
	sql = fmt.Sprintf("delete from user_torrent where %s and torrent_id=?", userIdsCond)
	for rows.Next() {
		var tid int64
		err := rows.Scan(&tid)
		if err != nil {
			return err
		}
		_, err = db.Exec(sql, tid)
		if err != nil {
			return err
		}
	}
	ut.mtx.Lock()
	delete(ut.userTorrents, infoHash)
	ut.mtx.Unlock()
	return nil
}
