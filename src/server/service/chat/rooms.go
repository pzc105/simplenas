package chat

import (
	"errors"
	"pnas/category"
	"pnas/utils"
	"sync"
)

type ChatRoomID int64

type Rooms struct {
	IRooms
	mtx     sync.Mutex
	rs      map[ChatRoomID]IRoom
	cid2rid map[category.ID]ChatRoomID

	roomIdPool utils.IdPool
}

func (cs *Rooms) Init() {
	cs.rs = make(map[ChatRoomID]IRoom)
	cs.cid2rid = make(map[category.ID]ChatRoomID)
	cs.roomIdPool.Init()
}

func (cs *Rooms) Join(itemId category.ID, sessionId int64, sendFunc SendFunc) {
	cs.mtx.Lock()
	rid, ok := cs.cid2rid[itemId]
	var r IRoom
	if !ok {
		rid = ChatRoomID(cs.roomIdPool.NewId())
		ri := &ChatRoomImpl{}
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

func (cs *Rooms) Leave(itemId category.ID, sessionId int64) {
	cs.mtx.Lock()
	var r IRoom
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

func (cs *Rooms) Broadcast(itemId category.ID, msg *ChatMessage) {
	cs.mtx.Lock()
	var r IRoom
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

func (cs *Rooms) GetRoom(itemId category.ID) (IRoom, error) {
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
