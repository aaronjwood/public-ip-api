package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"strings"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println("Failed to parse IP")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		forwarded := r.Header.Get("X-Forwarded-For")
		if forwarded != "" {
			originalIp := strings.Split(forwarded, ", ")[0]
			fmt.Fprintln(w, originalIp)
			return
		}

		fmt.Fprintln(w, ip)
	})

	log.Fatal(http.ListenAndServe(":8080", nil))
}
