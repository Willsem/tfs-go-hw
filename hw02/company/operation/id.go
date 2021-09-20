package operation

import (
	"strings"
)

type ID string

func (id *ID) UnmarshallJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	*id = ID(s)
	return nil
}
