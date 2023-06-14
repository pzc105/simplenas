package chat

import (
	"context"
	"pnas/user"
	"time"
)

type SendFunc func(*ChatMessage)

type ChatRoom interface {
	Context() context.Context
	Join(sessionId int64, sendFunc SendFunc)
	Leave(sessionId int64)
	Broadcast(*ChatMessage)
}

type ChatMessage struct {
	UserId   user.ID
	SentTime time.Time
	Msg      string
}
