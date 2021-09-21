package operation

import (
	"fmt"
	"strings"
)

type Type struct {
	Value string
	Err   error
}

func (t *Type) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	switch s {
	case "income", "+":
		t.Value = "+"
	case "outcome", "-":
		t.Value = "-"
	default:
		t.Err = fmt.Errorf("unexpected format of Type (value: %s)", s)
	}

	return nil
}
