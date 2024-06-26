package bt

import (
	"context"
	"errors"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/video"
	"sync"
	"time"
)

type InfoHash struct {
	Version int32
	Hash    string
}

type TorrentBase struct {
	Id          ptype.TorrentID
	InfoHash    InfoHash
	Name        string
	SavePath    string
	TotalSize   int64
	PieceLength int32
	NumPieces   int32
	Introduce   string
	Files       []File
	MagnetUri   string
}

type Torrent struct {
	mtx             sync.Mutex
	base            TorrentBase
	hasBase         bool
	updatedFileType bool
	state           prpc.BtStateEnum
	updateTime      time.Time
	whoHas          map[ptype.UserID]*userData
	removed         bool
	lastSt          *prpc.TorrentStatus

	ut *UserTorrentsImpl

	btClient *BtClient
	lastSave time.Time
}

func (t *Torrent) init(ut *UserTorrentsImpl) {
	t.whoHas = make(map[ptype.UserID]*userData)
	t.removed = false
	t.lastSt = &prpc.TorrentStatus{
		InfoHash: GetInfoHash(&t.base.InfoHash),
		Name:     t.base.Name,
		Total:    t.base.TotalSize,
		State:    t.state,
	}
	if t.state == prpc.BtStateEnum_seeding {
		t.lastSt.TotalDone = t.lastSt.Total
	}
	t.ut = ut
}

// must with UserTorrentsImpl.mtx
func (t *Torrent) addTask(ud *userData, params *DownloadTaskParams) error {
	ud.addTorrent(t)
	ud.setTaskCallback(&SetTaskCallbackParams{
		TaskId:    params.TaskId,
		TorrentId: t.base.Id,
		Callback:  params.Callback,
	})
	t.mtx.Lock()
	if t.removed {
		ud.removeTorrent(t.base.Id, true)
		t.mtx.Unlock()
		return errors.New("removed torrent")
	}
	t.whoHas[ud.userId] = ud
	lastSt := t.lastSt
	t.mtx.Unlock()
	log.Debugf("[bt] uid:%d add torrent:%d", ud.userId, t.base.Id)

	if params.Callback != nil {
		ud.callbackTaskQueue.Put(func() {
			params.Callback(nil, lastSt)
		})
	}
	return nil
}

// must with UserTorrentsImpl.mtx
func (t *Torrent) addUser(ud *userData) error {
	t.mtx.Lock()
	if t.removed {
		t.mtx.Unlock()
		return errors.New("removed torrent")
	}
	t.whoHas[ud.userId] = ud
	t.mtx.Unlock()
	log.Debugf("[bt] uid:%d add torrent:%d", ud.userId, t.base.Id)
	ud.addTorrent(t)
	return nil
}

// must with UserTorrentsImpl.mtx
func (t *Torrent) removeUser(uid ptype.UserID) error {
	t.mtx.Lock()
	ud, ok := t.whoHas[uid]
	delete(t.whoHas, uid)
	t.mtx.Unlock()
	if !ok {
		return errors.New("not found user")
	}
	log.Debugf("[bt] torrent:%d remove uid:%v", t.base.Id, uid)
	return ud.removeTorrent(t.base.Id, true)
}

func (t *Torrent) getUserCount() int {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return len(t.whoHas)
}

// must with UserTorrentsImpl.mtx
func (t *Torrent) remove() error {
	t.mtx.Lock()
	delResumeData(&t.base.InfoHash)
	t.removed = true
	var uids []ptype.UserID
	uds := []*userData{}
	for uid, ud := range t.whoHas {
		uids = append(uids, uid)
		uds = append(uds, ud)
	}
	t.whoHas = make(map[ptype.UserID]*userData)
	t.mtx.Unlock()
	for _, ud := range uds {
		ud.removeTorrent(t.base.Id, false)
	}
	log.Debugf("[bt] torrent:%d remove uids:%v", t.base.Id, uids)
	return deleteUserTorrentids(t.base.Id, uids...)
}

func (t *Torrent) hasBaseInfo() bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.hasBase
}

func (t *Torrent) GetBaseInfo() TorrentBase {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.base
}

func (t *Torrent) GetId() ptype.TorrentID {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.base.Id
}

func (t *Torrent) GetState() prpc.BtStateEnum {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.state
}

func (t *Torrent) GetInfoHash() InfoHash {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.base.InfoHash
}

func (t *Torrent) GetFiles() []File {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	ret := make([]File, len(t.base.Files))
	copy(ret, t.base.Files)
	return ret
}

func (t *Torrent) updateTorrentInfo(ti *prpc.TorrentInfo) {
	if len(ti.Name) == 0 && ti.NumPieces <= 0 {
		return
	}
	t.mtx.Lock()
	defer t.mtx.Unlock()

	sql := `update torrent set name=?, total_size=?, piece_length=?, num_pieces=? where version=? and info_hash=?`
	_, err := db.Exec(sql, ti.Name, ti.TotalSize, ti.PieceLength, ti.NumPieces, t.base.InfoHash.Version, t.base.InfoHash.Hash)
	if err != nil {
		log.Warnf("failed to update torrent err: %v", err)
		return
	}

	t.base.Name = ti.Name
	t.base.NumPieces = ti.NumPieces
	t.base.PieceLength = ti.PieceLength
	t.base.SavePath = ti.SavePath
	t.base.TotalSize = ti.TotalSize
	t.base.Files = make([]File, len(ti.Files))
	for i, f := range ti.Files {
		t.base.Files[i].Name = f.Name
		t.base.Files[i].Index = f.Index
		t.base.Files[i].St = f.St
		t.base.Files[i].TotalSize = f.TotalSize
	}
	t.hasBase = true

	if IsDownloadAll(t.state) && !t.updatedFileType {
		log.Infof("[bt] torrent: %d %s completed", t.base.Id, t.base.Name)
		for i := range t.base.Files {
			t.updateFileTypeLocked(i)
		}
		t.updatedFileType = true
	}
}

func (t *Torrent) updateStatus(s *prpc.TorrentStatus) {
	t.mtx.Lock()
	old := t.state
	t.state = s.State
	t.lastSt = s

	if old != s.State {
		sql := `update torrent set state=? where version=? and info_hash=?`
		_, err := db.Exec(sql, s.State, t.base.InfoHash.Version, t.base.InfoHash.Hash)
		if err != nil {
			log.Warnf("failed to update torrent:%d err:%v", t.base.Id, err)
		}
	}

	if IsDownloadAll(s.State) && !t.updatedFileType {
		log.Infof("[bt] torrent: %d %s completed", t.base.Id, t.base.Name)
		for i := range t.base.Files {
			t.updateFileTypeLocked(i)
		}
		t.updatedFileType = true
	}

	now := time.Now()
	needSaveResuem := func() bool {
		if t.removed {
			return false
		}
		if old != prpc.BtStateEnum_seeding {
			return true
		}
		if t.state == prpc.BtStateEnum_downloading && now.Sub(t.lastSave) > time.Second*10 {
			return true
		}
		return false
	}
	if needSaveResuem() {
		req := &prpc.GetResumeDataReq{
			InfoHash: GetInfoHash(&t.base.InfoHash),
		}
		rd, err := t.btClient.GetResumeData(context.Background(), req)
		if err == nil {
			saveResumeData(&t.base.InfoHash, rd.ResumeData)
		}
		t.lastSave = now
	}
	uds := []*userData{}
	hasAdmin := false
	for _, ud := range t.whoHas {
		uds = append(uds, ud)
		if ud.userId == ptype.AdminId {
			hasAdmin = true
		}
	}
	t.mtx.Unlock()

	if !hasAdmin {
		ud := t.ut.getUserData(ptype.AdminId)
		if ud != nil {
			ud.onBtStatus(t.base.Id, s)
		}
	}

	for _, ud := range uds {
		ud.onBtStatus(t.base.Id, s)
	}
}

func (t *Torrent) updateFileTypeLocked(index int) {
	fileName := t.base.Files[index].Name
	absFileName := t.base.SavePath + "/" + fileName
	meta, err := video.GetMetadata(absFileName)
	if err == nil {
		if video.IsVideo(meta) {
			t.base.Files[index].FileType |= FileVideoType
			t.base.Files[index].Meta = meta
		}
		if video.IsSubTitle(meta) {
			t.base.Files[index].FileType |= FileSubtitleType
		}
		if video.IsAudio(meta) {
			t.base.Files[index].FileType |= FileAudioType
		}
		log.Debugf("[bt] torrent:%d file:%s type:%d", t.base.Id, absFileName, t.base.Files[index].FileType)
	} else {
		log.Debugf("[bt] torrent:%d file:%s unkonwn type", t.base.Id, absFileName)
	}
}

func (t *Torrent) UpdateFileState(index int, st prpc.BtFile_State) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if int(index) >= len(t.base.Files) {
		return errors.New("file index out of range")
	}
	t.updateFileTypeLocked(index)
	return nil
}

func (t *Torrent) GetFileType(fileIndex int) (FileType, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if fileIndex >= len(t.base.Files) {
		return 0, errors.New("file index out of range")
	}
	return t.base.Files[fileIndex].FileType, nil
}

func (t *Torrent) UpdateFileType(fileIndex int, ft FileType) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if fileIndex >= len(t.base.Files) {
		return errors.New("file index out of range")
	}
	t.base.Files[fileIndex].FileType |= ft
	return nil
}

func (t *Torrent) UpdateVideoFileMeta(fileIndex int, meta *video.Metadata) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if fileIndex >= len(t.base.Files) {
		return errors.New("file index out of range")
	}
	t.base.Files[fileIndex].Meta = meta
	return nil
}

func (t *Torrent) GetVideoFileMeta(fileIndex int) (*video.Metadata, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if fileIndex >= len(t.base.Files) {
		return nil, errors.New("file index out of range")
	}
	return t.base.Files[fileIndex].Meta, nil
}

func (t *Torrent) GetFileState(fileIndex int) (prpc.BtFile_State, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if int(fileIndex) >= len(t.base.Files) {
		return 0, errors.New("file index out of range")
	}
	return t.base.Files[fileIndex].St, nil
}

func (t *Torrent) GetUpdateTime() time.Time {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.updateTime
}

func (t *Torrent) UpdateMagnetUri(magnetUri string) bool {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	sql := `update torrent set magnet_uri=? where id=?`
	r, err := db.Exec(sql, magnetUri, t.base.Id)
	if err != nil {
		log.Warnf("[bt] failed to update magnet uri err: %v", err)
		return false
	}
	af, err := r.RowsAffected()
	if err != nil || af == 0 {
		return false
	}
	return true
}

func (t *Torrent) GetMagnetUri() string {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.base.MagnetUri
}
