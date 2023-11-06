package chat

import (
	"errors"
	"sync"
)

type ChatRoomID int64

type Rooms struct {
	IRooms
	mtx sync.Mutex
	rs  map[string]IRoom
}

func (cs *Rooms) Init() {
	cs.rs = make(map[string]IRoom)
}

func (cs *Rooms) CreateRoom(params *CreateRoomParams) error {
	cs.mtx.Lock()
	if _, ok := cs.rs[params.RoomKey]; ok {
		return errors.New("existed room")
	}
	var nr ChatRoomImpl
	nr.Init(params)
	cs.rs[params.RoomKey] = &nr
	cs.mtx.Unlock()
	return nil
}

func (cs *Rooms) fetchRoom(roomKey string) IRoom {
	cs.mtx.Lock()
	r, ok := cs.rs[roomKey]
	cs.mtx.Unlock()
	if !ok {
		return nil
	}
	return r
}

func (cs *Rooms) Join(params *JoinParams) (int64, error) {
	r := cs.fetchRoom(params.RoomKey)
	if r != nil {
		return r.Join(params), nil
	}
	return -1, errors.New("not found room")
}

func (cs *Rooms) Leave(roomKey string, id int64) {
	r := cs.fetchRoom(roomKey)
	if r != nil {
		r.Leave(id)
	}
}

func (cs *Rooms) Broadcast(roomKey string, msg *ChatMessage) {
	r := cs.fetchRoom(roomKey)
	if r != nil {
		r.Broadcast(msg)
	}
}

func (cs *Rooms) GetRoom(roomKey string) IRoom {
	return cs.fetchRoom(roomKey)
}
