package task

import (
	"pnas/bt"
	"pnas/category"
	"pnas/prpc"
	"pnas/ptype"
	"pnas/utils"
	"sync"
)

type TasksIml struct {
	mtx    sync.Mutex
	tasks  map[ptype.TaskId]ITask
	idPool utils.IdPool

	hlsTasks gHlsTask
}

func (ts *TasksIml) newTaskId() ptype.TaskId {
	return ptype.TaskId(ts.idPool.NewId())
}

func (ts *TasksIml) Init() {
	ts.idPool.Init()
	ts.tasks = make(map[ptype.TaskId]ITask)
	ts.hlsTasks.init()
}

type NewBtHlsParams struct {
	UserId           ptype.UserID
	ParentId         ptype.CategoryID
	Bt               bt.UserTorrents
	DownloadReq      *prpc.DownloadRequest
	CategorySer      category.IService
	RecursiveNewPath bool
}

func (ts *TasksIml) NewVideoTask(params *NewVideoTaskParams) error {
	params.mgr = ts
	params.ghlsTask = &ts.hlsTasks
	task := newVideoTask(params)

	ts.mtx.Lock()
	ts.tasks[task.GetId()] = task
	ts.mtx.Unlock()
	task.Start()
	return nil
}

func (ts *TasksIml) NewBtHlsTask(params *NewBtHlsParams) error {
	task := newBtHlsTask(&newBtHlsTaskParams{
		userId:           params.UserId,
		parentId:         params.ParentId,
		bt:               params.Bt,
		downloadReq:      params.DownloadReq,
		categorySer:      params.CategorySer,
		mgr:              ts,
		ghlsTask:         &ts.hlsTasks,
		recursiveNewPath: params.RecursiveNewPath,
	})

	ts.mtx.Lock()
	ts.tasks[task.GetId()] = task
	ts.mtx.Unlock()
	task.Start()
	return nil
}
