package main

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strings"
)

func IpServer(w http.ResponseWriter, r *http.Request) {
	address := strings.Split(r.RemoteAddr, ":")[0]
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(struct {
		A string `json:"origin"`
	}{A: address})
}

func UserAgentServer(w http.ResponseWriter, r *http.Request) {
	userAgent := r.Header.Get("user-agent")
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(struct {
		UA string `json:"user-agent"`
	}{UA: userAgent})
}

func HeadersServer(w http.ResponseWriter, r *http.Request) {
	headers := make(map[string]string)
	for k, array := range r.Header {
		headers[k] = strings.Join(array, ",")
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(struct {
		H map[string]string `json:"headers"`
	}{H: headers})
}

func DeleteServer(w http.ResponseWriter, r *http.Request) {
	args := make(map[string]interface{})
	files := make(map[string]interface{})
	form := make(map[string]interface{})
	headers := make(map[string]string)
	origin := strings.Split(r.RemoteAddr, ":")[0]
	url := r.URL.String()
	for k, array := range r.URL.Query() {
		if len(array) == 1 {
			args[k] = array[0]
		} else {
			args[k] = array
		}
	}
	for k, array := range r.PostForm {
		if len(array) == 1 {
			form[k] = array[0]
		} else {
			form[k] = array
		}
	}
	for k, array := range r.Header {
		headers[k] = strings.Join(array, ";")
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(struct {
		Args    map[string]interface{} `json:"args"`
		Data    string                 `json:"data"`
		Files   map[string]interface{} `json:"files"`
		Form    map[string]interface{} `json:"form"`
		Headers map[string]string      `json:"headers"`
		Json    interface{}            `json:"json"`
		Origin  string                 `json:"origin"`
		Url     string                 `json:"url"`
	}{
		Args:    args,
		Data:    "",
		Files:   files,
		Form:    form,
		Headers: headers,
		Json:    nil,
		Origin:  origin,
		Url:     url,
	})
}

func PlayerServer(w http.ResponseWriter, r *http.Request) {
	player := r.PathValue("player")

	if player == "Pepper" {
		fmt.Fprint(w, "20")
		return
	}

	if player == "Floyd" {
		fmt.Fprint(w, "10")
		return
	}
}

func newAppMux() http.Handler {
	router := http.NewServeMux()
	router.HandleFunc("GET /players/{player}", PlayerServer)
	router.HandleFunc("GET /ip", IpServer)
	router.HandleFunc("GET /headers", HeadersServer)
	router.HandleFunc("GET /user-agent", UserAgentServer)
	router.HandleFunc("DELETE /delete", DeleteServer)
	return router
}

func main() {
	mux := newAppMux()

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		panic(err)
	}

	fmt.Println("Using port:", listener.Addr().(*net.TCPAddr).Port)

	panic(http.Serve(listener, mux))

}
