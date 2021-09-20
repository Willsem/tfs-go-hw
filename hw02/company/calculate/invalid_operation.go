package calculate

import (
	"fmt"
	"strconv"
)

type InvalidOperationID struct {
	*int
	*string
}

func (id InvalidOperationID) MarshalJSON() ([]byte, error) {
	switch {
	case id.int != nil:
		return []byte(strconv.Itoa(*id.int)), nil
	case id.string != nil:
		return []byte(*id.string), nil
	default:
		return nil, fmt.Errorf("Cannot marshal empty ID")
	}
}
