package domain

import (
	"fmt"
	"time"
)

type applicationType string

const (
	Buy  applicationType = "buy"
	Sell applicationType = "sell"
)

type Application struct {
	Id        uint64
	Ticker    string
	Cost      int64
	Size      int64
	CreatedAt time.Time
	Type      applicationType
}

func (app Application) String() string {
	return fmt.Sprintf("%s $%s\nСтоимость: %d\nРазмер: %d\nВремя заявки: %s\n\n",
		app.Type, app.Ticker, app.Cost, app.Size, app.CreatedAt.Format("01-02-2006 15:04:05"))
}

func (t applicationType) String() string {
	switch t {
	case Buy:
		return "Покупка"
	case Sell:
		return "Продажа"
	default:
		return "Неизвестный тип заявки"
	}
}
