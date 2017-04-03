package main

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine"
	"net"
)

// Our server that detects the connecting public IP.
func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		// Depending on our App Engine deployment type we may or may not have anything to split on.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err == nil {

			// Proxy awareness.
			// The official defined format is a comma space separated list of IP addresses.
			forwarded := r.Header.Get("X-Forwarded-For")
			if forwarded != "" {
				originalIp := strings.Split(forwarded, ", ")[0]
				fmt.Fprintln(w, originalIp)
				return
			}

			// If no proxy was found then give the parsed IP.
			fmt.Fprintf(w, ip)
			return
		}

		// If we could not split on an IP:Port then give back the raw remote address contained in the request.
		fmt.Fprintln(w, r.RemoteAddr)
	})

	appengine.Main()
}
