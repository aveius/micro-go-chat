// Go-to sample to get started on a simple webapp: https://golang.org/doc/articles/wiki/
package main

import (
	"log"
	"net/http"

	"github.com/gobuffalo/packr"
)

var staticBox packr.Box

// Serve single HTML page with embedded CSS and JS code
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	body, err := staticBox.Find("main.html")
	if err != nil {
		log.Fatal(err)
	}
	w.Write(body)
}

func main() {
	// magic from https://github.com/gobuffalo/packr
	staticBox = packr.NewBox("./static")

	initWS()

	openDb()

	// Init web server
	http.Handle("/", mdwAuth(http.HandlerFunc(handleMainPage)))
	http.Handle("/last_messages", mdwAuth(http.HandlerFunc(handleLastMessages)))
	http.Handle("/wschat", mdwAuth(http.HandlerFunc(handleInitChat)))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
