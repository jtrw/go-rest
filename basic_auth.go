package rest

import (
	"context"
	"crypto/subtle"
	"net/http"
)

const baContextKey = "authorizedWithBasicAuth"

type contextKey string

func BasicAuth(checker func(user, passwd string) bool) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			u, p, ok := r.BasicAuth()
			if !ok {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			if !checker(u, p) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey(baContextKey), true)))
		}
		return http.HandlerFunc(fn)
	}
}

func BasicAuthWithUserPasswd(user, passwd string) func(http.Handler) http.Handler {
	checkFn := func(reqUser, reqPasswd string) bool {
		matchUser := subtle.ConstantTimeCompare([]byte(user), []byte(reqUser))
		matchPass := subtle.ConstantTimeCompare([]byte(passwd), []byte(reqPasswd))
		return matchUser == 1 && matchPass == 1
	}
	return BasicAuth(checkFn)
}

func IsAuthorized(ctx context.Context) bool {
	v := ctx.Value(contextKey(baContextKey))
	return v != nil && v.(bool)
}

func BasicAuthWithPrompt(user, passwd string) func(http.Handler) http.Handler {
	checkFn := func(reqUser, reqPasswd string) bool {
		matchUser := subtle.ConstantTimeCompare([]byte(user), []byte(reqUser))
		matchPass := subtle.ConstantTimeCompare([]byte(passwd), []byte(reqPasswd))
		return matchUser == 1 && matchPass == 1
	}

	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {

			// extract basic auth from request
			u, p, ok := r.BasicAuth()
			if ok && checkFn(u, p) {
				h.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), contextKey(baContextKey), true)))
				return
			}
			// not authorized, prompt for basic auth
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
		}
		return http.HandlerFunc(fn)
	}
}
