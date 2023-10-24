package main

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"math/rand"
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
var activeCount atomic.Int64
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

func newSession() {
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
		ItemId: itemId,
	})
	if err != nil {
		fmt.Printf("join %v\n", err)
		return
	}
	go func() {
		fm := make(map[int64]bool)
		for {
			r, _ := s.Recv()
			rc.Add(int64(len(r.GetChatMsgs())))
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
	activeCount.Add(1)
	defer activeCount.Add(-1)

	for activeCount.Load() != sessCount {
		<-time.After(time.Duration(time.Millisecond * 10))
	}
	<-time.After(time.Duration(time.Second * 1))
	interval := time.Millisecond * 1000
	<-time.After(time.Duration(rand.Int63() % (int64(interval))))

	rLoopCount := 20
	for {
		lt := time.Now()
		lsc := utils.FetchAndAdd(&sc, int64(1))
		_, err := client.SendMsg2ChatRoom(ctx, &prpc.SendMsg2ChatRoomReq{
			ItemId: itemId,
			ChatMsg: &prpc.ChatMessage{
				Msg: fmt.Sprintf("%d", lsc),
			},
		})
		if err != nil {
			fmt.Printf("send %v\n", err)
		}
		sd := time.Since(lt) / time.Millisecond
		if sd < interval {
			<-time.After(interval - sd)
		}
		rLoopCount -= 1
		if rLoopCount <= 0 {
			break
		}
	}
}

func main() {
	setting.Init(".")

	for i := 0; i < sessCount; i++ {
		wg.Add(1)
		go newSession()
	}

	wg.Add(1)
	go func() {
		defer wg.Add(-1)
		t := time.Now()
		for {
			<-time.After(time.Second)
			fmt.Printf("sessCount: %d, sc: %d, rc: %d t: %f sec\n", activeCount.Load(), sc.Load(), rc.Load(), time.Since(t).Seconds())
			if activeCount.Load() == 0 {
				return
			}
		}
	}()

	wg.Wait()
}
