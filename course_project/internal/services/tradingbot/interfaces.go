package tradingbot

import "github.com/willsem/tfs-go-hw/course_project/internal/domain"

type TelegramBot interface {
	SendSubscribedMessage(message string)
}

type ApplicationsRepository interface {
	Add(application domain.Application) error
	GetAll() ([]domain.Application, error)
	GetByTicker(ticker string) ([]domain.Application, error)
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
