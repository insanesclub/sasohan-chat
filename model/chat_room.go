package model

import "sync"

// ChatRoom represents a chat room.
type ChatRoom struct {
	id    string              // chat room ID
	users map[string]struct{} // users participating
}

// NewChatRoom returns a new chat room.
func NewChatRoom(id string, users *sync.Map, uids ...string) *ChatRoom {
	room := &ChatRoom{
		id:    id,
		users: make(map[string]struct{}),
	}
	for _, uid := range uids {
		room.users[uid] = struct{}{}
	}
	return room
}

// Delete deletes c.
func (c *ChatRoom) Delete(users, rooms *sync.Map) {
	c.users = nil
	rooms.Delete(c.id)
}
