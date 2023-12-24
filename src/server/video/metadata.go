package video

import "time"

const (
	VideoType = "video"
	AudioType = "audio"
)

type Metadata struct {
	Streams []*Stream `json:"streams"`
	Format  Format    `json:"format"`
}
type Disposition struct {
	Default         int `json:"default"`
	Dub             int `json:"dub"`
	Original        int `json:"original"`
	Comment         int `json:"comment"`
	Lyrics          int `json:"lyrics"`
	Karaoke         int `json:"karaoke"`
	Forced          int `json:"forced"`
	HearingImpaired int `json:"hearing_impaired"`
	VisualImpaired  int `json:"visual_impaired"`
	CleanEffects    int `json:"clean_effects"`
	AttachedPic     int `json:"attached_pic"`
	TimedThumbnails int `json:"timed_thumbnails"`
}

type Tags struct {
	BPSEng                      string    `json:"BPS-eng"`
	DURATIONEng                 string    `json:"DURATION-eng"`
	NUMBEROFFRAMESEng           string    `json:"NUMBER_OF_FRAMES-eng"`
	NUMBEROFBYTESEng            string    `json:"NUMBER_OF_BYTES-eng"`
	STATISTICSWRITINGAPPEng     string    `json:"_STATISTICS_WRITING_APP-eng"`
	STATISTICSWRITINGDATEUTCEng string    `json:"_STATISTICS_WRITING_DATE_UTC-eng"`
	STATISTICSTAGSEng           string    `json:"_STATISTICS_TAGS-eng"`
	Encoder                     string    `json:"encoder"`
	CreationTime                time.Time `json:"creation_time"`
	Language                    string    `json:"language"`
	HandlerName                 string    `json:"handler_name"`
}

type Stream struct {
	Index              int         `json:"index"`
	CodecName          string      `json:"codec_name"`
	CodecLongName      string      `json:"codec_long_name"`
	Profile            string      `json:"profile,omitempty"`
	CodecType          string      `json:"codec_type"`
	CodecTagString     string      `json:"codec_tag_string"`
	CodecTag           string      `json:"codec_tag"`
	Width              int         `json:"width,omitempty"`
	Height             int         `json:"height,omitempty"`
	CodedWidth         int         `json:"coded_width,omitempty"`
	CodedHeight        int         `json:"coded_height,omitempty"`
	ClosedCaptions     int         `json:"closed_captions,omitempty"`
	HasBFrames         int         `json:"has_b_frames,omitempty"`
	SampleAspectRatio  string      `json:"sample_aspect_ratio,omitempty"`
	DisplayAspectRatio string      `json:"display_aspect_ratio,omitempty"`
	PixFmt             string      `json:"pix_fmt,omitempty"`
	Level              int         `json:"level,omitempty"`
	ColorRange         string      `json:"color_range,omitempty"`
	ColorSpace         string      `json:"color_space,omitempty"`
	ColorTransfer      string      `json:"color_transfer,omitempty"`
	ColorPrimaries     string      `json:"color_primaries,omitempty"`
	ChromaLocation     string      `json:"chroma_location,omitempty"`
	Refs               int         `json:"refs,omitempty"`
	RFrameRate         string      `json:"r_frame_rate"`
	AvgFrameRate       string      `json:"avg_frame_rate"`
	TimeBase           string      `json:"time_base"`
	StartPts           int         `json:"start_pts"`
	StartTime          string      `json:"start_time"`
	Disposition        Disposition `json:"disposition"`
	SampleFmt          string      `json:"sample_fmt,omitempty"`
	SampleRate         string      `json:"sample_rate,omitempty"`
	Channels           int         `json:"channels,omitempty"`
	ChannelLayout      string      `json:"channel_layout,omitempty"`
	BitsPerSample      int         `json:"bits_per_sample,omitempty"`
	BitsPerRawSample   string      `json:"bits_per_raw_sample,omitempty"`
	DurationTs         int         `json:"duration_ts,omitempty"`
	Duration           string      `json:"duration,omitempty"`
	Tags               Tags        `json:"tags,omitempty"`
}

type Format struct {
	Filename       string `json:"filename"`
	NbStreams      int    `json:"nb_streams"`
	NbPrograms     int    `json:"nb_programs"`
	FormatName     string `json:"format_name"`
	FormatLongName string `json:"format_long_name"`
	StartTime      string `json:"start_time"`
	Duration       string `json:"duration"`
	Size           string `json:"size"`
	BitRate        string `json:"bit_rate"`
	ProbeScore     int    `json:"probe_score"`
	Tags           Tags   `json:"tags"`
}
