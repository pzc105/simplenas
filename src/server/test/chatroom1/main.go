package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"pnas/prpc"
	"pnas/setting"
	"pnas/utils"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

var wg sync.WaitGroup
var sc atomic.Int64
var rc atomic.Int64

const (
	sessCount = 1500
	itemId    = 5
	email     = "12@12"
	pwd       = "123"
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

func newSession(first bool) {
	defer wg.Add(-1)

	creds, err := credentials.NewClientTLSFromFile(setting.GS().Server.CrtFile, "")
	if err != nil {
		fmt.Printf("load cred %v\n", err)
		return
	}
	conn, _ := grpc.Dial(fmt.Sprintf("%s:%d", setting.GS().Server.Domain, setting.GS().Server.Port), grpc.WithTransportCredentials(creds))
	client := prpc.NewUserServiceClient(conn)
	h := md5.New()
	h.Write([]byte(pwd))
	pwd := h.Sum(nil)
	var header metadata.MD
	_, err = client.Login(context.Background(), &prpc.LoginInfo{
		Email:      email,
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
		Room: &prpc.Room{
			Type: prpc.Room_Category,
			Id:   itemId,
		},
	})
	if err != nil {
		fmt.Printf("join %v\n", err)
		return
	}
	go func() {
		fm := make(map[int64]bool)
		for {
			r, _ := s.Recv()
			if first {
				rc.Add(int64(len(r.GetChatMsgs())))
			}
			for i := range r.GetChatMsgs() {

				m := r.GetChatMsgs()[i]
				c, _ := strconv.ParseInt(m.GetMsg(), 10, 64)
				if _, b := fm[c]; b {
					fmt.Printf("error msg: %d\n", c)
				}
				fm[c] = true
			}
		}
	}()

	for {
		lsc := utils.FetchAndAdd(&sc, int64(1))
		_, err := client.SendMsg2ChatRoom(ctx, &prpc.SendMsg2ChatRoomReq{
			Room: &prpc.Room{
				Type: prpc.Room_Category,
				Id:   itemId,
			},
			ChatMsg: &prpc.ChatMessage{
				Msg: fmt.Sprintf("%d", lsc),
			},
		})
		if err != nil {
			fmt.Printf("send %v\n", err)
		}
	}
}

func main() {
	setting.Init("../../server.yml")

	go newSession(true)
	go newSession(false)
	go newSession(false)

	wg.Add(1)
	go func() {
		defer wg.Add(-1)
		t := time.Now()
		for {
			<-time.After(time.Second)
			fmt.Printf("rc: %d now: %f sec\n", rc.Load(), time.Since(t).Seconds())
		}
	}()
	wg.Wait()
}
