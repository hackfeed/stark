package store

import "fmt"

type Identifier interface {
	fmt.Stringer

	AddKeyValue(key, value string) Identifier
	FormatID() string
	FormatIDWithPostfix(salt string) string
}
