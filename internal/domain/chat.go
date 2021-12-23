package domain

type Chat struct {
	Name     string
	Messages <-chan string
	IsActive bool
	Buffer   []string
}

func NewChat(name string) *Chat {
	return &Chat{
		Name:     name,
		Messages: make(<-chan string),
		IsActive: true,
		Buffer:   nil,
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

func (c *Chat) GetIsActive() bool {
	return c.IsActive
}

func (c *Chat) SetIsActive(isActive bool) {
	c.IsActive = isActive
}

func (c *Chat) GetBuffer() []string {
	return c.Buffer
}

func (c *Chat) SetBuffer(buffer []string) {
	c.Buffer = buffer
}
