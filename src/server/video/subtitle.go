package video

import (
	"fmt"
	"os"
	"os/exec"
	"pnas/log"
)

func GetStreams(streams []*Stream, pred func(s *Stream) bool) []*Stream {
	var vStreams []*Stream
	for _, s := range streams {
		if pred(s) {
			vStreams = append(vStreams, s)
		}
	}
	return vStreams
}

type GenSubtitleOpts struct {
	InputFileName string
	OutDir        string
	SubtitleName  string
	Format        string
	Suffix        string
}

func GenSubtitle(params *GenSubtitleOpts) error {
	os.MkdirAll(params.OutDir, 0755)
	meta, err := GetMetadata(params.InputFileName)
	if err != nil {
		return err
	}
	streams := GetStreams(meta.Streams, IsSubTitleStream)
	if len(streams) == 0 {
		return nil
	}

	os.MkdirAll(params.OutDir, 0755)

	var cmdParams []string
	cmdParams = append(cmdParams, "-threads")
	cmdParams = append(cmdParams, "4")
	cmdParams = append(cmdParams, "-i")
	cmdParams = append(cmdParams, params.InputFileName)
	cmdParams = append(cmdParams, "-threads")
	cmdParams = append(cmdParams, "4")
	cmdParams = append(cmdParams, "-v")
	cmdParams = append(cmdParams, "error")
	nameStatus := make(map[string]int)
	for i, s := range streams {
		cmdParams = append(cmdParams, "-map")
		cmdParams = append(cmdParams, fmt.Sprintf("0:%d", s.Index))
		cmdParams = append(cmdParams, fmt.Sprintf("-c:s:%d", i))
		cmdParams = append(cmdParams, params.Format)

		lang := s.Tags.Language
		name := fmt.Sprintf("%s/%s", params.OutDir, params.SubtitleName)
		if len(lang) > 0 {
			name += "." + lang
		}
		c, ok := nameStatus[name]
		if !ok {
			nameStatus[name] = 0
			name = fmt.Sprintf("%s.%s", name, params.Suffix)
		} else {
			nameStatus[name] = c + 1
			name = fmt.Sprintf("%s%d.%s", name, c+1, params.Suffix)
		}
		cmdParams = append(cmdParams, name)
	}
	cmdParams = append(cmdParams, "-y")

	cmd := exec.Command("ffmpeg", cmdParams...)
	log.Info(cmd.String())
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Warn(cmd.String())
		log.Warnf("%+v, %v", string(out), err)
		return err
	}
	return nil
}
