package operation

import (
	"fmt"
	"strconv"
	"strings"
)

type Value int

func (v *Value) UnmarshallJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	i, err := strconv.Atoi(s)
	if err != nil {
		return fmt.Errorf("unexpected format of Value: %s", err)
	}

	*v = Value(i)
	return nil
}
