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
