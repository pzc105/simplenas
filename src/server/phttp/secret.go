package phttp

import (
	"net/http"
	"time"
)

const (
	secretString = "peng"
)

func Secre(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	cipher := query.Get("cipher")
	if cipher == secretString {
		cookie := http.Cookie{
			Name:     "cipher",
			Value:    secretString,
			Expires:  time.Now().Add(24 * time.Hour * 7),
			HttpOnly: true,
		}
		http.SetCookie(w, &cookie)
		w.Write([]byte("ok"))
		return
	}
	cookie, err := r.Cookie("cipher")
	if err != nil {
		return
	}
	if cookie.Value == secretString {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("ok"))
	}
}

func SecretMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/secret" {
			next.ServeHTTP(w, r)
			return
		}
		cookie, err := r.Cookie("cipher")
		if err != nil {
			return
		}
		if cookie.Value == secretString {
			next.ServeHTTP(w, r)
		}
	})
}
