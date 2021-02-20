package model

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/insanesclub/sasohan-chat/dataservice/dbutil"
)

// User represents a user.
type User struct {
	id    string // user ID
	conn  *websocket.Conn
	read  chan *Message
	quit  chan struct{}
	rooms map[string]struct{} // rooms participating
}

// NewUser returns a new user.
func NewUser(id string, conn *websocket.Conn) *User {
	return &User{
		id:    id,
		conn:  conn,
		read:  make(chan *Message),
		quit:  make(chan struct{}),
		rooms: make(map[string]struct{}),
	}
}

// Run starts user activation.
func (u *User) Run(users, rooms *sync.Map) {
	msg := new(Message)

	// catch and send message
	for {
		select {
		case m := <-u.read:
			u.send(m)
		case <-u.quit:
			u.Quit(users)
			return
		default:
			if err := u.conn.ReadJSON(msg); err != nil {
				log.Println(err)
				u.quit <- struct{}{}
			} else {
				if err = dbutil.StoreJSON(*msg, "http://localhost:3000/store"); err != nil {
					log.Println(err)
				}
				if room, exists := rooms.Load(msg.ChatRoomID); exists {
					u.broadcast(room.(*ChatRoom), msg, users)
				}
			}
		}
	}
}

// broadcast broadcasts msg to each users in room.
func (u *User) broadcast(room *ChatRoom, msg *Message, users *sync.Map) {
	for uid := range room.users {
		go func(uid string) {
			if user, exists := users.Load(uid); exists && user.(*User) != nil {
				user.(*User).read <- msg
			}
		}(uid)
	}
}

// send sends msg to client.
func (u *User) send(msg *Message) {
	// if message was not sent, store failed message
	if err := u.conn.WriteJSON(*msg); err != nil {
		body := struct {
			ID      string  `json:"id"`
			Message Message `json:"message"`
		}{ID: u.id, Message: *msg}

		if err = dbutil.StoreJSON(body, "http://localhost:3000/storeFailed"); err != nil {
			log.Println(err)
		}

		u.quit <- struct{}{}
	}
}

// Quit alerts u to quit.
func (u *User) Quit(users *sync.Map) {
	users.Delete(u.id)
	close(u.read)
	close(u.quit)
	u.rooms = nil
	u.conn.Close()
}

// Leave lets u leave room.
func (u *User) Leave(room *ChatRoom) {
	delete(room.users, u.id)
	delete(u.rooms, room.id)
}
