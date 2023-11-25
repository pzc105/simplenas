package bt

import (
	"context"
	"io"
	"net/http"
	"pnas/setting"
	"strings"
	"sync"
	"time"
)

type trackerSub struct {
	mtx      sync.Mutex
	trackers []string

	shutDownCtx context.Context
	closeFunc   context.CancelFunc
	wg          sync.WaitGroup
}

func (t *trackerSub) Init() {
	t.shutDownCtx, t.closeFunc = context.WithCancel(context.Background())

	t.wg.Add(1)
	go t.background()
}

func (t *trackerSub) Close() {
	t.closeFunc()
	t.wg.Wait()
}

func (t *trackerSub) GetTrackers() []string {
	ret := []string{}
	t.mtx.Lock()
	ret = append(ret, t.trackers...)
	t.mtx.Unlock()
	return ret
}

func (t *trackerSub) refreshTrackers() {
	rsp, err := http.Get(setting.GS().Bt.TrackerSub)
	if err != nil {
		return
	}
	body, err := io.ReadAll(rsp.Body)
	if err != nil {
		return
	}
	lines := strings.Split(string(body), "\n")
	trackers := []string{}
	for _, line := range lines {
		if len(line) == 0 {
			continue
		}
		line = strings.Trim(line, " \t")
		if len(line) == 0 {
			continue
		}
		trackers = append(trackers, line)
	}
	t.mtx.Lock()
	t.trackers = trackers
	t.mtx.Unlock()
}

func (t *trackerSub) background() {
	defer t.wg.Done()
	t.refreshTrackers()
	ticker := time.NewTicker(time.Hour * 1)
loop:
	for {
		select {
		case <-ticker.C:
			t.refreshTrackers()
		case <-t.shutDownCtx.Done():
			break loop
		}
	}
}
