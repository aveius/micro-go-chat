// Contain all the logic to handle websockets
// https://godoc.org/github.com/gorilla/websocket
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Message struct {
	Author, Body string
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func handleInitChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	runChat(conn)
}

// Wonderland within the websocket connection
func runChat(conn *websocket.Conn) {
	for {
		// FIXME simple ping-pong, returns all received messages
		// FIXME need mutex, read/writes can't be concurrent!
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}
	}
}
