package utils

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
)

type TaskFunc func()

type tqClientOpts struct {
	maxTasks int32
}

type TqClientOpt interface {
	apply(*tqClientOpts)
}

type funcTqClientOpt struct {
	do func(opts *tqClientOpts)
}

func (f *funcTqClientOpt) apply(opts *tqClientOpts) {
	f.do(opts)
}

func WithMaxQueue(max int32) *funcTqClientOpt {
	return &funcTqClientOpt{
		do: func(opts *tqClientOpts) {
			opts.maxTasks = max
		},
	}
}

func defaultTqOpts(opts *tqClientOpts) {
	opts.maxTasks = -1
}

var (
	ErrNullTask   = errors.New("task is null")
	ErrOutMaxTask = errors.New("exhausted tasks")
	ErrNotFound   = errors.New("not found task")
	ErrClosed     = errors.New("shut down task")
)

type TaskId int64

type CloseWay int

const (
	Running           = 0
	CloseWayImmediate = 1
	CloseWayDrained   = 2
)

type Task struct {
	F        TaskFunc
	Id       TaskId
	Identity interface{}
}

type TaskQueue struct {
	opts tqClientOpts

	idPool IdPool

	cond     sync.Cond
	mtx      sync.Mutex
	tasks    list.List
	tNum     atomic.Int32
	tasksMap map[TaskId]*list.Element
	shutDown CloseWay

	wg sync.WaitGroup
}

func (tq *TaskQueue) Init(opts ...TqClientOpt) {
	defaultTqOpts(&tq.opts)
	for _, opt := range opts {
		opt.apply(&tq.opts)
	}

	tq.idPool.Init()
	tq.shutDown = Running
	tq.cond = *sync.NewCond(&tq.mtx)
	tq.tasks.Init()
	tq.tNum.Store(0)
	tq.tasksMap = make(map[TaskId]*list.Element)

	tq.wg.Add(1)
	go tq.handle()
}

func (tq *TaskQueue) Close(way CloseWay) {
	tq.mtx.Lock()
	tq.shutDown = way
	tq.cond.Signal()
	tq.mtx.Unlock()
	tq.wg.Wait()
}

func (tq *TaskQueue) put(f TaskFunc, identity interface{}) (TaskId, error) {
	if f == nil {
		return -1, ErrNullTask
	}
	id := TaskId(tq.idPool.NewId())

	tq.mtx.Lock()

	if tq.shutDown != Running {
		tq.mtx.Unlock()
		return -1, ErrClosed
	}

	for tq.opts.maxTasks > 0 && tq.tNum.Load() >= tq.opts.maxTasks {
		tq.cond.Wait()
	}
	e := tq.tasks.PushBack(&Task{
		F:        f,
		Id:       id,
		Identity: identity,
	})
	tq.tasksMap[id] = e
	tq.cond.Signal()
	tq.tNum.Add(1)
	tq.mtx.Unlock()
	return id, nil
}

func (tq *TaskQueue) Put(f TaskFunc) (TaskId, error) {
	return tq.put(f, nil)
}

func (tq *TaskQueue) PutWithIdentity(f TaskFunc, identity interface{}) (TaskId, error) {
	return tq.put(f, identity)
}

func (tq *TaskQueue) tryPut(f TaskFunc, identity interface{}) (TaskId, error) {
	if f == nil {
		return -1, errors.New("task is null")
	}
	if tq.opts.maxTasks > 0 && tq.tNum.Load() >= tq.opts.maxTasks {
		return -1, ErrOutMaxTask
	}

	id := TaskId(tq.idPool.NewId())
	tq.mtx.Lock()

	if tq.shutDown != Running {
		tq.mtx.Unlock()
		return -1, ErrClosed
	}

	if tq.opts.maxTasks > 0 && tq.tNum.Load() >= tq.opts.maxTasks {
		tq.mtx.Unlock()
		return -1, ErrOutMaxTask
	}
	e := tq.tasks.PushBack(&Task{
		F:        f,
		Id:       id,
		Identity: identity,
	})
	tq.tasksMap[id] = e
	tq.cond.Signal()
	tq.tNum.Add(1)
	tq.mtx.Unlock()
	return id, nil
}

func (tq *TaskQueue) TryPut(f TaskFunc) (TaskId, error) {
	return tq.tryPut(f, nil)
}

func (tq *TaskQueue) TryPutWithIdentity(f TaskFunc, identity interface{}) (TaskId, error) {
	return tq.tryPut(f, identity)
}

func (tq *TaskQueue) Remove(id TaskId) error {
	tq.mtx.Lock()
	e, ok := tq.tasksMap[id]
	if !ok {
		tq.mtx.Unlock()
		return ErrNotFound
	}
	tq.tasks.Remove(e)
	delete(tq.tasksMap, id)
	tq.tNum.Add(-1)
	tq.cond.Signal()
	tq.mtx.Unlock()
	return nil
}

func (tq *TaskQueue) Idle() bool {
	tq.mtx.Lock()
	ret := tq.shutDown == Running && tq.tNum.Load() == 0
	tq.mtx.Unlock()
	return ret
}

func (tq *TaskQueue) Len() int32 {
	tq.mtx.Lock()
	ret := tq.tNum.Load()
	tq.mtx.Unlock()
	return ret
}

func (tq *TaskQueue) Steal() *Task {
	tq.mtx.Lock()
	if tq.tasks.Len() == 0 {
		tq.mtx.Unlock()
		return nil
	}
	e := tq.tasks.Back()
	tq.tasks.Remove(e)
	t, ok := e.Value.(*Task)
	if !ok {
		panic("unknown task func")
	}
	delete(tq.tasksMap, t.Id)
	tq.mtx.Unlock()
	return t
}

func (tq *TaskQueue) handle() {
	defer tq.wg.Add(-1)
	for {
		tq.mtx.Lock()
		for tq.tasks.Len() == 0 && tq.shutDown == Running {
			tq.cond.Wait()
		}

		if tq.shutDown == CloseWayImmediate ||
			tq.tasks.Len() == 0 && tq.shutDown == CloseWayDrained {
			tq.mtx.Unlock()
			break
		}

		e := tq.tasks.Front()
		tq.tasks.Remove(e)
		t, ok := e.Value.(*Task)
		if !ok {
			panic("unknown task func")
		}
		delete(tq.tasksMap, t.Id)
		tq.mtx.Unlock()
		t.F()
		tq.tNum.Add(-1)
	}
}
