package domain

import (
	"fmt"
	"time"
)

type Messager interface {
	fmt.Stringer

	GetSent() time.Time
	GetAuthor() string
	GetMessage() string
}
