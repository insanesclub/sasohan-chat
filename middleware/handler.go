package middleware

import (
	"github.com/insanesclub/sasohan-chat/model"
	"github.com/labstack/echo/v4"
	"golang.org/x/net/websocket"
)

// HandleConnection handles /enter.
// connection creates a new user,
// and users send messages each other.
func HandleConnection(rooms map[string]*model.ChatRoom) func(echo.Context) error {
	return func(c echo.Context) error {
		// parse user ID from request body
		body := new(struct {
			ID string `json:"id"`
		})
		if err := c.Bind(body); err != nil {
			return err
		}

		websocket.Handler(func(conn *websocket.Conn) {
			user := model.NewUser(body.ID, conn)
			go user.Run(rooms)
		}).ServeHTTP(c.Response(), c.Request())
		return nil
	}
}
