package video

import "strings"

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

func IsSubTitle(meta *Metadata) bool {
	f := strings.Contains(meta.Format.FormatName, "srt") ||
		strings.Contains(meta.Format.FormatName, "vtt") ||
		strings.Contains(meta.Format.FormatName, "matroska")
	return f && HasSubtitleStream(meta)
}

func IsSubTitle2(absFileName string) bool {
	meta, err := GetMetadata(absFileName)
	if err != nil {
		return false
	}
	return IsSubTitle(meta)
}

func IsVideo(meta *Metadata) bool {
	f := strings.Contains(meta.Format.FormatName, "mp4") ||
		strings.Contains(meta.Format.FormatName, "matroska") ||
		strings.Contains(meta.Format.FormatName, "mpegts")
	return f && HasVideoStream(meta)
}

func IsVideo2(absFileName string) bool {
	meta, err := GetMetadata(absFileName)
	if err != nil {
		return false
	}
	return IsVideo(meta)
}
func IsAudio(meta *Metadata) bool {
	f := strings.Contains(meta.Format.FormatName, "mp3") ||
		strings.Contains(meta.Format.FormatName, "matroska")
	return f && HasAudioStream(meta)
}
func IsAudio2(absFileName string) bool {
	meta, err := GetMetadata(absFileName)
	if err != nil {
		return false
	}
	return IsAudio(meta)
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
