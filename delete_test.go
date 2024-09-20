package main

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"strings"
)

func TestDelete(t *testing.T) {
	path := "/delete"
	t.Run("GreyboxTest", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, path + "?param1=value1;param2=value2", nil)
		request.Header["user-agent"] = []string{"test"}
		request.RemoteAddr = "198.51.100.11"
		response := httptest.NewRecorder()

		DeleteServer(response, request)

		var got map[string]interface{}
		json.Unmarshal(response.Body.Bytes(), &got)
		want := map[string]interface{}{
			"args":  map[string]interface{}{},
			"data":  "",
			"files": map[string]interface{}{},
			"form":  map[string]interface{}{},
			"headers": map[string]interface{}{
				"user-agent": "test",
			},
			"json":   nil,
			"origin": "198.51.100.11",
			"url":    "/delete?param1=value1;param2=value2",
		}
		assert.Equal(t, want, got)
	})

	t.Run("BlackboxTest", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		urlline := ts.URL + path + "?p1=v1&p2=v2&p2=v21"
		request, _ := http.NewRequest(http.MethodDelete, urlline, nil)
		response, err := client.Do(request)
		assert.Nil(t, err)

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		var got map[string]interface{}
		json.Unmarshal(body, &got)
		assert.Nil(t, err)

		want := map[string]interface{}{
			"args":  map[string]interface{}{
				"p1": "v1",
				"p2": []interface{}{"v2", "v21"},
			},
			"data":  "",
			"files": map[string]interface{}{},
			"form":  map[string]interface{}{},
			"headers": map[string]interface{}{
				"Accept-Encoding": "gzip",
				"User-Agent":      "Go-http-client/1.1",
			},
			"json":   nil,
			"origin": "127.0.0.1",
			"url":    "/delete?p1=v1&p2=v2&p2=v21",
		}

		assert.Equal(t, want, got)
	})
	t.Run("BlackboxTest form", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		urlline := ts.URL + path
		formData := url.Values{
	        "f1": {"v1"},
	        "f2": []string{"v2", "v3"},
	    }

		request, _ := http.NewRequest(http.MethodDelete, urlline, strings.NewReader(formData.Encode()))
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		response, err := client.Do(request)
		assert.Nil(t, err)

		body, err := io.ReadAll(response.Body)
		response.Body.Close()
		var got map[string]interface{}
		json.Unmarshal(body, &got)
		assert.Nil(t, err)

		want := map[string]interface{}{
			"args":  map[string]interface{}{},
			"data":  "",
			"files": map[string]interface{}{},
			"form":  map[string]interface{}{
				"f1": "v1",
				"f2": []string{"v2", "v3"},
			},
			"headers": map[string]interface{}{
				"Accept-Encoding": "gzip",
				"User-Agent":      "Go-http-client/1.1",
				"Content-Length":  "17",
				"Content-Type":    "application/x-www-form-urlencoded",
			},
			"json":   nil,
			"origin": "127.0.0.1",
			"url":    "/delete",
		}

		assert.Equal(t, want, got)
	})
	t.Run("GET not implemented", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		urlline := ts.URL + path
		response, err := client.Post(urlline, "", nil)
		assert.Nil(t, err)

		assert.Equal(t, 405, response.StatusCode)
	})
	t.Run("POST not implemented", func(t *testing.T) {
		ts := httptest.NewServer(newAppMux())
		defer ts.Close()
		client := ts.Client()
		urlline := ts.URL + path
		response, err := client.Post(urlline, "", nil)
		assert.Nil(t, err)

		assert.Equal(t, 405, response.StatusCode)
	})
}
