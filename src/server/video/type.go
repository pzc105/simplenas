package video

import "strings"

type ID int64

const (
	VideoStreamType    = "video"
	AudioStreamType    = "audio"
	SubtitleStreamType = "subtitle"
)

func IsVideoStream(stream *Stream) bool {
	return stream.CodecType == "video"
}

func IsAudioStream(stream *Stream) bool {
	return stream.CodecType == "audio"
}

func IsSubTitleStream(stream *Stream) bool {
	return stream.CodecType == "subtitle"
}

func IsSubTitle(absFileName string) bool {
	meta, err := GetMetadata(absFileName)
	if err != nil {
		return false
	}
	f := strings.Contains(meta.Format.FormatName, "srt") ||
		strings.Contains(meta.Format.FormatName, "vtt") ||
		strings.Contains(meta.Format.FormatName, "matroska")
	return f && HasSubtitleStream(meta)
}

func IsVideo(absFileName string) bool {
	meta, err := GetMetadata(absFileName)
	if err != nil {
		return false
	}
	f := strings.Contains(meta.Format.FormatName, "mp4") ||
		strings.Contains(meta.Format.FormatName, "matroska") ||
		strings.Contains(meta.Format.FormatName, "mpegts")
	return f && HasVideoStream(meta)
}

func IsAudio(absFileName string) bool {
	meta, err := GetMetadata(absFileName)
	if err != nil {
		return false
	}
	f := strings.Contains(meta.Format.FormatName, "mp3") ||
		strings.Contains(meta.Format.FormatName, "matroska")
	return f && HasAudioStream(meta)
}

func HasVideoStream(meta *Metadata) bool {
	for _, s := range meta.Streams {
		if IsVideoStream(s) {
			return true
		}
	}
	return false
}

func HasAudioStream(meta *Metadata) bool {
	for _, s := range meta.Streams {
		if IsAudioStream(s) {
			return true
		}
	}
	return false
}

func HasSubtitleStream(meta *Metadata) bool {
	for _, s := range meta.Streams {
		if IsSubTitleStream(s) {
			return true
		}
	}
	return false
}
