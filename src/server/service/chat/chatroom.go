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
	sessionId   int64
	sendFunc    SendFunc
	nextReadPos uint64
	joinCount   uint32
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
	mtx       sync.Mutex
	usersData map[int64]*userData

	msgBuffers   []*versionChatMessage
	nextWritePos atomic.Uint64

	taskqueue utils.TaskQueue

	shutDownCtx context.Context
	closeFunc   context.CancelFunc
}

func (cr *ChatRoomImpl) Init() {
	cr.usersData = make(map[int64]*userData)
	cr.msgBuffers = make([]*versionChatMessage, utils.GetPow2_32(10000))
	cr.nextWritePos.Store(0)
	cr.shutDownCtx, cr.closeFunc = context.WithCancel(context.Background())
	cr.taskqueue.Init(utils.WithMaxQueue(1024))
}

func (cr *ChatRoomImpl) Close() {
	cr.closeFunc()
	cr.taskqueue.Close()
}

func (cr *ChatRoomImpl) Context() context.Context {
	return cr.shutDownCtx
}

func (cr *ChatRoomImpl) send2Session(ud *userData) {
	wp := cr.nextWritePos.Load()
	var msgs []*ChatMessage
	for ; ud.nextReadPos < wp; ud.nextReadPos++ {
		m := cr.msgBuffers[ud.nextReadPos&uint64(len(cr.msgBuffers)-1)]
		if m == nil || m.version < ud.nextReadPos {
			break
		}
		if m.version > ud.nextReadPos {
			continue
		}
		msgs = append(msgs, m.msg)
	}
	if len(msgs) > 0 {
		ud.sendFunc(msgs)
	}
}

func (cr *ChatRoomImpl) Join(sessionId int64, sendFunc SendFunc) {
	wp := cr.nextWritePos.Load()
	nr := uint64(0)
	if wp >= 100 {
		nr = wp - 100
	}
	ud := &userData{
		sessionId:   sessionId,
		sendFunc:    sendFunc,
		nextReadPos: nr,
		joinCount:   1,
	}

	cr.mtx.Lock()
	oud, ok := cr.usersData[sessionId]
	if !ok {
		cr.usersData[sessionId] = ud
		log.Debugf("[chatroom] sid:%d join count:%d", sessionId, ud.joinCount)
	} else {
		oud.sendFunc = sendFunc
		oud.joinCount += 1
		log.Debugf("[chatroom] sid:%d join count:%d", sessionId, oud.joinCount)
	}
	cr.mtx.Unlock()

	cr.taskqueue.Put(func() {
		cr.send2Session(ud)
	})
}

func (cr *ChatRoomImpl) Leave(sessionId int64) {
	cr.mtx.Lock()
	defer cr.mtx.Unlock()
	ud, ok := cr.usersData[sessionId]
	if ok {
		ud.joinCount -= 1
		log.Debugf("[chatroom] sid:%d leave count:%d", sessionId, ud.joinCount)
		if ud.joinCount == 0 {
			delete(cr.usersData, sessionId)
		}
	}
}

func (cr *ChatRoomImpl) Broadcast(m *ChatMessage) {
	pos := utils.FetchAndAdd(&cr.nextWritePos, uint64(1))
	cr.msgBuffers[pos&uint64(len(cr.msgBuffers)-1)] = &versionChatMessage{
		version: pos,
		msg:     m,
	}
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
