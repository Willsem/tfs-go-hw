package domain

import "time"

type Message struct {
	SenderID   string
	RecipentID string
	Content    string
	DateTime   time.Time
}
