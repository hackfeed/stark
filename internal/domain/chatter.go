package domain

type Chatter interface {
	GetName() string
	GetMessages() <-chan string
	SetMessages(<-chan string)
}
