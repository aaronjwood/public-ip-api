package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"net"
)

const (
	port = 5001
)

// Our server that detects the connecting public IP.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			// If we could not split on an IP:Port then give back the raw remote address contained in the request.
			fmt.Fprintln(w, r.RemoteAddr)
			return
		}

		// Proxy awareness.
		// The official defined format is a comma space separated list of IP addresses.
		forwarded := r.Header.Get("X-Forwarded-For")
		if forwarded != "" {
			originalIP := strings.Split(forwarded, ", ")[0]
			fmt.Fprintln(w, originalIP)
			return
		}

		// If no proxy was found then give the parsed IP.
		fmt.Fprintf(w, ip)
	})

	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
