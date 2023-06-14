package service

import (
	"errors"
	"pnas/category"
	"pnas/service/chat"
	"pnas/utils"
	"sync"
)

type ChatRoomID int64

type ChatRoomService struct {
	mtx     sync.Mutex
	rs      map[ChatRoomID]chat.ChatRoom
	cid2rid map[category.ID]ChatRoomID

	roomIdPool utils.IdPool
}

func (cs *ChatRoomService) Init() {
	cs.rs = make(map[ChatRoomID]chat.ChatRoom)
	cs.cid2rid = make(map[category.ID]ChatRoomID)
	cs.roomIdPool.Init()
}

func (cs *ChatRoomService) Join(itemId category.ID, sessionId int64, sendFunc chat.SendFunc) {
	cs.mtx.Lock()
	rid, ok := cs.cid2rid[itemId]
	var r chat.ChatRoom
	if !ok {
		rid = ChatRoomID(cs.roomIdPool.NewId())
		ri := &chat.ChatRoomImpl{}
		ri.Init()
		cs.rs[rid] = ri
		cs.cid2rid[itemId] = rid
		r = ri
	} else {
		r = cs.rs[rid]
	}
	cs.mtx.Unlock()
	r.Join(sessionId, sendFunc)
}

func (cs *ChatRoomService) Leave(itemId category.ID, sessionId int64) {
	cs.mtx.Lock()
	var r chat.ChatRoom
	rid, ok := cs.cid2rid[itemId]
	if ok {
		r = cs.rs[rid]
	} else {
		cs.mtx.Unlock()
		return
	}
	cs.mtx.Unlock()
	r.Leave(sessionId)
}

func (cs *ChatRoomService) Broadcast(itemId category.ID, msg *chat.ChatMessage) {
	cs.mtx.Lock()
	var r chat.ChatRoom
	rid, ok := cs.cid2rid[itemId]
	if ok {
		r = cs.rs[rid]
	} else {
		cs.mtx.Unlock()
		return
	}
	cs.mtx.Unlock()
	r.Broadcast(msg)
}

func (cs *ChatRoomService) GetRoom(itemId category.ID) (chat.ChatRoom, error) {
	cs.mtx.Lock()
	defer cs.mtx.Unlock()
	rid, ok := cs.cid2rid[itemId]
	if !ok {
		return nil, errors.New("not found room")
	}
	room, ok := cs.rs[rid]
	if !ok {
		return nil, errors.New("not found room")
	}
	return room, nil
}
