package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUserAgent(t *testing.T) {
	path := "/user-agent"
	t.Run("GreyboxTest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, path, nil)
		request.Header["User-Agent"] = []string{"test"}
		response := httptest.NewRecorder()

		UserAgentServer(response, request)

		var got map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &got)
		want := map[string]interface{}{
			"user-agent": "test",
		}
		assert.Equal(t, want, got)
	})
	t.Run("BlackboxTest", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		url := ts.URL + path
		request, err := http.NewRequest("GET", url, nil)
		assert.Nil(t, err)
		request.Header.Set("User-Agent", "Test")

		response, err := client.Do(request)
		assert.Nil(t, err)

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		var got map[string]interface{}
		json.Unmarshal(body, &got)
		assert.Nil(t, err)

		want := map[string]interface{}{
			"user-agent": "Test",
		}

		assert.Equal(t, want, got)
	})
	t.Run("POST not implemented", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		url := ts.URL + path
		response, err := client.Post(url, "", nil)
		assert.Nil(t, err)

		assert.Equal(t, 405, response.StatusCode)
	})
}
