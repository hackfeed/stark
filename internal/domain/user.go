package domain

type User struct {
	Name  string
	Chats map[string]*Chat
}

func NewUser(name string) *User {
	return &User{
		Name:  name,
		Chats: make(map[string]*Chat),
	}
}

func (u *User) GetName() string {
	return u.Name
}

func (u *User) GetChats() map[string]*Chat {
	return u.Chats
}

func (u *User) AddChat(chat *Chat) {
	if u.GetChats() != nil && u.GetActiveChat() != nil {
		u.GetActiveChat().SetIsActive(false)
	}
	u.GetChats()[chat.GetName()] = chat
}

func (u *User) RemoveChat(chatName string) {
	if u.GetChats()[chatName].GetIsActive() {
		for cn, c := range u.GetChats() {
			if cn != chatName {
				c.SetIsActive(true)
				break
			}
		}
	}
	delete(u.GetChats(), chatName)
}

func (u *User) GetActiveChat() *Chat {
	for _, chat := range u.GetChats() {
		if chat.GetIsActive() {
			return chat
		}
	}
	return nil
}

func (u *User) SetActiveChat(chatName string) {
	u.GetActiveChat().SetIsActive(false)
	u.GetChats()[chatName].SetIsActive(true)
}

func (u *User) GetInactiveChats() map[string]*Chat {
	chats := make(map[string]*Chat)
	for chatName, chat := range u.GetChats() {
		if !chat.GetIsActive() {
			chats[chatName] = chat
		}
	}
	return chats
}
