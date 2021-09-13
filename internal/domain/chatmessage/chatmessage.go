package chatmessage

import (
	"fmt"
	"time"

	"github.com/hackfeed/stark/internal/domain"
)

type chatMessage struct {
	Sent    time.Time
	Author  string
	Message string
}

func New(author, message string) domain.Messager {
	return &chatMessage{
		Sent:    time.Now(),
		Author:  author,
		Message: message,
	}
}

func (cm *chatMessage) GetSent() time.Time {
	return cm.Sent
}

func (cm *chatMessage) GetAuthor() string {
	return cm.Author
}

func (cm *chatMessage) GetMessage() string {
	return cm.Message
}

func (cm *chatMessage) String() string {
	return fmt.Sprintf("[%s] %s: %s", cm.Sent.Format("15:04:05"), cm.Author, cm.Message)
}
