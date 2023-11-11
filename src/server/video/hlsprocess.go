package video

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"pnas/log"
	"pnas/utils"
	"sync"
	"time"
)

var (
	ErrFailed = errors.New("failed to create hls task")
)

const (
	EnableQsv    = false
	CudaType     = 1
	QsvType      = 2
	SoftwareType = 3
)

type HlsTaskId int64
type HlsCallback func(error)

type hlsTask struct {
	id        HlsTaskId
	queueType int
	qtaskId   utils.TaskId
	params    *HlsGenParams
	cancel    bool
	mtx       sync.Mutex
	pid       int
	callback  HlsCallback
}

type HlsProcess struct {
	cudaQueue utils.TaskQueue
	qsvQueue  utils.TaskQueue
	soQueue   utils.TaskQueue

	mtx   sync.Mutex
	tasks map[HlsTaskId]*hlsTask

	idPool utils.IdPool

	shutDownCtx context.Context
	closeFunc   context.CancelFunc
	wg          sync.WaitGroup
}

type OnGenHlsCallback func(error)

func (h *HlsProcess) Init() {
	h.cudaQueue.Init(utils.WithMaxQueue(1024))
	h.qsvQueue.Init(utils.WithMaxQueue(1024))
	h.soQueue.Init(utils.WithMaxQueue(1024))

	h.idPool.Init()
	h.tasks = make(map[HlsTaskId]*hlsTask)

	h.shutDownCtx, h.closeFunc = context.WithCancel(context.Background())

	h.wg.Add(1)
	go func() {
		defer h.wg.Add(-1)
		ticker := time.NewTicker(time.Second * 3)
	loop:
		for {
			select {
			case <-ticker.C:
				if h.soQueue.Idle() {
					task := h.cudaQueue.Steal()
					if task != nil {
						htask, _ := task.Identity.(*hlsTask)
						log.Debugf("steal task: %s", htask.params.FullVideoFileName)
						h.useSoft(htask)
					}
				}
			case <-h.shutDownCtx.Done():
				break loop
			}
		}
	}()
}

type HlsGenParams struct {
	FullVideoFileName string
	FullAudioFileName []string
	OutDir            string
	Callback          HlsCallback
}

func (h *HlsProcess) onCallback(task *hlsTask, err error) {
	if err == nil {
		h.mtx.Lock()
		delete(h.tasks, task.id)
		h.mtx.Unlock()
		task.callback(nil)
		return
	}

	if task.cancel {
		task.callback(errors.New("failed to gen hsl"))
		return
	}

	if task.queueType == CudaType {
		if EnableQsv {
			h.useQsv(task)
		} else {
			h.useSoft(task)
		}
		return
	}

	if task.queueType == QsvType {
		h.useSoft(task)
		return
	}

	if task.queueType == SoftwareType {
		task.callback(errors.New("failed to gen hsl"))
		return
	}
}

func (h *HlsProcess) onHlsProcess(task *hlsTask, pid int) {
	task.mtx.Lock()
	task.pid = pid
	task.mtx.Unlock()
}

func (h *HlsProcess) Stop(id HlsTaskId) {
	h.mtx.Lock()
	task, ok := h.tasks[id]
	if ok {
		delete(h.tasks, id)
	}
	h.mtx.Unlock()
	if !ok {
		return
	}
	task.mtx.Lock()
	defer task.mtx.Unlock()
	task.cancel = true
	switch task.queueType {
	case CudaType:
		h.cudaQueue.Remove(task.qtaskId)
	case SoftwareType:
		h.soQueue.Remove(task.qtaskId)
	case QsvType:
		h.qsvQueue.Remove(task.qtaskId)
	}
	if task.pid > 0 {
		exec.Command("kill", "-9", fmt.Sprint(task.pid)).Run()
	}
}

func (h *HlsProcess) Gen(params *HlsGenParams) (HlsTaskId, error) {
	task := &hlsTask{
		id:        HlsTaskId(h.idPool.NewId()),
		queueType: CudaType,
		params:    params,
		cancel:    false,
		callback:  params.Callback,
		pid:       -1,
	}

	task.mtx.Lock()
	defer task.mtx.Unlock()

	onProcess := func(pid int) {
		h.onHlsProcess(task, pid)
	}

	tid, err := h.cudaQueue.TryPutWithIdentity(func() {
		err := GenHls(
			&GenHlsOpts{
				VideoFileName:     params.FullVideoFileName,
				AudioFileNames:    params.FullAudioFileName,
				WantedResolutions: CudaSplitEncoderParams,
				OutDir:            params.OutDir,
				Global:            CudaGlobalDecode,
				GlobalVideoParams: CudaH264GlobalVideoParams,
				GlobalAudioParams: GlobalAudioParams,
				OnProcess:         onProcess,
			})
		if err != nil {
			err = GenHls(
				&GenHlsOpts{
					VideoFileName:     params.FullVideoFileName,
					AudioFileNames:    params.FullAudioFileName,
					WantedResolutions: CudaEncoderParams2,
					OutDir:            params.OutDir,
					Global:            CudaGlobalDecode2,
					GlobalVideoParams: CudaH264GlobalVideoParams,
					GlobalAudioParams: GlobalAudioParams,
					OnProcess:         onProcess,
				})
		}
		h.onCallback(task, err)
	}, task)

	if err != nil {
		return -1, ErrFailed
	}
	task.qtaskId = tid

	h.mtx.Lock()
	h.tasks[task.id] = task
	h.mtx.Unlock()

	return task.id, nil
}

func (h *HlsProcess) useQsv(task *hlsTask) {
	var tid utils.TaskId
	var err error

	defer func() {
		if err != nil {
			h.onCallback(task, err)
		}
	}()

	task.mtx.Lock()
	defer task.mtx.Unlock()

	if task.cancel {
		return
	}

	onProcess := func(pid int) {
		h.onHlsProcess(task, pid)
	}

	task.queueType = QsvType
	params := task.params
	tid, err = h.soQueue.TryPut(func() {
		err := GenHls(
			&GenHlsOpts{
				VideoFileName:     params.FullVideoFileName,
				AudioFileNames:    params.FullAudioFileName,
				WantedResolutions: QsvSplitEncoderParams,
				Global:            QsvGlobalDecode,
				OutDir:            params.OutDir,
				GlobalVideoParams: QsvH264GlobalVideoParams,
				GlobalAudioParams: GlobalAudioParams,
				OnProcess:         onProcess,
			})
		h.onCallback(task, err)
	})
	if err != nil {
		task.cancel = true
		return
	}
	task.qtaskId = tid
}

func (h *HlsProcess) useSoft(task *hlsTask) {

	var tid utils.TaskId
	var err error

	defer func() {
		if err != nil {
			h.onCallback(task, err)
		}
	}()

	task.mtx.Lock()
	defer task.mtx.Unlock()

	if task.cancel {
		return
	}

	onProcess := func(pid int) {
		h.onHlsProcess(task, pid)
	}

	task.queueType = SoftwareType
	params := task.params
	tid, err = h.soQueue.TryPut(func() {
		err := GenHls(
			&GenHlsOpts{
				VideoFileName:     params.FullVideoFileName,
				AudioFileNames:    params.FullAudioFileName,
				WantedResolutions: SoSplitEncoderParams,
				Global:            SoGlobalDecode,
				OutDir:            params.OutDir,
				GlobalVideoParams: SoH264GlobalVideoParams,
				GlobalAudioParams: GlobalAudioParams,
				OnProcess:         onProcess,
			})
		h.onCallback(task, err)
	})
	if err != nil {
		task.cancel = true
		return
	}
	task.qtaskId = tid
}
