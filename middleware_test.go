package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"bytes"
	"encoding/json"
)

func TestMiddleware_Ping(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("blabla blabla"))
		require.NoError(t, err)
	})
	ts := httptest.NewServer(Ping(handler))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	require.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "pong", string(b))

	resp, err = http.Get(ts.URL + "/blabla")
	require.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "blabla blabla", string(b))
}

func TestMiddleware_PingPost(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("blabla blabla"))
		require.NoError(t, err)
	})
	ts := httptest.NewServer(Ping(handler))
	defer ts.Close()

    var jsonData = []byte("")

	resp, err := http.Post(ts.URL + "/ping", "application/json", bytes.NewBuffer(jsonData))
	require.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()

	var res map[string]interface{}

    json.NewDecoder(resp.Body).Decode(&res)

	assert.NoError(t, err)
	assert.Equal(t, "PONG", res["message"])
}