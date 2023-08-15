package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeaderTokenAuth(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        _, err := w.Write([]byte("blabla blabla"))
        require.NoError(t, err)
    })
    headerName := "Api-Token"
    ts := httptest.NewServer(Authentication("Api-Token", "1234567890")(handler))
    defer ts.Close()
    {
        req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
        require.NoError(t, err)
        req.Header.Set(headerName, "1234567890")
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
        assert.Equal(t, "Unauthorized\n", string(b))
     }
}
