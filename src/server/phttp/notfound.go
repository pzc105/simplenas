package phttp

import (
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("404 - Page not found"))
	w.WriteHeader(http.StatusNotFound)
}
