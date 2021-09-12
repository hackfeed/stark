package chatmessage

import (
	"time"

	"github.com/hackfeed/stark/internal/domain"
)

type ChatMessages []chatMessage

type chatMessage struct {
	LastUpdated time.Time `json:"last_updated"`
	Author      string    `json:"author"`
	Message     string    `json:"message"`
	IsEdited    bool      `json:"is_edited"`
}

func New(author, message string) domain.Messager {
	return &chatMessage{
		LastUpdated: time.Now(),
		Author:      author,
		Message:     message,
		IsEdited:    false,
	}
}

func (cm *chatMessage) GetLastUpdated() time.Time {
	return cm.LastUpdated
}

func (cm *chatMessage) GetAuthor() string {
	return cm.Author
}

func (cm *chatMessage) GetMessage() string {
	return cm.Message
}

func (cm *chatMessage) Edit(string) {
	cm.LastUpdated = time.Now()
	cm.Message += " (edited)"
	cm.IsEdited = true
}

func (cm *chatMessage) GetIsEdited() bool {
	return cm.IsEdited
}
