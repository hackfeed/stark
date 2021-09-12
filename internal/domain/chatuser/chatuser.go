package chatuser

import "github.com/hackfeed/stark/internal/domain"

type chatUser struct {
	Name  string
	Chats map[string]domain.Chatter
}

func New(name string) domain.User {
	return &chatUser{
		Name:  name,
		Chats: make(map[string]domain.Chatter),
	}
}

func (cu *chatUser) GetName() string {
	return cu.Name
}

func (cu *chatUser) GetChats() map[string]domain.Chatter {
	return cu.Chats
}

func (cu *chatUser) AddChat(chat domain.Chatter) {
	cu.GetChats()[chat.GetName()] = chat
}

func (cu *chatUser) RemoveChat(chatName string) {
	delete(cu.GetChats(), chatName)
}