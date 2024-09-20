package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIpAddress(t *testing.T) {
	t.Run("returns localhost", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/ip", nil)
		request.RemoteAddr = "198.51.100.10"
		response := httptest.NewRecorder()

		IpServer(response, request)

		var got map[string]string
		json.Unmarshal(response.Body.Bytes(), &got)
		want := map[string]string{
			"origin": "198.51.100.10",
		}

		assert.Equal(t, want, got)
	})
	t.Run("returns localhost on external interface", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		url := ts.URL + "/ip"
		response, err := client.Get(url)
		assert.Nil(t, err)

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		assert.Nil(t, err)
		var got map[string]string
		json.Unmarshal(body, &got)

		want := map[string]string{
			"origin": "127.0.0.1",
		}
		assert.Equal(t, want, got)
	})
	t.Run("POST not implemented", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		url := ts.URL + "/ip"

		response, err := client.Post(url, "", nil)

		assert.Nil(t, err)
		assert.Equal(t, 405, response.StatusCode)
	})
}
