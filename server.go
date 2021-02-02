package main

import (
	"github.com/insanesclub/sasohan-chat/middleware"
	"github.com/insanesclub/sasohan-chat/model"
	"github.com/labstack/echo/v4"
	. "github.com/labstack/echo/v4/middleware"
)

func main() {
	// memoize chat rooms
	rooms := make(map[string]*model.ChatRoom)

	e := echo.New()
	e.Use(Logger())
	e.Use(Recover())
	e.POST("/enter", middleware.HandleConnection(rooms))
	e.Logger.Fatal(e.Start(":1323"))
}
