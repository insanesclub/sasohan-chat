package middleware

import (
	"sync"

	"github.com/insanesclub/sasohan-chat/model"
	"golang.org/x/net/websocket"
)

// ConnectionHandlerGenerator generates connection handler.
func ConnectionHandlerGenerator(rooms *sync.Map) websocket.Handler {
	return func(c *websocket.Conn) {
		// get user ID
		body := new(struct {
			ID string `json:"id"`
		})

		// get payload from socket
		// payload is in-memory buffer
		// and max size of the buffer is set by Conn.MaxPayloadBytes
		// default max size is 32MB
		// see https://pkg.go.dev/golang.org/x/net/websocket
		//
		// if payload size exceeds limit, ErrFrameTooLarge is returned
		if err := websocket.JSON.Receive(c, body); err != nil {
			panic(err)
		}

		// create user with connection
		user := model.NewUser(body.ID, c)
		go user.Run(rooms)
	}
}
