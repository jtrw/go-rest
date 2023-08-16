package rest

import (
	"fmt"
    "net/http"
    "github.com/golang-jwt/jwt"
)

func AuthenticationJwt(headerName, secret string,  userCondition func(claims map[string]interface{}) error) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            if r.Header[headerName] == nil {
                http.Error(w, "Can not find token in header", http.StatusForbidden)
                return
            }

            token, _ := jwt.Parse(r.Header[headerName][0], func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("[ERROR] There was an error in parsing")
                }

                return []byte(secret), nil
            })

            if token == nil {
                http.Error(w, "Invalid token", http.StatusUnauthorized)
                return
            }

            if !token.Valid {
                http.Error(w, "Forbidden", http.StatusForbidden)
                return
            }

            claims, ok := token.Claims.(jwt.MapClaims)

            if !ok {
                w.Write([]byte("couldn't parse claims"));
                w.WriteHeader(http.StatusUnauthorized)
                return
            }

            if err := userCondition(claims); err != nil {
                http.Error(w, err.Error(), http.StatusUnauthorized)
                return
            }

            next.ServeHTTP(w, r)
        }
        return http.HandlerFunc(fn)
    }
}
