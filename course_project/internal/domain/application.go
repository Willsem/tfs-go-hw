package domain

import (
	"fmt"
	"time"
)

type applicationType string

const (
	BuyAppType  applicationType = "buy"
	SellAppType applicationType = "sell"
)

type Application struct {
	Id        uint64
	Ticker    string
	Cost      float64
	Size      uint64
	CreatedAt time.Time
	Type      applicationType
}

func (app Application) String() string {
	if app.CreatedAt.IsZero() {
		return fmt.Sprintf("%s $%s\nСтоимость: %f\nРазмер: %d\n\n",
			app.Type, app.Ticker, app.Cost, app.Size)
	} else {
		return fmt.Sprintf("%s $%s\nСтоимость: %f\nРазмер: %d\nВремя заявки: %s\n\n",
			app.Type, app.Ticker, app.Cost, app.Size, app.CreatedAt.Format("01-02-2006 15:04:05"))
	}
}

func (t applicationType) String() string {
	switch t {
	case BuyAppType:
		return "Покупка"
	case SellAppType:
		return "Продажа"
	default:
		return "Неизвестный тип заявки"
	}
}
