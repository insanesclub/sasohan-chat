package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/insanesclub/sasohan-chat/model"
	"golang.org/x/net/websocket"
)

func main() {
	origin := "http://localhost:1323/"
	url := "ws://localhost:1323/connect"
	conn, err := websocket.Dial(url, "", origin)

	if err != nil {
		log.Fatalf("connection failed: %v\n", err)
	}
	defer conn.Close()

	// request to create user with id
	body := struct {
		ID string `json:"id"`
	}{ID: "2"}

	if err = websocket.JSON.Send(conn, body); err != nil {
		log.Fatalf("error while sending message: %v\n", err)
	}

	// now user has been created
	fmt.Println("created user2")

	var msg model.Message

	// enjoy your chat!
	syncer := new(sync.WaitGroup)
	syncer.Add(2)

	go func() {
		defer syncer.Done()

		for {
			// receive
			if err := websocket.JSON.Receive(conn, &msg); err != nil {
				log.Fatalf("error while receiving message: %v\n", err)
			}
			fmt.Println(msg)
		}
	}()

	go func() {
		defer syncer.Done()

		for {
			// send
			fmt.Print("message: ")
			fmt.Scanf("%s", &msg.Text)
			msg.UserID = "2"
			msg.ChatRoomID = "A"
			msg.Time = time.Now().Unix()

			if err := websocket.JSON.Send(conn, msg); err != nil {
				log.Fatalf("error while sending message: %v\n", err)
			}
		}
	}()

	syncer.Wait()
}
