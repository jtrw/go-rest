package rest

import (
	"net/http"
	"strings"
)

func AuthBearer(token string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			authorization := r.Header.Get("Authorization")
			headerToken := strings.TrimSpace(strings.Replace(authorization, "Bearer", "", 1))

			if headerToken == "" {
				http.Error(w, "Token not Found", http.StatusUnauthorized)
				return
			}

			if headerToken != token {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
