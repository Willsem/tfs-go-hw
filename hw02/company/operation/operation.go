package operation

import "time"

type Operation struct {
	Type      Type      `json:"type,omitempty"`
	Value     Value     `json:"value,omitempty"`
	ID        ID        `json:"id,omitempty"`
	CreatedAt time.Time `json:"created_at,omitempty"`
}
