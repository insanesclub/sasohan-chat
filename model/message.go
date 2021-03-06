package model

import "encoding/json"

// Message represents a message.
type Message struct {
	UserID     string `json:"user_id"`      // who sent this message
	ChatRoomID string `json:"chat_room_id"` // where to send this message
	Time       int64  `json:"time"`         // when this message sent, unix time
	Text       string `json:"text"`         // text content
}

// String implements fmt.Stringer.
func (m Message) String() string {
	b, err := json.Marshal(m)
	if err != nil {
		return err.Error()
	}
	return string(b)
}
