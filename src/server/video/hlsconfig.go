package video

// ffmpeg -hwaccels shows
// ffmpeg -filters
// ffmpeg -h encoder=hevc_nvenc
// ffprobe -v 0 -of compact=p=0:nk=1 -show_entries stream=time_base -select_streams v:0 INPUT

var GlobalAudioParams = []KV{
	{"-c", "aac"},
	{"-b", "128k"},
	{"-ac", "2"},
}

var CudaGlobalDecode = []KV{
	{"-hwaccel", "cuda"},
	{"-hwaccel_output_format", "cuda"},
}

var CudaH264GlobalVideoParams = []KV{
	{"-c", "h264_nvenc"},
	{"-bf", "0"}, // disable b frame, for correcting pts and dts
	{"-fps_mode", "vfr"},
	{"-preset", "p7"},
	{"-profile", "high"},
	{"-cq", "26"},
	{"-level", "auto"},
}

var CudaSplitEncoderParams = []EncoderParams{
	{
		W:       3840,
		H:       2160,
		Filters: "scale_cuda=%d:%d:format=yuv420p",
		VCodecParams: []KV{
			{"-cq", "19"},
			{"-maxrate", "60000k"},
			{"-bufsize", "70000k"},
		},
		ACodecParams: []KV{
			{"-b", "400k"},
		},
	},
	{
		W:       1920,
		H:       1080,
		Filters: "scale_cuda=%d:%d:format=yuv420p",
		VCodecParams: []KV{
			{"-cq", "19"},
			{"-maxrate", "30000k"},
			{"-bufsize", "40000k"},
		},
		ACodecParams: []KV{
			{"-b", "300k"},
		},
	},
	{
		W:       1280,
		H:       720,
		Filters: "scale_cuda=%d:%d:format=yuv420p",
		VCodecParams: []KV{
			{"-maxrate", "5000k"},
			{"-bufsize", "6000k"},
		},
		ACodecParams: []KV{
			{"-b", "256k"},
		},
	},
}

var CudaGlobalDecode2 = []KV{
	{"-hwaccel", "cuda"},
}

var CudaEncoderParams2 = []EncoderParams{
	{
		W:       3840,
		H:       2160,
		Filters: "scale=%d:%d,format=yuv420p",
		VCodecParams: []KV{
			{"-cq", "19"},
			{"-maxrate", "60000k"},
			{"-bufsize", "70000k"},
		},
		ACodecParams: []KV{
			{"-b", "400k"},
		},
	},
	{
		W:       1920,
		H:       1080,
		Filters: "scale=%d:%d,format=yuv420p",
		VCodecParams: []KV{
			{"-cq", "19"},
			{"-maxrate", "30000k"},
			{"-bufsize", "40000k"},
		},
		ACodecParams: []KV{
			{"-b", "300k"},
		},
	},
	{
		W:       1280,
		H:       720,
		Filters: "scale=%d:%d,format=yuv420p",
		VCodecParams: []KV{
			{"-maxrate", "5000k"},
			{"-bufsize", "6000k"},
		},
		ACodecParams: []KV{
			{"-b", "256k"},
		},
	},
}

var QsvGlobalDecode = []KV{
	{"-init_hw_device", "qsv"},
	{"-hwaccel", "qsv"},
	{"-hwaccel_output_format", "qsv"},
}

var QsvH264GlobalVideoParams = []KV{
	{"-c", "h264_qsv"},
	{"-preset", "veryslow"},
	{"-profile", "high"},
	{"-global_quality", "25"},
}

var QsvSplitEncoderParams = []EncoderParams{
	{
		W:       3840,
		H:       2160,
		Filters: "scale_qsv=w=%d:h=%d:format=nv12",
		VCodecParams: []KV{
			{"-global_quality", "10"},
		},
		ACodecParams: []KV{
			{"-b", "400k"},
		},
	},
	{
		W:       1920,
		H:       1080,
		Filters: "scale_qsv=w=%d:h=%d:format=nv12",
		VCodecParams: []KV{
			{"-global_quality", "10"},
		},
		ACodecParams: []KV{
			{"-b", "300k"},
		},
	},
	{
		W:            1280,
		H:            720,
		Filters:      "scale_qsv=w=%d:h=%d:format=nv12",
		VCodecParams: []KV{},
		ACodecParams: []KV{
			{"-b", "256k"},
		},
	},
}

var SoGlobalDecode = []KV{
	//{"-threads", "3"},
}

var SoH264GlobalVideoParams = []KV{
	{"-c", "libx264"},
	{"-preset", "slower"},
	{"-crf", "26"},
	{"-threads", "6"},
	{"-pix_fmt", "yuv420p"},
}

var SoSplitEncoderParams = []EncoderParams{
	{
		W:       3840,
		H:       2160,
		Filters: "scale=%d:%d,format=yuv420p",
		VCodecParams: []KV{
			{"-crf", "19"},
			{"-maxrate", "60000k"},
			{"-bufsize", "70000k"},
		},
		ACodecParams: []KV{
			{"-b", "400k"},
		},
	},
	{
		W:       1920,
		H:       1080,
		Filters: "scale=%d:%d,format=yuv420p",
		VCodecParams: []KV{
			{"-crf", "19"},
			{"-maxrate", "30000k"},
			{"-bufsize", "40000k"},
		},
		ACodecParams: []KV{
			{"-b", "300k"},
		},
	},
	{
		W:       1280,
		H:       720,
		Filters: "scale=%d:%d,format=yuv420p",
		VCodecParams: []KV{
			{"-maxrate", "5000k"},
			{"-bufsize", "6000k"},
		},
		ACodecParams: []KV{
			{"-b", "256k"},
		},
	},
}
