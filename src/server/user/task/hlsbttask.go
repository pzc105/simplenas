package task

import (
	"pnas/bt"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/setting"
)

type BtHlsTask struct {
	RawTask
	magnetUri string
	userId    ptype.UserID
	parentId  ptype.CategoryID
	tid       ptype.TorrentID
	st        prpc.BtStateEnum
	bt        bt.UserTorrents
	infoHash  *bt.InfoHash
}

type newBtHlsTaskParams struct {
	userId   ptype.UserID
	parentId ptype.CategoryID
	tasks    *TasksIml
	callback TaskCallback
	bt       bt.UserTorrents
}

func newBtHlsTask(params *newBtHlsTaskParams) ITask {
	var bd BtHlsTask
	bd.userId = params.userId
	bd.parentId = params.parentId
	bd.id = params.tasks.newTaskId()
	bd.callback = bd.btCallback
	bd.status.Store(TaskStatusIniting)
	bd.bt = params.bt
	return &bd
}

func (bd *BtHlsTask) btCallback(error) {
	bd.bt.SetTaskCallback(&bt.SetTaskCallbackParams{
		UserId:    bd.userId,
		TaskId:    bd.id,
		TorrentId: bd.tid,
		Callback:  nil,
	})
}

func (bd *BtHlsTask) Start() {
	res, err := bd.bt.Download(&bt.DownloadParams{
		UserId: bd.userId,
		Req: &prpc.DownloadRequest{
			Type:     prpc.DownloadRequest_MagnetUri,
			Content:  []byte(bd.magnetUri),
			SavePath: setting.GS().Bt.SavePath,
		},
	})
	if err != nil {
		bd.into(TaskStatusInitFailed, err)
		return
	}
	infoHash := bt.TranInfoHash(res.InfoHash)
	t, err := bd.bt.GetTorrent(infoHash)
	if err != nil {
		bd.into(TaskStatusInitFailed, err)
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

func (bd *BtHlsTask) onBtStatus(err error, st *prpc.TorrentStatus) {
	if err != nil {
		bd.into(TaskStatusInitFailed, err)
		return
	}
	if !bt.IsDownloadAll(bd.st) && bt.IsDownloadAll(st.State) {

	}
	bd.st = st.State
}
