# go-rest

[![Build Status](https://github.com/jtrw/go-rest/workflows/Build/badge.svg)](https://github.com/jtrw/go-rest/actions)
[![Coverage Status](https://coveralls.io/repos/github/jtrw/go-rest/badge.svg?branch=master)](https://coveralls.io/github/jtrw/go-rest?branch=master)

Middleware for REST API

1. Ping
2. PanicRecovery
3. Authentication via header
4. Authentication via JWT token
5. RealIP
6. SizeLimit

### Authentication via header

Authentication is a middleware that checks for a header with a given name and value.

### Authentication via JWT token

AuthenticationJWT is a middleware that checks for a JWT token in the Authorization header.

### BasicAuth middleware

BasicAuth middleware checks for a username and password in the request's Authorization header.

### SizeLimit middleware

SizeLimit middleware checks if body size is above the limit and returns StatusRequestEntityTooLarge (413)

### RealIP middleware

RealIP is a middleware that sets a http.Request's RemoteAddr to the results of parsing either the X-Forwarded-For or X-Real-IP headers.
