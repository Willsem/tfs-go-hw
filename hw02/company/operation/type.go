package operation

import (
	"fmt"
	"strings"
)

type Type string

func (t *Type) UnmarshallJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	switch s {
	case "income", "+":
		*t = Type("+")
	case "outcome", "-":
		*t = Type("-")
	default:
		return fmt.Errorf("unexpected format of Type (value: %s)", s)
	}

	return nil
}
