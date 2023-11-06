package chat

import (
	"fmt"
	"math/rand"
	"pnas/user"
	"sync/atomic"
	"testing"
	"time"
)

func TestChatRoom(t *testing.T) {
	cp := CreateRoomParams{}
	cr := &ChatRoomImpl{}
	cr.Init(&cp)

	msgs := []*ChatMessage{}
	for i := 0; i < 300000; i++ {
		msgs = append(msgs, &ChatMessage{
			UserId:   user.ID(rand.Int63()),
			SentTime: time.Now(),
			Msg:      fmt.Sprintf("%d", i),
		})
	}
	var i1 atomic.Int32
	go func() {
		sid := int64(1)

		params := &JoinParams{
			SessionId: sid,
			SendFunc: func(cms []*ChatMessage) {
				cm := cms[0]
				if i1.Load() >= int32(len(msgs)) || msgs[i1.Load()] != cm {
					t.Errorf("i: %d, msg: %v", i1.Load(), cm)
				}
				i1.Add(1)
			},
		}
		cr.Join(params)
		<-cr.Context().Done()
	}()
	var i2 atomic.Int32
	go func() {
		sid := int64(2)

		i2.Store(0)
		params := &JoinParams{
			SessionId: sid,
			SendFunc: func(cms []*ChatMessage) {
				cm := cms[0]
				if i2.Load() >= int32(len(msgs)) || msgs[i2.Load()] != cm {
					t.Errorf("i: %d, msg: %v", i2.Load(), cm)
				}
				i2.Add(1)
			},
		}
		cr.Join(params)
		<-cr.Context().Done()
	}()
	for _, cm := range msgs {
		cr.Broadcast(cm)
	}
	time.Sleep(time.Second * 2)
	fmt.Printf("i1: %d, i2: %d\n", i1.Load(), i2.Load())
	cr.Close()
}
