package distributedchat

import (
	"pnas/prpc"
	"pnas/service/distributedChat/room"
	"sync"
)

type chatService struct {
	prpc.UnimplementedChatServiceServer
	mtx   sync.Mutex
	rooms map[string]*room.Room
}
