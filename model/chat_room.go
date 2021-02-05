package model

// ChatRoom represents a chat room.
type ChatRoom struct {
	id    string             // chat room ID
	users map[*User]struct{} // users participating
}

// NewChatRoom returns a new chat room.
func NewChatRoom(id string, users ...*User) *ChatRoom {
	room := &ChatRoom{
		id:    id,
		users: make(map[*User]struct{}),
	}
	for _, user := range users {
		room.users[user] = struct{}{}
		user.rooms[room] = struct{}{}
	}
	return room
}
