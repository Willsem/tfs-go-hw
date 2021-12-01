package tradingbot

import (
	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
	"github.com/willsem/tfs-go-hw/course_project/internal/dto"
)

type TelegramBot interface {
	SendSubscribedMessage(message string)
}

type ApplicationsRepository interface {
	Add(application domain.Application) error
}

type IndicatorService interface {
	MakeDecision(ticker domain.TickerInfo) domain.Decision
}

type SubscribeService interface {
	GetChan() <-chan domain.TickerInfo
	Subscribe(ticker string, candle domain.CandleType) error
	Unsubscribe(ticker string, candle domain.CandleType) error
	Close() error
}

type TradingService interface {
	OpenPositions() ([]dto.Position, error)
	SendOrder(order dto.Order) (domain.OrderStatus, error)
}
