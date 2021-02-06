package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
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
	}{ID: "1"}

	if err = websocket.JSON.Send(conn, body); err != nil {
		log.Fatalf("error while sending message: %v\n", err)
	}

	// now user has been created
	fmt.Println("created user1")

	var msg model.Message

	// if user1 and user2 are set, create chat room
	fmt.Print("Ready? ")
	fmt.Scanf("%s", &msg.Text)

	// request to create a new chat room
	rBody := struct {
		ID    string   `json:"chat_room_id"`
		Users []string `json:"users"`
	}{ID: "A", Users: []string{"1", "2"}}

	m, err := json.Marshal(rBody)
	if err != nil {
		log.Fatalf("error while marshaling json: %v\n", err)
	}
	buff := bytes.NewBuffer(m)
	_, err = http.Post("http://localhost:1323/newchat", "application/json", buff)
	if err != nil {
		log.Fatalf("error while creating chat room: %v\n", err)
	}

	// now chat room has been created
	fmt.Println("created chat room A")

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
			msg.UserID = "1"
			msg.ChatRoomID = "A"
			msg.Time = time.Now().Unix()

			if err := websocket.JSON.Send(conn, msg); err != nil {
				log.Fatalf("error while sending message: %v\n", err)
			}
		}
	}()

	syncer.Wait()
}
