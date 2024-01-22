package rest

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

	resp, err := http.Post(ts.URL+"/ping", "application/json", bytes.NewBuffer(jsonData))
	require.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()

	var res map[string]interface{}

	json.NewDecoder(resp.Body).Decode(&res)

	assert.NoError(t, err)
	assert.Equal(t, "PONG", res["message"])
}

func TestMiddleware_PanicRecovery(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("PANIC")
		//require.NoError(t, err)
	})
	ts := httptest.NewServer(PanicRecovery(handler))
	defer ts.Close()

	var jsonData = []byte("")

	resp, err := http.Post(ts.URL+"/error", "application/json", bytes.NewBuffer(jsonData))
	require.Nil(t, err)
	assert.Equal(t, 500, resp.StatusCode)
	defer resp.Body.Close()
}

func TestMiddleware_AppInfo(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("bla bla"))
		require.NoError(t, err)
	})
	ts := httptest.NewServer(AppInfo("app-name", "Nil", "12345")(handler))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/bla")
	require.Nil(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)

	assert.Equal(t, "bla bla", string(b))
	assert.Equal(t, "app-name", resp.Header.Get("App-Name"))
	assert.Equal(t, "12345", resp.Header.Get("App-Version"))
	assert.Equal(t, "Nil", resp.Header.Get("Author"))
}
