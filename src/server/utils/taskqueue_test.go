package utils

import (
	"math"
	"math/rand"
	"runtime"
	"sync"
	"testing"
	"time"
)

func TestTaskQUeue1(t *testing.T) {
	var tq TaskQueue
	tq.Init()

	var inc int
	var mtx sync.Mutex
	flag := make(map[int]bool)
	var wg sync.WaitGroup

	var bg = runtime.GOMAXPROCS(0)
	const incNumPer = 1000000

	for i := 0; i < bg; i++ {
		wg.Add(1)
		go func() {
			defer wg.Add(-1)
			tq.Put(func() {
				for j := 0; j < incNumPer; j++ {
					mtx.Lock()
					flag[inc] = true
					inc += 1
					mtx.Unlock()
				}
			})
		}()
	}
	wg.Wait()
	tq.Close(CloseWayDrained)

	if inc != bg*incNumPer {
		t.Errorf("not match inc: %d goal: %d", inc, bg*incNumPer)
	}

	for i := 0; i < inc; i++ {
		if _, ok := flag[i]; !ok {
			t.Errorf("not found %d", i)
			break
		}
	}
}

func TestTaskQUeue2(t *testing.T) {

	for i := 0; i < 100; i++ {
		var tq TaskQueue
		tq.Init()

		var f uint64
		var mtx sync.Mutex
		speculate := uint64(math.MaxUint64)

		for i := 0; i < 64; i++ {
			j := i
			id, _ := tq.Put(func() {
				mtx.Lock()
				f |= (1 << j)
				mtx.Unlock()
			})

			if rand.Int()%2 == 0 {
				time.Sleep(time.Microsecond * 2)
			}
			if tq.Remove(id) == nil {
				mtx.Lock()
				speculate &= ^uint64(1 << j)
				mtx.Unlock()
			}
		}
		tq.Close(CloseWayDrained)

		if f != speculate {
			t.Errorf("%d not match %d", f, speculate)
		}
	}
}

func TestTaskQUeue3(t *testing.T) {
	var tq TaskQueue
	tq.Init(WithMaxQueue(3))
	for i := 0; i < 3; i++ {
		_, err := tq.TryPut(func() {})
		if err != nil {
			t.Errorf("failed")
		}
	}
	tq.Close(CloseWayImmediate)
	_, err := tq.Put(func() {})
	if err != ErrClosed {
		t.Errorf("failed")
	}
}

func TestTaskQUeue4(t *testing.T) {
	var tq TaskQueue
	var mtx sync.Mutex
	mtx.Lock()
	tq.Init(WithMaxQueue(3))
	for i := 0; i < 100; i++ {
		_, err := tq.TryPut(func() {
			mtx.Lock()
		})
		if i == 99 && err != ErrOutMaxTask {
			t.Errorf("failed")
		}
	}
	tq.Close(CloseWayImmediate)
	mtx.Unlock()
}
