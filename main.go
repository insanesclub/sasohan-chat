package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/insanesclub/sasohan-chat/middleware"
)

// upgrader holds websocket connection
var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
}

func main() {
	// Map is concurrent-safe hash map.
	// see https://pkg.go.dev/sync#Map

	// memoize users[user ID, user]
	users := new(sync.Map)

	// memoize chat rooms[chat room ID, chat room]
	rooms := new(sync.Map)

	// register URIs
	http.Handle("/connect", middleware.Connect(users, rooms, upgrader))
	http.Handle("/disconnect", middleware.Disconnect(users))
	http.Handle("/newchat", middleware.NewChat(users, rooms))

	// start server
	if err := http.ListenAndServe(":3000", nil); err != nil {
		log.Fatalln(err)
	}
}
