// Go-to sample to get started on a simple webapp: https://golang.org/doc/articles/wiki/
package main

import (
	"io/ioutil"
	"log"
	"net/http"
)

// Serve single HTML page with embedded CSS and JS code
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("main.html")
	w.Write(body)
}

func main() {

	initWS()

	// Init web server
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/wschat", handleInitChat)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
