package chat

import (
	"context"
	"pnas/log"
	"pnas/user"
	"pnas/utils"
	"sync"
	"sync/atomic"
	"time"
)

const (
	HistoryMsgMaxCount = uint64(10)
)

type userData struct {
	id               int64
	sessionId        int64
	sendFunc         SendFunc
	nextReadPos      atomic.Uint64
	maxCacheNum      uint64
	maxCacheDuration time.Duration

	mtx      sync.Mutex
	lastPush time.Time
}

func (ud *userData) update(now time.Time) {
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	ud.lastPush = now
}

func (ud *userData) elapse(now time.Time) time.Duration {
	ud.mtx.Lock()
	defer ud.mtx.Unlock()
	return now.Sub(ud.lastPush)
}

type ChatMessage struct {
	UserId   user.ID
	SentTime time.Time
	Msg      string
}

type versionChatMessage struct {
	version uint64
	msg     *ChatMessage
}

type ChatRoomImpl struct {
	IRoom

	roomKey       string
	immediatePush bool
	interval      time.Duration

	mtx       sync.Mutex
	usersData map[int64]*userData

	idPool utils.IdPool

	msgBuffers   []*versionChatMessage
	nextWritePos atomic.Uint64

	taskqueue utils.TaskQueue

	shutDownCtx context.Context
	closeFunc   context.CancelFunc
}

func (cr *ChatRoomImpl) Init(params *CreateRoomParams) {
	cr.usersData = make(map[int64]*userData)
	cr.msgBuffers = make([]*versionChatMessage, utils.GetPow2_32(10000))
	cr.nextWritePos.Store(0)
	cr.shutDownCtx, cr.closeFunc = context.WithCancel(context.Background())
	cr.taskqueue.Init(utils.WithMaxQueue(1024))
	cr.idPool.Init()

	cr.roomKey = params.RoomKey
	cr.immediatePush = params.ImmediatePush
	cr.interval = params.Interval

	go cr.tick()
}

func (cr *ChatRoomImpl) Close() {
	cr.closeFunc()
	cr.taskqueue.Close(utils.CloseWayImmediate)
}

func (cr *ChatRoomImpl) Context() context.Context {
	return cr.shutDownCtx
}

func (cr *ChatRoomImpl) send2Session(ud *userData) {
	wp := cr.nextWritePos.Load()
	var msgs []*ChatMessage
	for {
		nr := ud.nextReadPos.Load()
		if nr >= wp {
			break
		}

		m := cr.msgBuffers[nr&uint64(len(cr.msgBuffers)-1)]
		if m == nil || m.version < nr {
			break
		}
		if m.version > nr {
			continue
		}
		msgs = append(msgs, m.msg)

		ud.nextReadPos.Add(1)
	}
	if len(msgs) > 0 {
		ud.sendFunc(msgs)
		ud.update(time.Now())
	}
}

func (cr *ChatRoomImpl) Join(params *JoinParams) int64 {
	wp := cr.nextWritePos.Load()
	nr := uint64(0)
	if !params.NeedRecent {
		nr = wp
	} else if wp >= 100 {
		nr = wp - 100
	}
	ud := &userData{
		id:               cr.idPool.NewId(),
		sessionId:        params.SessionId,
		sendFunc:         params.SendFunc,
		nextReadPos:      atomic.Uint64{},
		maxCacheNum:      params.MaxCacheNum,
		maxCacheDuration: params.MaxCacheDuration,
		lastPush:         time.Now(),
	}
	ud.nextReadPos.Store(nr)

	log.Debugf("[chatroom] sid:%d id:%d", params.SessionId, ud.id)

	cr.mtx.Lock()
	cr.usersData[ud.id] = ud
	cr.mtx.Unlock()

	if params.NeedRecent {
		cr.taskqueue.Put(func() {
			cr.send2Session(ud)
		})
	}

	return ud.id
}

func (cr *ChatRoomImpl) Leave(id int64) {
	cr.mtx.Lock()
	defer cr.mtx.Unlock()
	ud, ok := cr.usersData[id]
	if ok {
		log.Debugf("[chatroom] sid:%d id:%d", ud.sessionId, ud.id)
		delete(cr.usersData, ud.id)
	}
}

func (cr *ChatRoomImpl) Broadcast(m *ChatMessage) {
	pos := utils.FetchAndAdd(&cr.nextWritePos, uint64(1))
	cr.msgBuffers[pos&uint64(len(cr.msgBuffers)-1)] = &versionChatMessage{
		version: pos,
		msg:     m,
	}

	if cr.immediatePush {
		cr.taskqueue.Put(func() {
			cr.mtx.Lock()
			udtmp := make([]*userData, 0, len(cr.usersData))
			for _, ud := range cr.usersData {
				udtmp = append(udtmp, ud)
			}
			cr.mtx.Unlock()
			for _, ud := range udtmp {
				cr.send2Session(ud)
			}
		})
	}
}

func (cr *ChatRoomImpl) tick() {
	ticker := time.NewTicker(cr.interval)

loop:
	for {
		select {
		case <-ticker.C:
			cr.mtx.Lock()
			udtmp := make([]*userData, 0, len(cr.usersData))
			now := time.Now()
			nw := cr.nextWritePos.Load()
			for _, ud := range cr.usersData {
				diff := nw - ud.nextReadPos.Load()
				if diff != 0 && (ud.elapse(now) > ud.maxCacheDuration || diff > uint64(ud.maxCacheNum)) {
					udtmp = append(udtmp, ud)
				}
			}
			cr.mtx.Unlock()
			cr.taskqueue.Put(func() {
				for _, ud := range udtmp {
					cr.send2Session(ud)
				}
			})
		case <-cr.shutDownCtx.Done():
			break loop
		}
	}
}
