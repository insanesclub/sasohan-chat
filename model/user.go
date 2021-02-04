package model

import (
	"sync"

	"golang.org/x/net/websocket"
)

// User represents a user.
type User struct {
	id   string // identifier
	conn *websocket.Conn
	read chan *Message
	quit chan struct{}
}

// NewUser returns a new user.
func NewUser(id string, conn *websocket.Conn) *User {
	return &User{
		id:   id,
		conn: conn,
		read: make(chan *Message),
		quit: make(chan struct{}),
	}
}

// Run starts user activation.
func (u *User) Run(rooms *sync.Map) {
	go u.listen(rooms)

	defer u.exit()

	// catch and send message
	// TODO: when user logouts, close connection and channels
	for {
		select {
		case msg := <-u.read:
			u.send(msg)
		case <-u.quit:
			return
		}
	}
}

// listen listens message from client.
// when message comes, broadcasts it to each users in the chat room.
func (u *User) listen(rooms *sync.Map) error {
	msg := new(Message)
	for {
		// parse message and broadcast it
		if err := websocket.Message.Receive(u.conn, msg); err != nil {
			return err
		}
		if room, exists := rooms.Load(msg.ChatRoomID); exists {
			u.broadcast(room.(*ChatRoom), msg)
		}
	}

	// TODO: when u quits, this goroutine should stop (context?)
}

// broadcast broadcasts msg to each users in room.
func (u *User) broadcast(room *ChatRoom, msg *Message) {
	for user := range room.Users {
		user.read <- msg
	}
}

// send sends msg to client.
// if error occurs or msg was sent less, returns it.
func (u *User) send(msg *Message) error {
	if err := websocket.JSON.Send(u.conn, *msg); err != nil {
		return err
	}
	return nil
}

// exit removes user from chat rooms and close its connection, channels.
func (u *User) exit() {
	// TODO
}
