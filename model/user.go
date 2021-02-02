package model

import (
	"encoding/json"
	"errors"

	"golang.org/x/net/websocket"
)

// User represents a user.
type User struct {
	id   string // identifier
	conn *websocket.Conn
	read chan Message
	quit chan struct{}
}

// NewUser returns a new user with connection c.
func NewUser(id string, conn *websocket.Conn) *User {
	return &User{
		id:   id,
		conn: conn,
		read: make(chan Message),
		quit: make(chan struct{}),
	}
}

// Run starts user activation.
func (u *User) Run(rooms map[string]*ChatRoom) {
	go u.receive(rooms)

	// catch message and send
	// if user logout, then quit and close connection.
	for {
		select {
		case msg := <-u.read:
			u.send(msg)
		case <-u.quit:
			u.conn.Close()
			return
		}
	}
}

// receive listens message from socket.
// when message comes, broadcasts it to users in the chat room.
func (u *User) receive(rooms map[string]*ChatRoom) error {
	msg := new(Message)
	for {
		// parse message and broadcast it
		if err := websocket.Message.Receive(u.conn, msg); err != nil {
			return err
		}
		u.broadCast(rooms[msg.ChatRoomID], msg)
	}
}

// broadcast delivers msg to each users in room.
func (u *User) broadCast(room *ChatRoom, msg *Message) {
	for _, user := range room.Users {
		user.read <- *msg
	}
}

// send sends msg to client.
// if error occurs or msg was sent less, returns it.
func (u *User) send(msg Message) error {
	m, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	n, err := u.conn.Write(m)
	if err != nil {
		return err
	}
	if n < len(m) {
		return errors.New("less message sent")
	}
	return nil
}

// exit occurs when user logout.
// alert user to quit and close connection.
func (u *User) exit() { u.quit <- struct{}{} }
