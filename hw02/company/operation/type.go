package operation

import (
	"fmt"
	"strings"
)

type Type struct {
	string
	error
}

func (t *Type) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	switch s {
	case "income", "+":
		t.string = "+"
	case "outcome", "-":
		t.string = "-"
	default:
		t.error = fmt.Errorf("unexpected format of Type (value: %s)", s)
	}

	return nil
}
