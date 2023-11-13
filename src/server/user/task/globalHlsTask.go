package task

import (
	"errors"
	"pnas/ptype"
	"pnas/video"
	"sync"
)

var (
	ErrExistedVideo    = errors.New("existed video")
	ErrFailed2NewVideo = errors.New("failed to new video")
)

type htask struct {
	hid       ptype.HlsTaskId
	callbacks []TaskCallback
}

type gHlsTask struct {
	hlsProcess video.HlsProcess

	mtx        sync.Mutex
	genHslTask map[ptype.VideoID]*htask
}

func (h *gHlsTask) init() {
	h.hlsProcess.Init()
	h.genHslTask = make(map[ptype.VideoID]*htask)
}

func (h *gHlsTask) callback(vid ptype.VideoID, err error) {
	h.mtx.Lock()
	if err == nil {
		video.VideoHasHls(vid)
	}
	task, ok := h.genHslTask[vid]
	if !ok {
		h.mtx.Unlock()
		panic("failed to get hls task")
	}
	delete(h.genHslTask, vid)
	h.mtx.Unlock()
	for _, c := range task.callbacks {
		c(err)
	}
}

type addHlsTaskParams struct {
	videoFullName string
	callback      TaskCallback
	audioTracksFN []string
	outDir        string
}

func (h *gHlsTask) add(params *addHlsTaskParams) (ptype.VideoID, error) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	v, err := video.GetVideoByFileName(params.videoFullName)
	if err != nil {
		vid, err := video.New(params.videoFullName)
		if err != nil {
			return -1, ErrFailed2NewVideo
		}
		v.Id = vid
	} else if v.HlsCreated {
		return -1, ErrExistedVideo
	}

	if r, ok := h.genHslTask[v.Id]; ok {
		r.callbacks = append(r.callbacks, params.callback)
		return v.Id, nil
	}

	hlsCallback := func(err error) {
		h.callback(v.Id, err)
	}

	hid, err := h.hlsProcess.Gen(&video.HlsGenParams{
		FullVideoFileName: params.videoFullName,
		FullAudioFileName: params.audioTracksFN,
		OutDir:            params.outDir,
		Callback:          hlsCallback,
	})
	if err != nil {
		return -1, errors.New("failed to create hls process")
	}
	h.genHslTask[v.Id] = &htask{
		hid:       hid,
		callbacks: []TaskCallback{params.callback},
	}
	return v.Id, nil
}
