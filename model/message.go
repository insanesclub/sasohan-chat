package model

import (
	"fmt"
	"time"
)

// Message represents a message.
type Message struct {
	UserID     string `json:"user_id"`      // who sent this message
	ChatRoomID string `json:"chat_room_id"` // where to send this message
	Time       int64  `json:"time"`         // when this message sent, unix time
	Text       string `json:"text"`         // text content
}

// String implements fmt.Stringer.
func (m Message) String() string {
	return fmt.Sprintf("from: %s\nto: %s\nwhen: %v\ntext: %s",
		m.UserID,
		m.ChatRoomID,
		time.Unix(m.Time, 0),
		m.Text)
}
