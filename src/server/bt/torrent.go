package bt

import (
	"context"
	"encoding/hex"
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
	mtx        sync.Mutex
	base       TorrentBase
	hasBase    bool
	state      prpc.BtStateEnum
	updateTime time.Time
	whoHas     map[ptype.UserID]bool

	btClient *BtClient
	lastSave time.Time
}

func (t *Torrent) init() {
	t.whoHas = make(map[ptype.UserID]bool)
}

func (t *Torrent) addUser(uid ptype.UserID) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.whoHas[uid] = true
}

func (t *Torrent) removeUser(uid ptype.UserID) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	delete(t.whoHas, uid)
}

func (t *Torrent) getAllUser() []ptype.UserID {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	var ret []ptype.UserID
	for uid := range t.whoHas {
		ret = append(ret, uid)
	}
	return ret
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

	sql := `update torrent set name=?, total_size=?, piece_length=?, num_pieces=?`
	_, err := db.Exec(sql, ti.Name, ti.TotalSize, ti.PieceLength, ti.NumPieces)
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
}

func (t *Torrent) updateStatus(s *prpc.TorrentStatus) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	old := t.state
	t.state = s.State

	if !IsDownloadAll(old) && IsDownloadAll(s.State) {
		log.Infof("[bt] torrent: [%s] %s completed", hex.EncodeToString([]byte(t.base.InfoHash.Hash)), t.base.Name)
		for i := range t.base.Files {
			t.updateFileTypeLocked(i)
		}
	}

	if s.State == prpc.BtStateEnum_downloading ||
		old != prpc.BtStateEnum_seeding && s.State == prpc.BtStateEnum_seeding {
		now := time.Now()
		if now.Sub(t.lastSave) > time.Second*10 {
			req := &prpc.GetResumeDataReq{
				InfoHash: GetInfoHash(&t.base.InfoHash),
			}
			rd, err := t.btClient.GetResumeData(context.Background(), req)
			if err == nil {
				saveResumeData(&t.base.InfoHash, rd.ResumeData)
			}
			t.lastSave = now
		}
	}
}

func (t *Torrent) updateFileTypeLocked(index int) {
	log.Infof("[bt] torrent: [%s] file: %s completed", hex.EncodeToString([]byte(t.base.InfoHash.Hash)), t.base.Files[index].Name)
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
		log.Debugf("[bt] torrent:%s file: %s type: %d", hex.EncodeToString([]byte(t.base.InfoHash.Hash)),
			absFileName, t.base.Files[index].FileType)
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
	r, err := db.Exec(sql, t.base.Id, magnetUri)
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
