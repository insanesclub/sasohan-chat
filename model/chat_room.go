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
		if user, exists := users.Load(uid); exists && user.(*User) != nil {
			user.(*User).rooms[id] = struct{}{}
		}
	}
	return room
}

// Delete deletes c.
func (c *ChatRoom) Delete(users, rooms *sync.Map) {
	for uid := range c.users {
		if user, exists := users.Load(uid); exists && user.(*User) != nil {
			delete(user.(*User).rooms, c.id)
		}
	}
	c.users = nil
	rooms.Delete(c.id)
}
