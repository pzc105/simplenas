package task

import (
	"pnas/ptype"
	"sync/atomic"
)

type TaskStatus int

const (
	TaskStatusIniting = iota
	TaskStatusRunning
	TaskStatusFinished
	TaskStatusInitFailed
)

type TaskCallback func(error)

type ITask interface {
	Start()
	GetStatus() TaskStatus
	GetProgress() int32
	Stop() error
}

type RawTask struct {
	id       ptype.TaskId
	status   atomic.Value
	progress atomic.Int32
	callback TaskCallback
}

func (r *RawTask) into(st TaskStatus, err error) {
	r.status.Store(st)
	if st == TaskStatusInitFailed && r.callback != nil {
		r.callback(err)
	} else if st == TaskStatusFinished && r.callback != nil {
		r.callback(err)
	}
}

func (r *RawTask) GetStatus() TaskStatus {
	return r.status.Load().(TaskStatus)
}

func (r *RawTask) GetProgress() int32 {
	return r.progress.Load()
}

func (r *RawTask) Stop() error {
	return nil
}
