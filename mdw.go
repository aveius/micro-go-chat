// Auth middleware for our handlers
package main

import (
	"log"
	"net"
	"net/http"
	"os"
)

var whitelistedIP string

// As usual, thanks for the pointers:
// https://www.alexedwards.net/blog/making-and-using-middleware

// Simple middleware limiting access to the configured IP. If none is configured,
// will restrict to localhost calls.
func mdwAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Fatal(err)
		}
		if ip != whitelistedIP {
			log.Printf("(W) Rejected client %s\n", ip)
			http.Error(w, "", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func init() {
	whitelistedIP = os.Getenv("WHITELISTED_IP")
	if whitelistedIP == "" {
		whitelistedIP = "127.0.0.1"
	}
}
