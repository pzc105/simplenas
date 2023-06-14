package utils

import (
	"sync"
	"sync/atomic"
)

type tqClientOpts struct {
	max int32
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
			opts.max = max
		},
	}
}

func defaultTqOpts(opts *tqClientOpts) {
	opts.max = 1
}

type TaskQueue struct {
	opts  tqClientOpts
	queue chan func()
	used  atomic.Int32

	wg sync.WaitGroup
}

func (tq *TaskQueue) Init(opts ...TqClientOpt) {
	defaultTqOpts(&tq.opts)
	for _, opt := range opts {
		opt.apply(&tq.opts)
	}

	tq.used.Store(0)
	tq.queue = make(chan func(), tq.opts.max)

	tq.wg.Add(1)
	go tq.handle()
}

func (tq *TaskQueue) Close() {
	close(tq.queue)
	tq.wg.Wait()
}

func (tq *TaskQueue) TryPut(f func()) bool {
	for {
		v := tq.used.Load()
		if v >= tq.opts.max {
			return false
		}
		if tq.used.CompareAndSwap(v, v+1) {
			tq.queue <- f
			return true
		}
	}
}

func (tq *TaskQueue) Put(f func()) {
	FetchAndAdd(&tq.used, int32(1))
	tq.queue <- f
}

func (tq *TaskQueue) handle() {
	defer tq.wg.Add(-1)
	for {
		f, ok := <-tq.queue
		if !ok {
			return
		}
		FetchAndAdd(&tq.used, int32(-1))
		f()
	}
}
