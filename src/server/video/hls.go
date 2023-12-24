package video

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"pnas/log"
	"pnas/ptype"
	"pnas/setting"
	"sort"
	"strings"
)

type KV [2]string

type EncoderParams struct {
	W            int
	H            int
	Filters      string
	VCodecParams []KV
	ACodecParams []KV
}

func mapStreams(params *GenHlsOpts, cmdParams *[]string) (
	videoStreamCount int, audioLangInGroup []string, err error) {

	meta, err := GetMetadata(params.VideoFileName)
	if err != nil {
		return 0, []string{}, err
	}

	wantResolutions := params.WantedResolutions

	vStreams := GetStreams(meta.Streams, IsVideoStream)
	if len(vStreams) == 0 {
		return 0, []string{}, errors.New("no video stream")
	}
	sort.Slice(vStreams, func(i, j int) bool {
		return vStreams[i].Width > vStreams[j].Width
	})
	vStream := vStreams[0]
	var subRes []EncoderParams
	for _, ws := range wantResolutions {
		if ws.W <= vStream.Width {
			ws.H = -2
			subRes = append(subRes, ws)
		} else if ws.H <= vStream.Height {
			ws.W = -2
			subRes = append(subRes, ws)
		}
	}
	var buf bytes.Buffer
	for i, r := range subRes {
		*cmdParams = append(*cmdParams, "-map")
		*cmdParams = append(*cmdParams, fmt.Sprintf("%d:v:0", vStream.Index))
		if len(r.Filters) > 0 {
			*cmdParams = append(*cmdParams, fmt.Sprintf("-filter:v:%d", i))
			buf.WriteString(fmt.Sprintf(r.Filters, r.W, r.H))
		}
		*cmdParams = append(*cmdParams, buf.String())
		buf.Reset()
		for _, p := range r.VCodecParams {
			*cmdParams = append(*cmdParams, fmt.Sprintf("%s:v:%d", p[0], i))
			*cmdParams = append(*cmdParams, p[1])
		}
	}

	type pair struct {
		findex   int
		sindex   int
		channels int
		lang     string
	}

	var selectedAudioStreams []pair
	selectAudioStreams := func(i int, meta *Metadata) {
		aStreams := GetStreams(meta.Streams, IsAudioStream)
		for _, s := range aStreams {
			lang := strings.TrimSpace(s.Tags.Language)
			if len(lang) == 0 || lang == "und" {
				lang = fmt.Sprintf("audio%d", s.Index)
			}
			selectedAudioStreams = append(selectedAudioStreams,
				pair{
					findex:   i,
					sindex:   s.Index,
					channels: s.Channels,
					lang:     lang,
				})
			audioLangInGroup = append(audioLangInGroup, lang)
		}
	}
	selectAudioStreams(0, meta)
	for i, afn := range params.AudioFileNames {
		meta, err := GetMetadata(afn)
		if err != nil {
			continue
		}
		selectAudioStreams(i+1, meta)
	}

	aIndex := 0
	for _, r := range subRes {
		for _, s := range selectedAudioStreams {
			*cmdParams = append(*cmdParams, "-map")
			*cmdParams = append(*cmdParams, fmt.Sprintf("%d:%d", s.findex, s.sindex))

			for _, p := range r.ACodecParams {
				*cmdParams = append(*cmdParams, fmt.Sprintf("%s:a:%d", p[0], aIndex))
				*cmdParams = append(*cmdParams, p[1])
			}
			aIndex += 1
		}
	}
	return len(subRes), audioLangInGroup, nil
}

type GenHlsOpts struct {
	VideoFileName     string
	AudioFileNames    []string
	WantedResolutions []EncoderParams
	Global            []KV
	GlobalVideoParams []KV
	GlobalAudioParams []KV
	OutDir            string
	BaseUrl           string
	OnProcess         func(pid int)
}

func GenHls(params *GenHlsOpts) error {
	var cmdParams []string
	cmdParams = append(cmdParams, "-hide_banner")
	cmdParams = append(cmdParams, "-loglevel")
	cmdParams = append(cmdParams, "warning")
	for _, p := range params.Global {
		cmdParams = append(cmdParams, p[0])
		cmdParams = append(cmdParams, p[1])
	}
	cmdParams = append(cmdParams, "-i")
	cmdParams = append(cmdParams, params.VideoFileName)
	for _, afn := range params.AudioFileNames {
		cmdParams = append(cmdParams, "-i")
		cmdParams = append(cmdParams, afn)
	}
	for _, p := range params.GlobalVideoParams {
		cmdParams = append(cmdParams, p[0]+":v")
		cmdParams = append(cmdParams, p[1])
	}
	for _, p := range params.GlobalAudioParams {
		cmdParams = append(cmdParams, p[0]+":a")
		cmdParams = append(cmdParams, p[1])
	}
	videoStreamCount, audioLangInGroup, err := mapStreams(params, &cmdParams)
	if err != nil {
		return err
	}

	cmdParams = append(cmdParams, "-f")
	cmdParams = append(cmdParams, "hls")
	if len(params.BaseUrl) > 0 {
		cmdParams = append(cmdParams, "-hls_base_url")
		cmdParams = append(cmdParams, params.BaseUrl)
	}
	var sms []string
	for vi := 0; vi < videoStreamCount; vi++ {
		for ai, lang := range audioLangInGroup {
			sms = append(sms, fmt.Sprintf("a:%d,agroup:ag%d,name:a%d,language:%s",
				vi*len(audioLangInGroup)+ai, vi, vi*len(audioLangInGroup)+ai, lang))
		}
		str := fmt.Sprintf("v:%d,agroup:ag%d,name:v%d", vi, vi, vi)
		sms = append(sms, str)
	}
	cmdParams = append(cmdParams, "-var_stream_map")
	cmdParams = append(cmdParams, strings.Join(sms, " "))
	cmdParams = append(cmdParams, "-hls_time")
	cmdParams = append(cmdParams, "10")
	cmdParams = append(cmdParams, "-hls_flags")
	cmdParams = append(cmdParams, "independent_segments")
	cmdParams = append(cmdParams, "-hls_segment_type")
	cmdParams = append(cmdParams, "mpegts")
	cmdParams = append(cmdParams, "-hls_playlist_type")
	cmdParams = append(cmdParams, "event")
	cmdParams = append(cmdParams, "-hls_list_size")
	cmdParams = append(cmdParams, "0")

	cmdParams = append(cmdParams, "-master_pl_name")
	cmdParams = append(cmdParams, "master.m3u8")
	cmdParams = append(cmdParams, "-hls_segment_filename")
	cmdParams = append(cmdParams, params.OutDir+"/stream_%v/data%03d.ts")
	cmdParams = append(cmdParams, params.OutDir+"/stream_%v/stream.m3u8")
	cmdParams = append(cmdParams, "-y")

	cmd := exec.Command("ffmpeg", cmdParams...)
	log.Info(cmd.String())

	var stdBuffer bytes.Buffer
	cmd.Stdout = &stdBuffer
	cmd.Stderr = &stdBuffer
	err = cmd.Start()
	if err != nil {
		out := stdBuffer.Bytes()
		log.Warnf("%+v, %v", string(out), err)
		return err
	}
	pid := cmd.Process.Pid
	if params.OnProcess != nil {
		params.OnProcess(pid)
	}
	err = cmd.Wait()
	if err != nil {
		out := stdBuffer.Bytes()
		log.Warnf("%+v, %v", string(out), err)
		return err
	}

	return nil
}

func GetHlsPlayListPath(vid ptype.VideoID) string {
	return setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%d", vid)
}
