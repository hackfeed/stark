package identifier

import (
	"fmt"
	"sort"
	"strings"

	"github.com/hackfeed/stark/internal/store"
)

type id struct {
	keyValues map[string]string
}

func New() store.Identifier {
	return &id{make(map[string]string)}
}

func (id *id) AddKeyValue(key, value string) store.Identifier {
	id.keyValues[key] = value
	return id
}

func (id *id) FormatID() string {
	keys := make([]string, len(id.keyValues))

	i := 0
	for key := range id.keyValues {
		keys[i] = key
		i++
	}
	sort.Strings(keys)

	ids := make([]string, len(id.keyValues)*2)
	i = 0
	for _, key := range keys {
		ids[i] = key
		ids[i+1] = id.keyValues[key]
		i += 2
	}

	return strings.Join(ids, ":")
}

func (id *id) FormatIDWithPostfix(postfix string) string {
	return fmt.Sprintf("%s:%s", id.FormatID(), postfix)
}

func (id *id) String() string {
	return id.FormatID()
}
