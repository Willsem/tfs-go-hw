package operation

type Operation struct {
	Type      *Type     `json:"type,omitempty"`
	Value     *Value    `json:"value,omitempty"`
	ID        *ID       `json:"id,omitempty"`
	CreatedAt *Datetime `json:"created_at,omitempty"`
}
