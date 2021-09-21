package operation

import (
	"fmt"
	"strconv"
	"strings"
)

type Value struct {
	Value int
	Err   error
}

func (v *Value) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	i, err := strconv.Atoi(s)
	if err != nil {
		v.Err = fmt.Errorf("unexpected format of Value: %s", err)
	}

	v.Value = i
	return nil
}
