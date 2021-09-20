package operation

import (
	"strings"
)

type ID struct {
	string
	error
}

func (id *ID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	id.string = s
	return nil
}
