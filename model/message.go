package model

// Message represents a message from user.
type Message struct {
	// who sent this message
	UserID string `json:"user_id"`
	// where to send this message
	ChatRoomID string `json:"chat_room_id"`
	// when this message sent, unix time
	Time uint64 `json:"time"`
	// text content
	Content string `json:"content"`
}
