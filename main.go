package main

import (
	"io"
	"log"
	"net/http"
	"time"
)

// upstream is set at build time:
//
//	go build -ldflags "-s -w -X main.upstream=sdk-server-nlb:9090"
var upstream string

func main() {
	if upstream == "" {
		log.Fatal("upstream not set; build with -ldflags \"-X main.upstream=<host:port>\"")
	}

	// single shared client
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		resp, err := client.Get("http://" + upstream)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// simple stream copy; Content-Type is always text/plain from mock-server
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		_, _ = io.Copy(w, resp.Body)
	})

	const port = "8080"
	log.Printf("mock-gateway proxy ready → %s on :%s", upstream, port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
