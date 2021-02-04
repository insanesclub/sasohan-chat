package main

import (
	"log"
	"net/http"
	"sync"

	"github.com/insanesclub/sasohan-chat/middleware"
)

func main() {
	// memoize chat rooms
	rooms := new(sync.Map)

	http.Handle("/connect", middleware.ConnectionHandlerGenerator(rooms))

	if err := http.ListenAndServe(":3000", nil); err != nil {
		panic(err)
	}
	if r := recover(); r != nil {
		log.Fatalln(r)
	}
}
