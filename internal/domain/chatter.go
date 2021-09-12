package domain

type Chatter interface {
	GetName() string
	GetUsers() map[string]User
	GetMessages() []Messager
	AddUser(User)
	RemoveUser(string)
}
