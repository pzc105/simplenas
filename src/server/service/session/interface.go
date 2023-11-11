package session

import (
	"context"
	"net/http"
	"pnas/ptype"
)

type ISessions interface {
	NewSession(params *NewSessionParams) *Session
	GetSession(*http.Request) (*Session, error)
	GetSession2(context.Context) (*Session, error)
	GetSession3(sessionId ptype.SessionID) (*Session, error)
}
