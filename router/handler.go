package router

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/insanesclub/sasohan-chat/dataservice/dbutil"
	"github.com/insanesclub/sasohan-chat/model"
	"github.com/labstack/echo/v4"
)

// Connect creates a new connection.
func Connect(users, rooms *sync.Map, upgrader websocket.Upgrader) echo.HandlerFunc {
	return func(c echo.Context) error {
		conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
		if err != nil {
			return err
		}

		// get user ID
		body := new(struct {
			ID string `json:"id"`
		})

		// get payload from socket
		// payload is an in-memory buffer
		// and max size of the buffer is set by Conn.SetReadLimit
		// default max size is 4KB
		// see https://pkg.go.dev/github.com/gorilla/websocket
		//
		// if payload size exceeds limit, ErrReadLimit occurs
		if err = conn.ReadJSON(body); err != nil {
			return err
		}

		buf, err := json.Marshal(*body)
		if err != nil {
			return err
		}

		// restore unsent messages
		restoredMessages := new(struct {
			Messages []model.Message `json:"messages"`
			Success  bool            `json:"success"`
			ErrorMsg string          `json:"error_msg"`
		})

		dbutil.RestoreJSON(restoredMessages, "http://localhost:3000/restore", bytes.NewBuffer(buf))

		if !restoredMessages.Success {
			return errors.New(restoredMessages.ErrorMsg)
		}

		// send restored messages
		for _, msg := range restoredMessages.Messages {
			if err = conn.WriteJSON(msg); err != nil {
				return err
			}
		}

		// create a user
		user := model.NewUser(body.ID, conn)
		users.Store(body.ID, user)

		go user.Run(users, rooms)

		return nil
	}
}

// Disconnect closes a connection.
func Disconnect(users *sync.Map) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get user ID
		body := new(struct {
			ID string `json:"id"`
		})

		if err := c.Bind(body); err != nil {
			return err
		}

		// alert user to quit
		if user, exists := users.Load(body.ID); exists && user.(*model.User) != nil {
			user.(*model.User).Quit(users)
		}
		return nil
	}
}

// NewChat creates a new chat room.
func NewChat(users, rooms *sync.Map) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get user ID
		body := new(struct {
			ChatRoomID string   `json:"chat_room_id"`
			Users      []string `json:"users"`
		})

		if err := c.Bind(body); err != nil {
			return err
		}

		// create a chat room
		room := model.NewChatRoom(body.ChatRoomID, users, body.Users...)
		rooms.Store(body.ChatRoomID, room)

		return nil
	}
}

// LeaveChat lets the user leave the chat room.
func LeaveChat(users, rooms *sync.Map) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get user ID and chat room ID
		body := new(struct {
			UserID     string `json:"user_id"`
			ChatRoomID string `json:"chat_room_id"`
		})

		if err := c.Bind(body); err != nil {
			return err
		}

		user, exists := users.Load(body.UserID)
		if !exists || user.(*model.User) == nil {
			return fmt.Errorf("user %s does not exist", body.UserID)
		}

		room, exists := rooms.Load(body.ChatRoomID)
		if !exists {
			return fmt.Errorf("chat room %s does not exist", body.ChatRoomID)
		}

		user.(*model.User).Leave(room.(*model.ChatRoom))

		return nil
	}
}

// DeleteChat deletes the chat room.
func DeleteChat(users, rooms *sync.Map) echo.HandlerFunc {
	return func(c echo.Context) error {
		// get chat room ID
		body := new(struct {
			ChatRoomID string `json:"chat_room_id"`
		})

		if err := c.Bind(body); err != nil {
			return err
		}

		room, exists := rooms.Load(body.ChatRoomID)
		if !exists {
			return fmt.Errorf("chat room %s does not exist", body.ChatRoomID)
		}

		room.(*model.ChatRoom).Delete(users, rooms)

		return nil
	}
}
