package bt

import (
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
}

type Torrent struct {
	base TorrentBase

	mtx        sync.Mutex
	hasBase    bool
	files      []File
	state      prpc.BtStateEnum
	updateTime time.Time
	whoHas     map[ptype.UserID]bool

	btClient *BtClient
	lastSave time.Time
}

func (t *Torrent) init() {
	t.whoHas = make(map[ptype.UserID]bool)
	t.whoHas[ptype.AdminId] = true
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
	ret := make([]File, len(t.files))
	copy(ret, t.files)
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
	t.files = make([]File, len(ti.Files))
	for i, f := range ti.Files {
		t.files[i].Name = f.Name
		t.files[i].Index = f.Index
		t.files[i].St = f.St
		t.files[i].TotalSize = f.TotalSize
	}
	t.hasBase = true
}

func (t *Torrent) updateStatus(s *prpc.TorrentStatus) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	old := t.state
	t.state = s.State

	if t.state != s.State && s.State == prpc.BtStateEnum_finished {
		for i := range t.files {
			t.updateFileTypeLocked(i)
		}
	}

	if t.state == prpc.BtStateEnum_downloading ||
		old != prpc.BtStateEnum_finished && t.state == prpc.BtStateEnum_finished {
		now := time.Now()
		if now.Sub(t.lastSave) > time.Second*10 {
			// req := &prpc.GetResumeDataReq{
			// 	InfoHash: GetInfoHash(&t.base.InfoHash),
			// }
			// rd, err := t.btClient.GetResumeData(context.Background(), req)
			// if err == nil {
			// 	saveResumeData(&t.base.InfoHash, rd.ResumeData)
			// }
			t.lastSave = now
		}
	}
}

func (t *Torrent) updateFileTypeLocked(index int) {
	log.Infof("[bt] torrent: %s file: %s completed", hex.EncodeToString([]byte(t.base.InfoHash.Hash)), t.files[index].Name)
	fileName := t.files[index].Name
	absFileName := t.base.SavePath + "/" + fileName
	meta, err := video.GetMetadata(absFileName)
	if err == nil {
		if video.IsVideo(meta) {
			t.files[index].FileType |= FileVideoType
			t.files[index].Meta = meta
		}
		if video.IsSubTitle(meta) {
			t.files[index].FileType |= FileSubtitleType
		}
		if video.IsAudio(meta) {
			t.files[index].FileType |= FileAudioType
		}
		log.Debugf("[bt] torrent:%s file: %s type: %d", hex.EncodeToString([]byte(t.base.InfoHash.Hash)),
			absFileName, t.files[index].FileType)
	}
}

func (t *Torrent) UpdateFileState(index int, st prpc.BtFile_State) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if int(index) >= len(t.files) {
		return errors.New("file index out of range")
	}
	t.updateFileTypeLocked(index)
	return nil
}

func (t *Torrent) GetFileType(fileIndex int) (FileType, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if fileIndex >= len(t.files) {
		return 0, errors.New("file index out of range")
	}
	return t.files[fileIndex].FileType, nil
}

func (t *Torrent) UpdateFileType(fileIndex int, ft FileType) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if fileIndex >= len(t.files) {
		return errors.New("file index out of range")
	}
	t.files[fileIndex].FileType |= ft
	return nil
}

func (t *Torrent) UpdateVideoFileMeta(fileIndex int, meta *video.Metadata) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if fileIndex >= len(t.files) {
		return errors.New("file index out of range")
	}
	t.files[fileIndex].Meta = meta
	return nil
}

func (t *Torrent) GetVideoFileMeta(fileIndex int) (*video.Metadata, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if fileIndex >= len(t.files) {
		return nil, errors.New("file index out of range")
	}
	return t.files[fileIndex].Meta, nil
}

func (t *Torrent) GetFileState(fileIndex int) (prpc.BtFile_State, error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if int(fileIndex) >= len(t.files) {
		return 0, errors.New("file index out of range")
	}
	return t.files[fileIndex].St, nil
}

func (t *Torrent) GetUpdateTime() time.Time {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.updateTime
}
