package filesrepo

type FilesRepository interface {
	GetFiles(string) (map[string][]byte, error)
	SetFiles(string, map[string][]byte) error
}
