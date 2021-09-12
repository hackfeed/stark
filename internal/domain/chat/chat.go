package chat

import "github.com/hackfeed/stark/internal/domain"

type chat struct {
	Name     string
	Users    map[string]domain.User
	Messages []domain.Messager
}

func New(name string) domain.Chatter {
	return &chat{
		Name:     name,
		Users:    make(map[string]domain.User),
		Messages: make([]domain.Messager, 0),
	}
}

func (c *chat) GetName() string {
	return c.Name
}

func (c *chat) GetUsers() map[string]domain.User {
	return c.Users
}

func (c *chat) GetMessages() []domain.Messager {
	return c.Messages
}

func (c *chat) AddUser(user domain.User) {
	c.GetUsers()[user.GetName()] = user
}

func (c *chat) RemoveUser(name string) {
	delete(c.GetUsers(), name)
}
