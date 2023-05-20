package phttp

import "net/http"

func CorsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		headers.Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
		headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		headers.Set("Access-Control-Allow-Headers", r.Header.Get("Access-Control-Request-Headers"))
		headers.Set("Access-Control-Expose-Headers", "Access-Control-Allow-Credentials, Set-Cookie, Date, Content-Type, Vary, Access-Control-Allow-Origin, grpc-status, grpc-message")
		headers.Set("Access-Control-Allow-Credentials", "true")

		if r.Method == "OPTIONS" && headers.Get("Access-Control-Request-Method") != "" {
			headers.Add("Vary", "Origin")
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
