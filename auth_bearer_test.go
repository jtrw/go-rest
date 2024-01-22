package rest

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestBearer_AuthBearer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("blabla blabla"))
		require.NoError(t, err)
	})
	ts := httptest.NewServer(AuthBearer("123456")(handler))
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/ping")
	require.Nil(t, err)
	assert.Equal(t, 401, resp.StatusCode)
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "Token not Found\n", string(b))

	req, err := http.NewRequest("GET", ts.URL+"/ping", nil)
	require.NoError(t, err)
	req.Header.Set("Authorization", "Bearer 123456")
	resp, err = http.DefaultClient.Do(req)
	require.NoError(t, err)
	assert.Equal(t, 200, resp.StatusCode)
	defer resp.Body.Close()
	b, err = io.ReadAll(resp.Body)
	assert.NoError(t, err)
	assert.Equal(t, "blabla blabla", string(b))
}
