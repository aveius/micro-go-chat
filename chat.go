// Contain all the logic to handle websockets
// https://godoc.org/github.com/gorilla/websocket
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// type Message struct {
// 	Author, Body string
// }

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var openedWS []*websocket.Conn

// Init global states
func initWS() {
	openedWS = make([]*websocket.Conn, 0, 10)
}

// Handles WS upgrade requests
func handleInitChat(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	// Upgrade successful! Register our new WS in the pool of WSs, and start playing with it
	openedWS = append(openedWS, conn)
	runChat(conn)
	for i, el := range openedWS {
		if el == conn {
			openedWS[i] = openedWS[len(openedWS)-1]
			openedWS[len(openedWS)-1] = nil
			openedWS = openedWS[:len(openedWS)-1]
			break
		}
	}
}

// Wonderland within the websocket connection
func runChat(conn *websocket.Conn) {
	log.Printf("We're on! (%d connected clients) \n", len(openedWS))
	for {
		// FIXME simple ping-pong, returns all received messages
		// FIXME need mutex, read/writes can't be concurrent!
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// Try to notify all opened WSs
		log.Println("WS in: " + string(p))
		messageAllWS(messageType, p)
	}
}

func messageAllWS(messageType int, message []byte) {

	for _, conn := range openedWS {
		log.Println("Sending to WS " + conn.RemoteAddr().String())
		if err := conn.WriteMessage(messageType, message); err != nil {
			log.Println(err)
			return
		}
	}

}
