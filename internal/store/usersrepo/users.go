package usersrepo

type UsersRepository interface {
	GetUsers(string) (map[string]struct{}, error)
	SetUsers(string, map[string]struct{}) error
}
