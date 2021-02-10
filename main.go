package main

import (
	"sync"

	"github.com/gorilla/websocket"
	"github.com/insanesclub/sasohan-chat/router"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Map is concurrent-safe hash map.
	// see https://pkg.go.dev/sync#Map

	// memoize users[user ID, user]
	users := new(sync.Map)

	// memoize chat rooms[chat room ID, chat room]
	rooms := new(sync.Map)

	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// upgrader holds websocket connection
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1 << 12,
		WriteBufferSize: 1 << 12,
	}

	// register URIs
	e.GET("/connect", router.Connect(users, rooms, upgrader))
	e.POST("/disconnect", router.Disconnect(users))
	e.POST("/newchat", router.NewChat(users, rooms))
	e.POST("/leave", router.Leave(users, rooms))
	e.POST("/deletechat", router.DeleteChat(rooms))

	// start server
	e.Logger.Fatal(e.Start(":1323"))
}
