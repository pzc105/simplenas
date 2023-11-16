package task

import (
	"errors"
	"fmt"
	"pnas/ptype"
	"pnas/setting"
	"pnas/video"
	"sync"
)

var (
	ErrExistedVideo    = errors.New("existed video")
	ErrFailed2NewVideo = errors.New("failed to new video")
)

type htask struct {
	hid       ptype.HlsTaskId
	callbacks map[ptype.TaskId]TaskCallback
}

type gHlsTask struct {
	HlsProcess video.HlsProcess

	mtx        sync.Mutex
	genHslTask map[ptype.VideoID]*htask
}

func (h *gHlsTask) init() {
	h.HlsProcess.Init()
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

type AddHlsTaskParams struct {
	VideoFullName string
	AudioTracksFN []string
	MyTaskId      ptype.TaskId
	Callback      TaskCallback
}

func (h *gHlsTask) add(params *AddHlsTaskParams) (ptype.VideoID, error) {
	h.mtx.Lock()
	defer h.mtx.Unlock()

	v, err := video.GetVideoByFileName(params.VideoFullName)
	if err != nil {
		vid, err := video.New(params.VideoFullName)
		if err != nil {
			return -1, ErrFailed2NewVideo
		}
		v.Id = vid
	} else if v.HlsCreated {
		return v.Id, ErrExistedVideo
	}

	if r, ok := h.genHslTask[v.Id]; ok {
		if params.MyTaskId <= 0 {
			return v.Id, nil
		}
		if _, ok := r.callbacks[params.MyTaskId]; ok {
			panic("duplicate")
		}
		if params.Callback != nil {
			r.callbacks[params.MyTaskId] = params.Callback
		}
		return v.Id, nil
	}

	hlsCallback := func(err error) {
		h.callback(v.Id, err)
	}

	outDir := setting.GS().Server.HlsPath + fmt.Sprintf("/vid_%d", v.Id)

	hid, err := h.HlsProcess.Gen(&video.HlsGenParams{
		FullVideoFileName: params.VideoFullName,
		FullAudioFileName: params.AudioTracksFN,
		OutDir:            outDir,
		Callback:          hlsCallback,
	})
	if err != nil {
		return -1, errors.New("failed to create hls process")
	}
	ht := &htask{
		hid: hid,
	}
	ht.callbacks = make(map[ptype.TaskId]TaskCallback)
	if params.Callback != nil {
		ht.callbacks[params.MyTaskId] = params.Callback
	}
	h.genHslTask[v.Id] = ht
	return v.Id, nil
}
