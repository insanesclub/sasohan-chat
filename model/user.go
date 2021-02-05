package model

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

// User represents a user.
type User struct {
	id    string // user ID
	conn  *websocket.Conn
	read  chan *Message
	quit  chan struct{}
	rooms map[*ChatRoom]struct{} // rooms participating
}

// NewUser returns a new user.
func NewUser(id string, conn *websocket.Conn) *User {
	return &User{
		id:    id,
		conn:  conn,
		read:  make(chan *Message),
		quit:  make(chan struct{}),
		rooms: make(map[*ChatRoom]struct{}),
	}
}

// Run starts user activation.
func (u *User) Run(users, rooms *sync.Map) {
	msg := new(Message)

	// catch and send message
	for {
		select {
		case msg := <-u.read:
			fmt.Printf("got message from user%s\n", msg.UserID)
			u.send(msg)
		case <-u.quit:
			for room := range u.rooms {
				delete(room.users, u)
			}
			users.Delete(u.id)
			close(u.read)
			close(u.quit)
			u.conn.Close()
			return
		default:
			if err := u.conn.ReadJSON(msg); err != nil {
				log.Println(err)
				u.Quit()
			} else {
				fmt.Println(msg)
			}
			if room, exists := rooms.Load(msg.ChatRoomID); exists {
				fmt.Printf("found room %s\n", room.(*ChatRoom).id)
				u.broadcast(room.(*ChatRoom), msg)
			}
		}
	}
}

// broadcast broadcasts msg to each users in room.
func (u *User) broadcast(room *ChatRoom, msg *Message) {
	fmt.Printf("broadcast to %v\n", room.users)
	for user := range room.users {
		go func(user *User) {
			user.read <- msg
		}(user)
	}
}

// send sends msg to client.
func (u *User) send(msg *Message) {
	fmt.Printf("send message to user%s\n", u.id)
	if err := u.conn.WriteJSON(*msg); err != nil {
		log.Println(err)
		u.Quit()
	} else {
		fmt.Println(msg)
	}
}

// Quit alerts user to quit.
func (u *User) Quit() {
	u.quit <- struct{}{}
	fmt.Printf("user%s quit\n", u.id)
}

// String implements fmt.Stringer.
func (u User) String() string {
	return fmt.Sprintf("user%s", u.id)
}
