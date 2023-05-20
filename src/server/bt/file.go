package bt

import (
	"pnas/prpc"
	"pnas/video"
)

const (
	FileUnknownType = 0

	FileVideoType = 1 << (iota - 1)
	FileSubtitleType
	FileAudioType
)

type FileType int

type File struct {
	Name       string
	Index      int32
	St         prpc.BtFile_State
	TotalSize  int64
	Downloaded int64
	FileType   FileType
	Meta       *video.Metadata
}
