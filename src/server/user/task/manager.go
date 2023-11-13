package task

import (
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

