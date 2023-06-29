package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
	"pnas/prpc"
	"pnas/setting"
	"strings"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var wg sync.WaitGroup

const (
	itemId = 5
)

func getCookie(header metadata.MD, key string) string {
	cookies := header["set-cookie"]
	for _, cookie := range cookies {
		i := strings.Index(cookie, key)
		if i >= 0 {
			j := strings.Index(cookie[i+len(key):], ";")
			return cookie[i : i+len(key)+j]
		}
	}
	return ""
}

func newSession() {
	defer wg.Add(-1)

	creds, err := credentials.NewClientTLSFromFile(setting.GS.Server.CrtFile, "")
	if err != nil {
		fmt.Printf("load cred %v\n", err)
		return
	}
	conn, _ := grpc.Dial(fmt.Sprintf("%s:%d", setting.GS.Server.Domain, setting.GS.Server.Port), grpc.WithTransportCredentials(creds))
	client := prpc.NewUserServiceClient(conn)
	h := md5.New()
	h.Write([]byte("123"))
	pwd := h.Sum(nil)
	var header metadata.MD
	_, err = client.Login(context.Background(), &prpc.LoginInfo{
		Email:      "12@12",
		Passwd:     hex.EncodeToString(pwd),
		RememberMe: true,
	}, grpc.Header(&header))
	if err != nil {
		fmt.Printf("login %v\n", err)
		return
	}
	cookie := fmt.Sprintf("%s;%s", getCookie(header, "token"), getCookie(header, "sessionid"))
	md := metadata.Pairs("cookie", cookie)
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	s, err := client.JoinChatRoom(ctx, &prpc.JoinChatRoomReq{
		ItemId: itemId,
	})
	if err != nil {
		fmt.Printf("join %v\n", err)
		return
	}
	go func() {
		for {
			r, _ := s.Recv()
			t, _ := time.Parse(time.RFC3339Nano, r.GetChatMsgs()[0].GetMsg())
			d := time.Since(t)
			if d > time.Millisecond*100 && d%4 == 0 {
				fmt.Printf("overload, millisec: %d\n", d/time.Millisecond)
			}
		}
	}()
	interval := time.Millisecond * 1000
	<-time.After(time.Duration(rand.Int63() % (int64(interval))))
	for {
		_, err := client.SendMsg2ChatRoom(ctx, &prpc.SendMsg2ChatRoomReq{
			ItemId: itemId,
			ChatMsg: &prpc.ChatMessage{
				Msg: time.Now().Format(time.RFC3339Nano),
			},
		})
		if err != nil {
			fmt.Printf("send %v\n", err)
		}
		<-time.After(interval)
	}
}

func main() {
	setting.Init(".")

	sc := 3000

	for i := 0; i < sc; i++ {
		wg.Add(1)
		go newSession()
	}

	wg.Wait()
}