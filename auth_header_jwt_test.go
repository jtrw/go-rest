package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"fmt"
)

func TestHeaderJwtTokenAuth(t *testing.T) {
    jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.WKQfGgHiRhXdkdz6Qy90gMQhYf3uK-GMeyAQBEs1EbQ"
    jwtFail := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.e30.1F5StBaWKNe53iB2919Agg3nMcCdwINDWlT0sNBaMbE"

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, err := w.Write([]byte("blabla blabla"))
        require.NoError(t, err)
    })
    headerName := "Api-Token"
    ts := httptest.NewServer(AuthenticationJwt(
    headerName,
    "1234567890",
    func(claims map[string]interface{}) error {return nil})(handler))
    defer ts.Close()
    {
        req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
        require.NoError(t, err)
        req.Header.Set(headerName, jwtToken)
        resp, err := http.DefaultClient.Do(req)
        require.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
        defer resp.Body.Close()
    }
    {
        req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
        require.NoError(t, err)
        req.Header.Set(headerName, "invalid")
        resp, err := http.DefaultClient.Do(req)
        require.NoError(t, err)
        assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
        defer resp.Body.Close()
        b, err := io.ReadAll(resp.Body)
        assert.NoError(t, err)
        assert.Equal(t, "Invalid token\n", string(b))
    }
    {
        req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
        require.NoError(t, err)
        req.Header.Set(headerName, jwtFail)
        resp, err := http.DefaultClient.Do(req)
        require.NoError(t, err)
        assert.Equal(t, http.StatusForbidden, resp.StatusCode)
        defer resp.Body.Close()
        b, err := io.ReadAll(resp.Body)
        assert.NoError(t, err)
        assert.Equal(t, "Forbidden\n", string(b))
    }
    {
        req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
        require.NoError(t, err)
        resp, err := http.DefaultClient.Do(req)
        require.NoError(t, err)
        assert.Equal(t, http.StatusForbidden, resp.StatusCode)
        defer resp.Body.Close()
        b, err := io.ReadAll(resp.Body)
        assert.NoError(t, err)
        assert.Equal(t, "Can not find token in header\n", string(b))
    }
}

func TestHeaderJwtTokenAuthCheckClaim(t *testing.T) {
    jwtToken := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTIzIn0.tsuWzcw0zCYzHoq0Kflun7cxVJWKMdQwWczNhU2h2IQ"
    jwtFail := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJub3RfdXNlciI6IjEyMyJ9.wGj6rh93KK83eaehCoxwmMCyEvyEQXadeJykayMkEd8"

    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, err := w.Write([]byte("blabla blabla"))
        require.NoError(t, err)
    })
    headerName := "Api-Token"
    ts := httptest.NewServer(AuthenticationJwt(
    headerName,
    "1234567890",
    func(claims map[string]interface{}) error {
        if claims["user_id"] == nil {
            return fmt.Errorf("user_id not found")
        }
        return nil
    })(handler))
    defer ts.Close()
    {
        req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
        require.NoError(t, err)
        req.Header.Set(headerName, jwtToken)
        resp, err := http.DefaultClient.Do(req)
        require.NoError(t, err)
        assert.Equal(t, 200, resp.StatusCode)
        defer resp.Body.Close()
    }
    {
        req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
        require.NoError(t, err)
        req.Header.Set(headerName, jwtFail)
        resp, err := http.DefaultClient.Do(req)
        require.NoError(t, err)
        assert.Equal(t, http.StatusUnauthorized, resp.StatusCode)
        defer resp.Body.Close()
        b, err := io.ReadAll(resp.Body)
        assert.NoError(t, err)
        assert.Equal(t, "user_id not found\n", string(b))
    }
}
