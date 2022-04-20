package rest

import (
	"net/http"
	"strings"
)

// Ping middleware response with pong to /ping. Stops chain if ping request detected
func Ping(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "GET" && strings.HasSuffix(strings.ToLower(r.URL.Path), "/ping") {
			w.Header().Set("Content-Type", "text/plain")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("pong"))
			return
		}
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}

// RealIP is a middleware that sets a http.Request's RemoteAddr to the results
// of parsing either the X-Forwarded-For or X-Real-IP headers.
//
// This middleware should only be used if user can trust the headers sent with request.
// If reverse proxies are configured to pass along arbitrary header values from the client,
// or if this middleware used without a reverse proxy, malicious clients could set anything
// as X-Forwarded-For header and attack the server in various ways.
func RealIP(h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		if rip, err := realip.Get(r); err == nil {
			r.RemoteAddr = rip
		}
		h.ServeHTTP(w, r)
	}

	return http.HandlerFunc(fn)
}