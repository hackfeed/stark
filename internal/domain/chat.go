package domain

type Chat struct {
	Name     string
	Messages <-chan string
}

func NewChat(name string) *Chat {
	return &Chat{
		Name:     name,
		Messages: make(<-chan string),
	}
}

func (c *Chat) GetName() string {
	return c.Name
}

func (c *Chat) GetMessages() <-chan string {
	return c.Messages
}

func (c *Chat) SetMessages(messages <-chan string) {
	c.Messages = messages
}
