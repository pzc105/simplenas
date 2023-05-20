package service

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"pnas/db"
	"pnas/log"
	"pnas/prpc"
	"pnas/user"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/metadata"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.abcdefghijk"
	letterBytesLen = len(letterBytes)
	letterIdxBits  = 6
	letterIdxMask  = 1<<letterIdxBits - 1
	token_len      = 24

	SessionIdFieldName = "sessionid"
	ToeknFieldName     = "token"

	SessionRedisKey = "session"
)

var id_pool IdPool

type session struct {
	Id        int64
	UserId    user.ID
	Token     string
	ExpiresAt time.Time

	needPush   atomic.Bool
	btStatusCh chan *prpc.StatusRespone
}

func InitIdPool() {
	id_pool.Init()

	rand.Seed(time.Now().UnixNano())
	result, err := db.GREDIS.Keys(context.Background(), SessionRedisKey+"*").Result()
	if err != nil {
		log.Errorf("failed to init id pool, err: %v", err)
		return
	}
	for _, k := range result {
		id, err := strconv.ParseInt(k[len(SessionRedisKey):], 10, 64)
		if err != nil {
			continue
		}
		id_pool.Allocated(id)
	}

	go clearIdTick()
}

func clearIdTick() {
	timer := time.NewTicker(60 * time.Second)
	for {
		<-timer.C
		for _, id := range id_pool.GetAllocatedIds() {
			_, err := db.GREDIS.Get(context.Background(), genSessionRedisKey(id)).Result()
			if err == redis.Nil {
				id_pool.ReleaseId(id)
			}
		}
	}
}

func NewToken() string {
	ret := make([]byte, token_len)
	for i := range ret {
		ret[i] = letterBytes[rand.Int63()&int64(letterIdxMask)]
	}
	return string(ret)
}

func GenCookieSessionById(id int64, rememberMe bool, userId user.ID) (metadata.MD, *session) {
	var session session
	session.Id = id
	session.UserId = userId
	session.Token = NewToken()
	session.ExpiresAt = time.Now().Add(time.Hour * 24 * 7)
	var expiresPair string
	if rememberMe {
		expiresPair = "Expires=" + session.ExpiresAt.Format(time.RFC1123)
	}

	cookieToken := fmt.Sprintf("%s=%s;SameSite=Strict;Path=/;HttpOnly;%s",
		ToeknFieldName,
		session.Token,
		expiresPair)
	cookieTokenId := fmt.Sprintf("%s=%d;SameSite=Strict;Path=/;HttpOnly;%s",
		SessionIdFieldName,
		id,
		expiresPair)
	return metadata.Pairs("Set-Cookie", cookieToken, "Set-Cookie", cookieTokenId), &session
}

func GenCookieSession(rememberMe bool, userId user.ID) (metadata.MD, *session) {
	return GenCookieSessionById(id_pool.NewId(), rememberMe, userId)
}

func GetTokenAndIdByCookie(cookie string) (token string, id int64) {
	header := http.Header{}
	header.Add("cookie", cookie)
	request := http.Request{Header: header}
	cookies := request.Cookies()

	for i := range cookies {
		cookie := cookies[i]
		if cookie.Name == SessionIdFieldName {
			var err error
			id, err = strconv.ParseInt(cookie.Value, 10, 64)
			if err != nil {
				return
			}
		} else if cookie.Name == ToeknFieldName {
			token = cookie.Value
		}
	}
	return token, id
}

func GetTokenAndId(ctx context.Context) (token string, id int64) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", -1
	}

	raw := md.Get("cookie")
	if len(raw) == 0 {
		return "", -1
	}
	return GetTokenAndIdByCookie(raw[0])
}

func genSessionRedisKey(id int64) string {
	return fmt.Sprintf("%s%d", SessionRedisKey, id)
}

func saveSession(session *session) error {
	jsonStr, err := json.Marshal(session)
	if err != nil {
		return err
	}
	err = db.GREDIS.Set(context.Background(), genSessionRedisKey(session.Id), string(jsonStr), time.Until(session.ExpiresAt)).Err()
	if err != nil {
		return err
	}
	return nil
}

func loadSession(id int64) (*session, error) {
	jsonStr, err := db.GREDIS.Get(context.Background(), genSessionRedisKey(id)).Result()
	if err != nil {
		return nil, err
	}
	var session session
	err = json.Unmarshal([]byte(jsonStr), &session)
	return &session, err
}
