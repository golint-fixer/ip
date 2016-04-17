package main

import (
	"fmt"
	"net/http"

	"gopkg.in/vinxi/ip.v0"
	"gopkg.in/vinxi/vinxi.v0"
)

const port = 3100

func main() {
	// Create a new vinxi proxy
	vs := vinxi.NewServer(vinxi.ServerOptions{Port: port})

	// Create the IP filter
	filter := ip.New("127.0.0.1/8", "::1/64")

	// Attach a filter level middleware handler
	filter.Use(func(w http.ResponseWriter, r *http.Request, h http.Handler) {
		w.Header().Set("IP-Allowed", r.RemoteAddr)
		h.ServeHTTP(w, r)
	})

	// Register the filter in vinxi
	vs.Use(filter)

	// Target server to forward
	vs.Forward("http://httpbin.org")

	fmt.Printf("Server listening on port: %d\n", port)
	err := vs.Listen()
	if err != nil {
		fmt.Errorf("Error: %s\n", err)
	}
}
