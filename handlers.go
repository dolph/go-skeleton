package main

import (
	"fmt"
	"net"
	"net/http"
)

// Returns the client's IP address.
func echoRemoteAddr(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		fmt.Fprintf(w, "%s\n", r.RemoteAddr)
	} else {
		fmt.Fprintf(w, "%s\n", ip)
	}
}

// Returns a 404 Not Found page.
func notFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	fmt.Fprintf(w, "%d Not Found\n", http.StatusNotFound)
}
