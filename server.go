package main

import (
	"fmt"
	"strings"
	"encoding/json"
	"net"
	"net/http"
)

func IpServer(w http.ResponseWriter, r *http.Request) {
	address := strings.Split(r.RemoteAddr, ":")[0]
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(struct{
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
		fmt.Println(k , array)
		headers[k] = strings.Join(array, ",")
	}
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "  ")
	encoder.Encode(struct{
		H map[string]string `json:"headers"`
	}{H: headers})
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

func newAppMux() http.Handler{
	router := http.NewServeMux()
	router.HandleFunc("GET /players/{player}", PlayerServer)
	router.HandleFunc("GET /ip", IpServer)
	router.HandleFunc("GET /headers", HeadersServer)
	router.HandleFunc("GET /user-agent", UserAgentServer)
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
