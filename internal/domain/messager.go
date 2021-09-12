package domain

import "time"

type Messager interface {
	GetLastUpdated() time.Time
	GetAuthor() string
	GetMessage() string
	GetIsEdited() bool
	Edit(string)
}
