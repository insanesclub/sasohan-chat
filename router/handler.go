package router

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
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
		// payload is in-memory buffer
		// and max size of the buffer is set by Conn.SetReadLimit
		// default max size is 4KB
		// see https://pkg.go.dev/github.com/gorilla/websocket
		//
		// if payload size exceeds limit, ErrReadLimit occurs
		if err = conn.ReadJSON(body); err != nil {
			log.Fatalf("error while getting user ID: %v\n", err)
		}

		// create user
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
		if user, exists := users.Load(body.ID); exists {
			user.(*model.User).Quit()
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

		// create new chat room
		us := make([]*model.User, len(body.Users))
		for i, uid := range body.Users {
			if user, exists := users.Load(uid); exists {
				us[i] = user.(*model.User)
			}
		}

		room := model.NewChatRoom(body.ChatRoomID, us...)
		rooms.Store(body.ChatRoomID, room)
		return nil
	}
}

// Leave lets the user leave the chat room.
func Leave(users, rooms *sync.Map) echo.HandlerFunc {
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
		if !exists {
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
func DeleteChat(rooms *sync.Map) echo.HandlerFunc {
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

		room.(*model.ChatRoom).Delete()
		rooms.Delete(body.ChatRoomID)
		return nil
	}
}
