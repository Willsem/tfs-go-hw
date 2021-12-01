package dto

import (
	"time"
)

type Position struct {
	Side     string    `json:"side"`
	Symbol   string    `json:"symbol"`
	Price    float64   `json:"price"`
	FillTime time.Time `json:"fillTime"`
	Size     uint64    `json:"size"`
}
