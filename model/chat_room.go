package model

// ChatRoom represents a chat room.
type ChatRoom struct {
	ID    string             // identifier
	Users map[*User]struct{} // user set
}

// NewChatRoom returns a new chat room.
func NewChatRoom(id string) *ChatRoom {
	return &ChatRoom{
		ID:    id,
		Users: make(map[*User]struct{}),
	}
}
