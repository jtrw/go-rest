package rest

import (
    "log"
	"net/http"
	"strings"
	"encoding/json"
	"runtime/debug"
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

func PanicRecovery(next http.Handler) http.Handler {
  return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    defer func() {
      if err := recover(); err != nil {
        http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
        log.Println("An error occurred:", err)
        log.Println(string(debug.Stack()))
      }
    }()
    next.ServeHTTP(w, req)
  })
}