package model

// ChatRoom represents a chat room.
type ChatRoom struct {
	ID    string  // identifier
	Users []*User // participants
}

// NewChatRoom
func NewChatRoom(id string) *ChatRoom {
	return &ChatRoom{
		ID:    id,
		Users: make([]*User, 2), // there's always at least 2 users in a chat room
	}
}
