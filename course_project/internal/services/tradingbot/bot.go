package tradingbot

import (
	"github.com/willsem/tfs-go-hw/course_project/internal/services/indicator"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/subscribe"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/trading"
	"github.com/willsem/tfs-go-hw/course_project/pkg/log"
)

type TradingBotImpl struct {
	subscribeService subscribe.SubscribeService
	tradingService   trading.TradingService
	indicatorService indicator.IndicatorService
	logger           log.Logger
}

func New(
	subscribeService subscribe.SubscribeService,
	tradingService trading.TradingService,
	indicatorService indicator.IndicatorService,
	logger log.Logger,
) *TradingBotImpl {
	return &TradingBotImpl{
		subscribeService: subscribeService,
		tradingService:   tradingService,
		indicatorService: indicatorService,
		logger:           logger,
	}
}

func (bot *TradingBotImpl) Start() error {
	return nil
}

func (bot *TradingBotImpl) Configure(command Command) error {
	return nil
}
