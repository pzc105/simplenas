package chat

import (
	"context"
	"pnas/category"
)

type SendFunc func([]*ChatMessage)

type IRoom interface {
	Context() context.Context
	Join(sessionId int64, sendFunc SendFunc)
	Leave(sessionId int64)
	Broadcast(*ChatMessage)
}

type IRooms interface {
	Join(itemId category.ID, sessionId int64, sendFunc SendFunc)
	Leave(itemId category.ID, sessionId int64)
	Broadcast(itemId category.ID, msg *ChatMessage)
	GetRoom(itemId category.ID) (IRoom, error)
}
