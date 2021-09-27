package operation

import (
	"strings"
)

type ID struct {
	Value string
	Err   error
}

func (id *ID) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	id.Value = s
	return nil
}
