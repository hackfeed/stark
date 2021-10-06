package domain

import (
	"fmt"
	"time"
)

type Message struct {
	Sent    time.Time
	Author  string
	Message string
}

func NewMessage(author, message string) *Message {
	return &Message{
		Sent:    time.Now(),
		Author:  author,
		Message: message,
	}
}

func (m *Message) GetSent() time.Time {
	return m.Sent
}

func (m *Message) GetAuthor() string {
	return m.Author
}

func (m *Message) GetMessage() string {
	return m.Message
}

func (m *Message) String() string {
	return fmt.Sprintf("[%s] %s: %s", m.Sent.Format("15:04:05"), m.Author, m.Message)
}
