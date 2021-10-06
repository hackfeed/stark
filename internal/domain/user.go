package domain

type User struct {
	Name  string          `json:"name"`
	Chats map[string]Chat `json:"chats"`
}

func NewUser(name string) *User {
	return &User{
		Name:  name,
		Chats: make(map[string]Chat),
	}
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetChats() map[string]Chat {
	return u.Chats
}

func (u *User) AddChat(chat Chat) {
	u.GetChats()[chat.GetName()] = chat
}

func (u *User) RemoveChat(chatName string) {
	delete(u.GetChats(), chatName)
}
