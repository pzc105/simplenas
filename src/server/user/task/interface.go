package task

import (
	"pnas/ptype"
	"sync/atomic"
)

type TaskStatus int

const (
	TaskStatusIniting TaskStatus = iota
	TaskStatusRunning
	TaskStatusFinished
	TaskStatusFailed
	TaskStatusStopped
)

func isEndStatus(st TaskStatus) bool {
	return st == TaskStatusFinished || st == TaskStatusFailed || st == TaskStatusStopped
}

type TaskCallback func(error)

type ITask interface {
	GetId() ptype.TaskId
	Start()
	GetStatus() TaskStatus
	GetProgress() int32
	Stop() error
}

type ITasks interface {
	NewBtHlsTask(params *NewBtHlsParams) error
	NewVideoTask(params *NewVideoTaskParams) error
	newTaskId() ptype.TaskId
}

type RawTask struct {
	id       ptype.TaskId
	status   atomic.Value
	progress atomic.Int32
	callback TaskCallback
}

func (r *RawTask) init(mgr ITasks) {
	r.id = mgr.newTaskId()
	r.status.Store(TaskStatusIniting)
}

func (r *RawTask) GetId() ptype.TaskId {
	return r.id
}

func (r *RawTask) into(st TaskStatus, err error) {
	old := r.status.Load()
	if isEndStatus(old.(TaskStatus)) {
		return
	}
	for !r.status.CompareAndSwap(old, st) {
		old = r.status.Load()
		if isEndStatus(old.(TaskStatus)) {
			return
		}
	}
	if r.callback != nil && isEndStatus(r.status.Load().(TaskStatus)) {
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
