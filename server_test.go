package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestIpAddress(t *testing.T) {
	t.Run("returns localhost", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/ip", nil)
		request.RemoteAddr = "198.51.100.10"
		response := httptest.NewRecorder()

		IpServer(response, request)

		var got  map[string]string
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
		var got  map[string]string
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

func TestHeaders(t *testing.T) {
	path := "/headers"
	t.Run("GreyboxTest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, path, nil)
		request.RemoteAddr = "198.51.100.10"
		request.Header["foo"] = []string{"bar"}
		response := httptest.NewRecorder()

		HeadersServer(response, request)

		var got  map[string]interface{}
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
		var got  map[string]interface{}
		json.Unmarshal(body, &got)
		assert.Nil(t, err)

		want := map[string]interface{}{
			"headers": map[string]interface{}{
				"Accept-Encoding": "gzip",
				"User-Agent": "Go-http-client/1.1",
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

func TestUserAgent(t *testing.T) {
	path := "/user-agent"
	t.Run("GreyboxTest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, path, nil)
		request.Header["User-Agent"] = []string{"test"}
		response := httptest.NewRecorder()

		UserAgentServer(response, request)

		var got  map[string]interface{}
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
		var got  map[string]interface{}
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

func TestGETPlayers(t *testing.T) {
	t.Run("returns Pepper's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Pepper", nil)
		request.SetPathValue("player", "Pepper")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "20"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("returns Floyd's score", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/Floyd", nil)
		request.SetPathValue("player", "Floyd")
		response := httptest.NewRecorder()

		PlayerServer(response, request)

		got := response.Body.String()
		want := "10"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})
}
