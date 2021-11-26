package subscribe

import "github.com/willsem/tfs-go-hw/course_project/internal/domain"

type candleType string

type SubscribeService interface {
	GetChan() <-chan domain.TickerInfo
	Subscribe(ticker string, candle candleType) error
	Unsubscribe(ticker string, candle candleType) error
	Close() error
}
