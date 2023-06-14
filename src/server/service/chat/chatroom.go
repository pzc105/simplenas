package chat

import (
	"context"
	"pnas/utils"
	"sync"
	"sync/atomic"
)

const (
	HistoryMsgMaxCount = uint64(10)
)

type userData struct {
	sendFunc    SendFunc
	nextReadPos uint64
}

type versionChatMessage struct {
	version uint64
	msg     *ChatMessage
}

type ChatRoomImpl struct {
	ChatRoom
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
	cr.taskqueue.Init(utils.WithMaxQueue(126))
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
	for ; ud.nextReadPos < wp; ud.nextReadPos++ {
		m := cr.msgBuffers[ud.nextReadPos&uint64(len(cr.msgBuffers)-1)]
		if m == nil || m.version != ud.nextReadPos {
			break
		}
		ud.sendFunc(m.msg)
	}
}

func (cr *ChatRoomImpl) Join(sessionId int64, sendFunc SendFunc) {
	ud := &userData{
		sendFunc:    sendFunc,
		nextReadPos: 0,
	}

	cr.mtx.Lock()
	cr.usersData[sessionId] = ud
	cr.mtx.Unlock()

	cr.taskqueue.Put(func() {
		cr.send2Session(ud)
	})
}

func (cr *ChatRoomImpl) Leave(sessionId int64) {
	cr.mtx.Lock()
	defer cr.mtx.Unlock()
	_, ok := cr.usersData[sessionId]
	if ok {
		delete(cr.usersData, sessionId)
	}
}

func (cr *ChatRoomImpl) Broadcast(m *ChatMessage) {
	pos := utils.FetchAndAdd(&cr.nextWritePos, uint64(1))
	cr.msgBuffers[pos&uint64(len(cr.msgBuffers)-1)] = &versionChatMessage{
		version: pos,
		msg:     m,
	}

	cr.mtx.Lock()
	defer cr.mtx.Unlock()
	for _, ud := range cr.usersData {
		ud2 := ud
		cr.taskqueue.Put(func() {
			cr.send2Session(ud2)
		})
	}
}
