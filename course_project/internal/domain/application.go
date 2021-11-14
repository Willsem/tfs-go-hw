package domain

import "time"

type Application struct {
	Id        uint64
	Ticker    string
	Cost      int64
	CreatedAt time.Time
}
