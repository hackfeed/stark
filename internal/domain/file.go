package domain

type File struct {
	Name    string
	Content []byte
}

func NewFile(name string, content []byte) *File {
	return &File{
		Name:    name,
		Content: content,
	}
}

func (f *File) GetName() string {
	return f.Name
}

func (f *File) GetContent() []byte {
	return f.Content
}
