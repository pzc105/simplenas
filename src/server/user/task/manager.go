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
}

func (ts *TasksIml) newTaskId() ptype.TaskId {
	return ptype.TaskId(ts.idPool.NewId())
}
