package bt

import (
	"errors"
	"pnas/log"
	"pnas/prpc"
	"pnas/video"
	"sync"
)

type InfoHash struct {
	Version int32
	Hash    string
}

type TorrentBase struct {
	InfoHash    InfoHash
	Name        string
	SavePath    string
	TotalSize   int64
	PieceLength int32
	NumPieces   int32
}

type Torrent struct {
	base TorrentBase

	mtx        sync.Mutex
	files      []File
	state      prpc.BtStateEnum
	resumeData []byte
}

func NewTorrent(base *TorrentBase) *Torrent {
	return &Torrent{
		base: *base,
	}
}

func (t *Torrent) UpdateBaseInfo(base *TorrentBase) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.base = *base
}

func (t *Torrent) GetBaseInfo() TorrentBase {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.base
}

func (t *Torrent) UpdateState(newState prpc.BtStateEnum) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.state = newState
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

func (t *Torrent) UpdateFilesInfo(fs []File) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if len(t.files) != len(fs) {
		if len(t.files) > 0 {
			log.Warnf("torrent:%s change file:%v, new name:%v", t.base.Name, t.files, fs)
		}
		t.files = make([]File, len(fs))
		copy(t.files, fs)
	}
}

func (t *Torrent) GetFiles() []File {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	ret := make([]File, len(t.files))
	copy(ret, t.files)
	return ret
}

func (t *Torrent) UpdateFileState(index int, st prpc.BtFile_State) error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	if int(index) >= len(t.files) {
		return errors.New("file index out of range")
	}
	t.files[index].St = st
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

func (t *Torrent) UpdateResumeData(r []byte) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.resumeData = r
}

func (t *Torrent) GetResumeData() []byte {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.resumeData
}
