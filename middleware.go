package rest

import (
	"net/http"
	"strings"
	"encoding/json"
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
		if r.Method == "POST" && strings.HasSuffix(strings.ToLower(r.URL.Path), "/ping") {
            w.Header().Set("Content-Type", "application/json")
            w.WriteHeader(http.StatusOK)
            resp := make(map[string]string)
            resp["message"] = "PONG"
            jsonResp, _ := json.Marshal(resp)
            _, _ = w.Write(jsonResp)
            return
        }
		next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}