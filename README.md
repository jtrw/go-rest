# go-rest

[![Build Status](https://github.com/jtrw/go-rest/workflows/Build/badge.svg)](https://github.com/jtrw/go-rest/actions)
[![Coverage Status](https://coveralls.io/repos/github/jtrw/go-rest/badge.svg?branch=master)](https://coveralls.io/github/jtrw/go-rest?branch=master)

Middleware for REST API

1. [Ping](#ping-middleware)
2. [PanicRecovery](#panicrecovery-middleware)
3. [Authentication via header](#authentication-via-header)
4. [Authentication via Bearer token](#authentication-via-bearer-token)
5. [Authentication via JWT token](#authentication-via-jwt-token)
6. [BasicAuth](#basicauth-middleware)
7. [RealIP](#realip-middleware)
8. [SizeLimit](#sizelimit-middleware)

## Usage

### Ping middleware

Ping middleware is a simple middleware that returns a 200 OK response.

### Authentication via header

Authentication is a middleware that checks for a header with a given name and value.

### Authentication via JWT token

AuthenticationJWT is a middleware that checks for a JWT token in the Authorization header.

Example:
```
func main() {
	http.Handle("/", middleware.AuthenticationJWT("Jwt-Token", "secret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})))

	http.ListenAndServe(":8080", nil)
}
``


### Authentication via bearer token

AuthenticationBearer is a middleware that checks for a bearer token in the Authorization header.

Example:
```
func main() {
	http.Handle("/", middleware.AuthenticationBearer("Bearer", "secret", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})))

	http.ListenAndServe(":8080", nil)
}
```


### BasicAuth middleware

BasicAuth middleware checks for a username and password in the request's Authorization header.

For form use BasicAuthWithPrompt middleware, which prompts for username and password if they are not provided.
Example:
```

func main() {
	http.Handle("/", middleware.BasicAuth("username", "password", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Hello, World!")
	})))

	http.ListenAndServe(":8080", nil)
}
```


### SizeLimit middleware

SizeLimit middleware checks if body size is above the limit and returns StatusRequestEntityTooLarge (413)

### RealIP middleware

RealIP is a middleware that sets a http.Request's RemoteAddr to the results of parsing either the X-Forwarded-For or X-Real-IP headers.
