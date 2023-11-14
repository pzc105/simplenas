package task

import (
	"pnas/bt"
	"pnas/category"
	"pnas/prpc"
	"pnas/ptype"
	"sort"
	"sync/atomic"
)

type btHlsTask struct {
	RawTask
	downloadReq      *prpc.DownloadRequest
	userId           ptype.UserID
	parentId         ptype.CategoryID
	mgr              ITasks
	st               prpc.BtStateEnum
	bt               bt.UserTorrents
	infoHash         *bt.InfoHash
	tid              ptype.TorrentID
	ghlsTasks        *gHlsTask
	videoTaskCount   atomic.Int32
	categorySer      category.IService
	recursiveNewPath bool
	videoTasks       []ITask
}

type newBtHlsTaskParams struct {
	callback         TaskCallback
	userId           ptype.UserID
	parentId         ptype.CategoryID
	mgr              ITasks
	bt               bt.UserTorrents
	downloadReq      *prpc.DownloadRequest
	ghlsTask         *gHlsTask
	categorySer      category.IService
	recursiveNewPath bool
}

func newBtHlsTask(params *newBtHlsTaskParams) ITask {
	var bd btHlsTask
	bd.RawTask.init(params.mgr)
	bd.userId = params.userId
	bd.parentId = params.parentId
	bd.bt = params.bt
	bd.downloadReq = params.downloadReq
	bd.ghlsTasks = params.ghlsTask
	bd.categorySer = params.categorySer
	bd.mgr = params.mgr
	bd.callback = params.callback
	bd.recursiveNewPath = params.recursiveNewPath
	return &bd
}

func (bd *btHlsTask) videoTaskCallback(err error) {
	bd.videoTaskCount.Add(-1)
	if bd.videoTaskCount.Load() == 0 {
		bd.into(TaskStatusFinished, nil)
	}
}

func (bd *btHlsTask) Start() {
	bd.downloadReq.StopAfterGotMeta = false
	res, err := bd.bt.Download(&bt.DownloadParams{
		UserId: bd.userId,
		Req:    bd.downloadReq,
	})
	if err != nil {
		if err == bt.ErrDownloaded {
			bd.infoHash = bt.TranInfoHash(res.InfoHash)
			bd.downloaded()
			return
		}
		bd.into(TaskStatusFailed, err)
		return
	}
	infoHash := bt.TranInfoHash(res.InfoHash)
	t, err := bd.bt.GetTorrent(infoHash)
	if err != nil {
		bd.into(TaskStatusFailed, err)
	}
	bd.infoHash = infoHash
	bd.tid = t.GetBaseInfo().Id
	bd.bt.SetTaskCallback(&bt.SetTaskCallbackParams{
		UserId:    bd.userId,
		TaskId:    bd.id,
		TorrentId: bd.tid,
		Callback:  bd.onBtStatus,
	})
}

func (bd *btHlsTask) downloaded() {
	bd.bt.SetTaskCallback(&bt.SetTaskCallbackParams{
		UserId:    bd.userId,
		TaskId:    bd.id,
		TorrentId: bd.tid,
		Callback:  nil,
	})
	t, err := bd.bt.GetTorrent(bd.infoHash)
	if err != nil {
		bd.into(TaskStatusFailed, err)
		return
	}
	files := t.GetFiles()
	sort.Slice(files, func(a, b int) bool {
		return files[a].Name < files[b].Name
	})
	for _, f := range files {
		if f.FileType&bt.FileVideoType != 0 {
			task := newVideoTask(&NewVideoTaskParams{
				Callback:    bd.videoTaskCallback,
				UserId:      bd.userId,
				ParentId:    bd.parentId,
				mgr:         bd.mgr,
				ghlsTask:    bd.ghlsTasks,
				CategorySer: bd.categorySer,
				Bt:          bd.bt,
				InfoHash:    bd.infoHash,
				BtFileIndex: int(f.Index),
			})
			bd.videoTaskCount.Add(1)
			bd.videoTasks = append(bd.videoTasks, task)
		}
	}
	if len(bd.videoTasks) == 0 {
		bd.into(TaskStatusFinished, nil)
	}
	for _, task := range bd.videoTasks {
		task.Start()
	}
}

func (bd *btHlsTask) onBtStatus(err error, st *prpc.TorrentStatus) {
	if err != nil {
		bd.into(TaskStatusFailed, err)
		bd.bt.SetTaskCallback(&bt.SetTaskCallbackParams{
			UserId:    bd.userId,
			TaskId:    bd.id,
			TorrentId: bd.tid,
			Callback:  nil,
		})
		return
	}
	if !bt.IsDownloadAll(bd.st) && bt.IsDownloadAll(st.State) {
		bd.downloaded()
	}
	bd.st = st.State
}
