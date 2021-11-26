package tradingbot

import (
	"github.com/willsem/tfs-go-hw/course_project/internal/repositories/applications"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/indicator"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/subscribe"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/telegram"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/trading"
	"github.com/willsem/tfs-go-hw/course_project/pkg/log"
)

type TradingBotImpl struct {
	subscribeService       subscribe.SubscribeService
	tradingService         trading.TradingService
	indicatorService       indicator.IndicatorService
	applicationsRepository applications.ApplicationsRepository
	telegramBot            telegram.Bot
	logger                 log.Logger
	tickers                map[string]uint64
}

func New(
	subscribeService subscribe.SubscribeService,
	tradingService trading.TradingService,
	indicatorService indicator.IndicatorService,
	applicationsRepository applications.ApplicationsRepository,
	telegramBot telegram.Bot,
	logger log.Logger,
) *TradingBotImpl {
	return &TradingBotImpl{
		subscribeService:       subscribeService,
		tradingService:         tradingService,
		indicatorService:       indicatorService,
		applicationsRepository: applicationsRepository,
		telegramBot:            telegramBot,
		logger:                 logger,
		tickers:                make(map[string]uint64),
	}
}

func (bot *TradingBotImpl) Start() error {
	return nil
}

func (bot *TradingBotImpl) Continue() error {
	return nil
}

func (bot *TradingBotImpl) Pause() error {
	return nil
}

func (bot *TradingBotImpl) AddTicker(ticker string) error {
	bot.telegramBot.SendSubscribedMessage(ticker)
	return nil
}

func (bot *TradingBotImpl) RemoveTicker(ticker string) error {
	return nil
}
