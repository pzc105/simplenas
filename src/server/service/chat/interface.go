package chat

import (
	"context"
	"pnas/ptype"
	"time"
)

type SendFunc func([]*ChatMessage)

type IRoom interface {
	Context() context.Context
	Join(*JoinParams) int64
	Leave(id int64)
	Broadcast(*ChatMessage)
}

type IRooms interface {
	CreateRoom(*CreateRoomParams) error
	Join(*JoinParams) (int64, error)
	Leave(roomKey string, id int64)
	Broadcast(roomKey string, msg *ChatMessage)
	GetRoom(roomKey string) IRoom
}

type JoinParams struct {
	RoomKey          string
	SessionId        ptype.SessionID
	SendFunc         SendFunc
	MaxCacheNum      uint64
	MaxCacheDuration time.Duration
	NeedRecent       bool
}

type CreateRoomParams struct {
	RoomKey       string
	ImmediatePush bool
	Interval      time.Duration
}
