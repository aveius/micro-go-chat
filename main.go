// Go-to sample to get started on a simple webapp: https://golang.org/doc/articles/wiki/
package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

var db *sql.DB

// Serve single HTML page with embedded CSS and JS code
func handleMainPage(w http.ResponseWriter, r *http.Request) {
	body, _ := ioutil.ReadFile("main.html")
	w.Write(body)
}

func main() {
	// Init DB
	// FIXME use envvar
	connStr := "postgres://postgres:*******g@127.0.0.1:5432/microchat"
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}

	initWS()

	// Init web server
	http.HandleFunc("/", handleMainPage)
	http.HandleFunc("/wschat", handleInitChat)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
