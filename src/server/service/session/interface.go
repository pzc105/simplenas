package session

import (
	"context"
	"net/http"
)

type SessionsInterface interface {
	Init()
	NewSession(params *NewSessionParams) *Session
	GetSession(*http.Request) (*Session, error)
	GetSession2(context.Context) (*Session, error)
	GetSession3(sessionId int64) (*Session, error)
}
