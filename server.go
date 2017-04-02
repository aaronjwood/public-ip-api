package main

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/appengine"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}

		forwarded := r.Header.Get("X-Forwarded-For")
		if forwarded != "" {
			originalIp := strings.Split(forwarded, ", ")[0]
			fmt.Fprintln(w, originalIp)
			return
		}

		fmt.Fprintln(w, r.RemoteAddr)
	})

	appengine.Main()
}
