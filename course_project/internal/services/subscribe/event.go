package subscribe

import "github.com/willsem/tfs-go-hw/course_project/internal/domain"

type event struct {
	Event      string        `json:"event,omitempty"`
	Feed       string        `json:"feed,omitempty"`
	ProductIds []string      `json:"product_ids,omitempty"`
	ProductId  string        `json:"product_id,omitempty"`
	Message    string        `json:"message,omitempty"`
	Candle     domain.Candle `json:"candle,omitempty"`
}
