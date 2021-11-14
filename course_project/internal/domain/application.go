package domain

import (
	"fmt"
	"time"
)

type Application struct {
	Id        uint64
	Ticker    string
	Cost      int64
	CreatedAt time.Time
}

func (app Application) String() string {
	return fmt.Sprintf("[%s] Стоимость: %d, Время заявки: %s\n",
		app.Ticker, app.Cost, app.CreatedAt.Format("01-02-2006 15:04:05"))
}
