package chat

import "github.com/hackfeed/stark/internal/domain"

type chat struct {
	Name     string
	Messages <-chan string
}

func New(name string) domain.Chatter {
	return &chat{
		Name:     name,
		Messages: make(<-chan string),
	}
}

func (c *chat) GetName() string {
	return c.Name
}

func (c *chat) GetMessages() <-chan string {
	return c.Messages
}

func (c *chat) SetMessages(messages <-chan string) {
	c.Messages = messages
}
