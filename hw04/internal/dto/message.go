package dto

import "time"

type Message struct {
	User     User      `json:"user"`
	Content  string    `json:"content"`
	DateTime time.Time `json:"date_time"`
}
