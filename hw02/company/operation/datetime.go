package operation

import (
	"fmt"
	"strings"
	"time"
)

type Datetime struct {
	time.Time
	error
}

func (date *Datetime) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		date.error = fmt.Errorf("%s", err)
	}

	date.Time = t
	return nil
}

func (date Datetime) MarshallJSON() ([]byte, error) {
	return nil, nil
}
