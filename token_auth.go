package rest

import (
	"fmt"
    "net/http"
    "github.com/golang-jwt/jwt"
)

const TOKEN_NAME = "Api-Token"

func Authentication(token string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            apiToken := r.Header.Get(TOKEN_NAME)
            if apiToken != token {
                w.Write([]byte("Unauthorized"));
                w.WriteHeader(http.StatusUnauthorized)
                return
            }
        }
        return http.HandlerFunc(fn)
    }
}

func AuthenticationJwt(secret string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        fn := func(w http.ResponseWriter, r *http.Request) {
            if r.Header[TOKEN_NAME] == nil {
                w.Write([]byte("Can not find token in header"));
                w.WriteHeader(http.StatusUnauthorized)
                return
            }

            token, _ := jwt.Parse(r.Header[TOKEN_NAME][0], func(token *jwt.Token) (interface{}, error) {
                if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                    return nil, fmt.Errorf("[ERROR] There was an error in parsing")
                }

                return []byte(secret), nil
            })

            if token == nil {
                w.Write([]byte("Invalid token"));
                w.WriteHeader(http.StatusUnauthorized)
                return
            }

            if !token.Valid {
                w.WriteHeader(http.StatusForbidden)
                return
            }

            _, ok := token.Claims.(jwt.MapClaims)

            if !ok {
                w.Write([]byte("couldn't parse claims"));
                w.WriteHeader(http.StatusUnauthorized)
                return
            }

//             if claims["user_id"] == nil {
//                 w.Write([]byte("user_id not found"));
//                 w.WriteHeader(http.StatusUnauthorized)
//                 return
//             }
            next.ServeHTTP(w, r)
        }
        return http.HandlerFunc(fn)
    }
}
