package domain

type User interface {
	GetName() string
	GetChats() map[string]Chatter
	AddChat(Chatter)
	RemoveChat(string)
}
