package session

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"pnas/db"
	"pnas/log"
	"pnas/ptype"
	"pnas/utils"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc/metadata"
)

const (
	letterBytes    = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ.0123456789?"
	letterBytesLen = len(letterBytes)
	letterIdxBits  = 6
	letterIdxMask  = 1<<letterIdxBits - 1
	token_len      = 24

	SessionIdFieldName = "sessionid"
	ToeknFieldName     = "token"

	SessionRedisKey = "session"
)

type Session struct {
	Id        ptype.SessionID
	UserId    ptype.UserID
	Token     string
	ExpiresAt time.Time
}

var (
	session_rand *rand.Rand
	randMtx      sync.Mutex

	oneChecker atomic.Bool
)

func init() {
	session_rand = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func NewToken() string {
	ret := make([]byte, token_len)
	for i := range ret {
		randMtx.Lock()
		rv := session_rand.Int63()
		randMtx.Unlock()
		ret[i] = letterBytes[rv&int64(letterIdxMask)]
	}
	return string(ret)
}

type Sessions struct {
	ISessions
	idPool utils.IdPool
}

func (ss *Sessions) Init() {
	if !oneChecker.CompareAndSwap(false, true) {
		panic("only one sessions")
	}

	ss.idPool.Init()
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
		ss.idPool.Allocated(id)
	}

	go ss.clearIdTick()
}

func (ss *Sessions) clearIdTick() {
	timer := time.NewTicker(60 * time.Second)
	for {
		<-timer.C
		for _, id := range ss.idPool.GetAllocatedIds() {
			id2 := ptype.SessionID(id)
			_, err := db.GREDIS.Get(context.Background(), genSessionRedisKey(id2)).Result()
			if err == redis.Nil {
				ss.idPool.ReleaseId(id)
			}
		}
	}
}

type NewSessionParams struct {
	OldId     ptype.SessionID
	ExpiresAt time.Time
	UserId    ptype.UserID
}

func (ss *Sessions) NewSession(params *NewSessionParams) *Session {
	var session Session
	if params.OldId < 0 {
		session.Id = ptype.SessionID(ss.idPool.NewId())
	} else {
		session.Id = params.OldId
	}
	session.UserId = params.UserId
	session.Token = NewToken()
	if !params.ExpiresAt.IsZero() {
		session.ExpiresAt = params.ExpiresAt
	} else {
		session.ExpiresAt = time.Now().Add(time.Hour * 1)
	}
	saveSession(&session)
	return &session
}

func (ss *Sessions) getSession(cookie string) (*Session, error) {
	header := http.Header{}
	header.Add("cookie", cookie)
	request := http.Request{Header: header}
	cookies := request.Cookies()

	var id ptype.SessionID
	var err error
	var token string
	for i := range cookies {
		cookie := cookies[i]
		if cookie.Name == SessionIdFieldName {
			idt, err := strconv.ParseInt(cookie.Value, 10, 64)
			if err != nil {
				return nil, err
			}
			id = ptype.SessionID(idt)
		} else if cookie.Name == ToeknFieldName {
			token = cookie.Value
		}
	}
	s, err := loadSession(id)
	if err != nil || s.Token != token {
		return nil, errors.New("not found session")
	}
	return s, nil
}

func (ss *Sessions) GetSession(r *http.Request) (*Session, error) {
	return ss.getSession(r.Header.Get("cookie"))
}

func (ss *Sessions) GetSession2(ctx context.Context) (*Session, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, errors.New("not found context")
	}

	raw := md.Get("cookie")
	if len(raw) == 0 {
		return nil, errors.New("not found cookie")
	}
	return ss.getSession(raw[0])
}

func (ss *Sessions) GetSession3(sessionId ptype.SessionID) (*Session, error) {
	return loadSession(sessionId)
}

func genSessionRedisKey(id ptype.SessionID) string {
	return fmt.Sprintf("%s%d", SessionRedisKey, id)
}

func saveSession(session *Session) error {
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

func loadSession(id ptype.SessionID) (*Session, error) {
	jsonStr, err := db.GREDIS.Get(context.Background(), genSessionRedisKey(id)).Result()
	if err != nil {
		return nil, err
	}
	var session Session
	err = json.Unmarshal([]byte(jsonStr), &session)
	return &session, err
}

func GenSessionTokenCookie(s *Session) string {
	expiresPair := "Expires=" + s.ExpiresAt.Format(time.RFC1123)
	tokenCookie := fmt.Sprintf("%s=%s;SameSite=Strict;Path=/;HttpOnly;%s",
		ToeknFieldName,
		s.Token,
		expiresPair)
	return tokenCookie
}

func GenSessionIdCookie(s *Session) string {
	expiresPair := "Expires=" + s.ExpiresAt.Format(time.RFC1123)
	tokenIdCookie := fmt.Sprintf("%s=%d;SameSite=Strict;Path=/;HttpOnly;%s",
		SessionIdFieldName,
		s.Id,
		expiresPair)
	return tokenIdCookie
}
