package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaders(t *testing.T) {
	path := "/headers"
	t.Run("GreyboxTest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, path, nil)
		request.RemoteAddr = "198.51.100.10"
		request.Header["foo"] = []string{"bar"}
		response := httptest.NewRecorder()

		HeadersServer(response, request)

		var got map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &got)
		want := map[string]interface{}{
			"headers": map[string]interface{}{
				"foo": "bar",
			},
		}
		assert.Equal(t, want, got)
	})
	t.Run("BlackboxTest", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		url := ts.URL + path
		response, err := client.Get(url)
		assert.Nil(t, err)

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		var got map[string]interface{}
		json.Unmarshal(body, &got)
		assert.Nil(t, err)

		want := map[string]interface{}{
			"headers": map[string]interface{}{
				"Accept-Encoding": "gzip",
				"User-Agent":      "Go-http-client/1.1",
			},
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
