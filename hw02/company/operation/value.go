package operation

import (
	"fmt"
	"strconv"
	"strings"
)

type Value struct {
	int
	error
}

func (v *Value) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	i, err := strconv.Atoi(s)
	if err != nil {
		v.error = fmt.Errorf("unexpected format of Value: %s", err)
	}

	v.int = i
	return nil
}
