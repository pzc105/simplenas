package chat

import (
	"fmt"
	"math/rand"
	"pnas/user"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestChatRoom(t *testing.T) {
	cr := &ChatRoomImpl{}
	cr.Init()

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

		cr.Join(sid, func(cm *ChatMessage) {
			if i1.Load() >= int32(len(msgs)) || msgs[i1.Load()] != cm {
				t.Errorf("i: %d, msg: %v", i1.Load(), cm)
			}
			i1.Add(1)
		})
		<-cr.Context().Done()
	}()
	var i2 atomic.Int32
	go func() {
		sid := int64(2)

		i2.Store(0)
		cr.Join(sid, func(cm *ChatMessage) {
			if i2.Load() >= int32(len(msgs)) || msgs[i2.Load()] != cm {
				t.Errorf("i: %d, msg: %v", i2.Load(), cm)
			}
			i2.Add(1)
		})
		<-cr.Context().Done()
	}()
	for _, cm := range msgs {
		cr.Broadcast(cm)
	}
	time.Sleep(time.Second * 2)
	fmt.Printf("i1: %d, i2: %d\n", i1.Load(), i2.Load())
	cr.Close()
}

func TestBenchmark(t *testing.T) {
	cr := &ChatRoomImpl{}
	cr.Init()

	maxDuration := time.Duration(0)

	for i := 0; i < 10000; i++ {
		cr.Join(int64(i), func(cm *ChatMessage) {
			d := time.Since(cm.SentTime)
			if d > maxDuration {
				maxDuration = d
			}
		})
	}
	var wg sync.WaitGroup
	for i := 0; i < 6; i++ {
		wg.Add(1)
		go func() {
			for j := 0; j < 100000; j++ {
				cr.Broadcast(&ChatMessage{
					SentTime: time.Now(),
				})
			}
			wg.Add(-1)
		}()
	}
	wg.Wait()
	fmt.Printf("maxDuration: %d\n", maxDuration)
}
