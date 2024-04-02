package distributedchat

import (
	"sync"
	"sync/atomic"
	"time"
)

type ChatMessage struct {
	Id       int64
	UserId   int64
	SentTime time.Time
	Msg      string
}

type room struct {
	mtx           sync.Mutex
	users         map[int64]bool
	msgBuffers    []*ChatMessage
	nextWritePos  atomic.Uint64
	oldestMsgId   int64
	oldestMsgTime time.Time
}
