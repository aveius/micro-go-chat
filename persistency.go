package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Handles DB connection and I/O

var db *sql.DB

type message struct {
	Author string `json:"author"`
	Body   string `json:"body"`
}

type datedMessage struct {
	message
	Timestamp int64 `json:"date"`
}

// Opens connection to DB, and initializes schema & table if necessary
func openDb() {
	// FIXME use envvar
	var err error
	db, err = sql.Open("postgres", os.Getenv("PG_CONN_URL"))
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		log.Printf("Connectivity to DB failed, persistency disabled: %s", err)
		db.Close()
		db = nil
		return
	}

	_, err = db.Exec("CREATE SCHEMA IF NOT EXISTS microchat;")
	if err != nil {
		log.Printf("DB: schema creation error, %s", err)
		db.Close()
		db = nil
		return
	}

	_, err = db.Exec("CREATE TABLE IF NOT EXISTS microchat.messages (" +
		"date TIMESTAMP DEFAULT CURRENT_TIMESTAMP," +
		"author TEXT," +
		"message TEXT" +
		")")
	if err != nil {
		log.Printf("DB: table creation error, %s", err)
		db.Close()
		db = nil
		return
	}

	log.Println("DB connected")
	// defer db.Close()
}

// Takes a JSON message from the frontend, and serializes it into the DB
func saveMessage(p []byte) {
	if db == nil {
		return
	}

	var msg message
	json.Unmarshal(p, &msg)

	stmt, err := db.Prepare("INSERT INTO microchat.messages(author, message) VALUES($1, $2)")
	if err != nil {
		log.Fatal(err)
	}
	_, err = stmt.Exec(msg.Author, msg.Body)
	if err != nil {
		log.Fatal(err)
	}
}

// if DB available, retrieve last 5m messages (20 at most), return as JSON
func handleLastMessages(w http.ResponseWriter, r *http.Request) {
	if db == nil {
		http.Error(w, "DB unavailable", http.StatusServiceUnavailable)
		return
	}

	var res = make([]datedMessage, 0, 20)

	rows, err := db.Query("SELECT * FROM microchat.messages " +
		"WHERE date >= CURRENT_TIMESTAMP - interval '5 minutes' " +
		"ORDER BY date DESC " +
		"LIMIT 20")
	// TODO probably should be a wrapper on all queries; downgrade gracefully on any error
	if err != nil {
		log.Printf("DB: data retrieval failed, %s", err)
		db.Close()
		db = nil

		http.Error(w, "DB error", http.StatusServiceUnavailable)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var m datedMessage
		var t time.Time // We need a temporary time type for timestamp retrieval from PG
		err := rows.Scan(&t, &m.Author, &m.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Javascript uses timestamps with millis
		m.Timestamp = t.UnixNano() / 1e6

		res = append(res, m)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Sending %d recent messages\n", len(res))
	out, err := json.Marshal(res)
	if err != nil {
		log.Fatal(err)
	}
	w.Write(out)
}
