package store

import "github.com/hackfeed/stark/internal/domain"

type MessagesRepository interface {
	GetMessages(Identifier) ([]domain.Messager, error)
	SetMessages(Identifier, []domain.Messager) error
}
