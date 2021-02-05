package middleware

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/insanesclub/sasohan-chat/model"
)

// Connect creates a new connection.
func Connect(users, rooms *sync.Map, upgrader websocket.Upgrader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Fatalf("upgrader.Upgrade: %v\n", err)
		}

		// get user ID
		body := new(struct {
			ID string `json:"id"`
		})

		// get payload from socket
		// payload is in-memory buffer
		// and max size of the buffer is set by Conn.MaxPayloadBytes
		// default max size is 32MB
		// see https://pkg.go.dev/github.com/gorilla/websocket
		//
		// if payload size exceeds limit, ErrFrameTooLarge is returned
		if err = conn.ReadJSON(body); err != nil {
			log.Fatalf("error while getting user ID: %v\n", err)
		}

		// create user
		user := model.NewUser(body.ID, conn)
		fmt.Printf("created user%s\n", body.ID)

		users.Store(body.ID, user)
		go user.Run(users, rooms)
	}
}

// Disconnect closes a connection.
func Disconnect(users *sync.Map) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user ID
		body := new(struct {
			ID string `json:"id"`
		})

		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			log.Fatalln(err)
		}

		// alert user to quit
		if user, exists := users.Load(body.ID); exists {
			user.(*model.User).Quit()
		}
	}
}

// NewChat creates a new chat room
func NewChat(users, rooms *sync.Map) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// get user ID
		body := new(struct {
			ChatRoomID string   `json:"chat_room_id"`
			Users      []string `json:"users"`
		})

		if err := json.NewDecoder(r.Body).Decode(body); err != nil {
			log.Fatalln(err)
		}

		// create new chat room
		us := make([]*model.User, len(body.Users))
		for i, uid := range body.Users {
			if user, exists := users.Load(uid); exists {
				us[i] = user.(*model.User)
			}
		}

		room := model.NewChatRoom(body.ChatRoomID, us...)
		rooms.Store(body.ChatRoomID, room)

		fmt.Printf("created chat room %s\n", body.ChatRoomID)
		fmt.Printf("users in chat room %s: %v\n", body.ChatRoomID, room)
	}
}
