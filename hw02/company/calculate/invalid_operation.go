package calculate

import (
	"fmt"
	"strconv"
)

type InvalidOperationID struct {
	*int
	*string
}

func NewInvalidOperationID(id string) InvalidOperationID {
	opID := InvalidOperationID{}

	val, err := strconv.Atoi(id)
	if err != nil {
		opID.int = &val
		opID.string = nil
	} else {
		opID.int = nil
		opID.string = &id
	}

	return opID
}

func (id InvalidOperationID) MarshalJSON() ([]byte, error) {
	switch {
	case id.string != nil:
		return []byte(*id.string), nil
	case id.int != nil:
		return []byte(strconv.Itoa(*id.int)), nil
	default:
		return nil, fmt.Errorf("Cannot marshal empty ID")
	}
}
