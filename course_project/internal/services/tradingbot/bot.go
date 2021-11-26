package tradingbot

import (
	"context"
	"fmt"

	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
	"github.com/willsem/tfs-go-hw/course_project/internal/repositories/applications"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/indicator"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/subscribe"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/telegram"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/trading"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/trading/tradingdto"
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
	buySize                uint64
	isWorking              bool
	cancelCtx              func()
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
		buySize:                1,
		isWorking:              false,
	}
}

func (bot *TradingBotImpl) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	bot.cancelCtx = cancel

	go bot.workWithTickers(ctx)
	bot.isWorking = true

	return nil
}

func (bot *TradingBotImpl) Stop() {
	bot.cancelCtx()
}

func (bot *TradingBotImpl) IsWorking() bool {
	return bot.isWorking
}

func (bot *TradingBotImpl) Continue() error {
	bot.isWorking = true
	return nil
}

func (bot *TradingBotImpl) Pause() error {
	bot.isWorking = false
	return nil
}

func (bot *TradingBotImpl) AddTicker(ticker string) error {
	err := bot.subscribeService.Subscribe(ticker, subscribe.Candle1m)
	if err != nil {
		return err
	}

	return nil
}

func (bot *TradingBotImpl) RemoveTicker(ticker string) error {
	err := bot.subscribeService.Unsubscribe(ticker, subscribe.Candle1m)
	if err != nil {
		return err
	}

	return nil
}

func (bot *TradingBotImpl) ChangeSize(newSize uint64) {
	bot.buySize = newSize
}

func (bot *TradingBotImpl) workWithTickers(ctx context.Context) {
	tickers := bot.subscribeService.GetChan()
	for {
		select {
		case <-ctx.Done():
			bot.subscribeService.Close()

		case ticker := <-tickers:
			decision := bot.indicatorService.MakeDecision(ticker)

			switch decision {
			case indicator.Buy:
				resp, err := bot.tradingService.SendOrder(tradingdto.Order{
					OrderType:  "ioc",
					Symbol:     ticker.ProductId,
					Side:       "buy",
					Size:       bot.buySize,
					LimitPrice: float64(ticker.Candle.Close) * 1.1,
				})
				if err != nil {
					bot.logger.Error(err)
				}

				if resp == trading.Placed {
					app := domain.Application{
						Ticker: ticker.ProductId,
						Cost:   float64(ticker.Candle.Close),
						Size:   bot.buySize,
						Type:   domain.Buy,
					}

					bot.applicationsRepository.Add(app)
					bot.telegramBot.SendSubscribedMessage(app.String())

					if _, ok := bot.tickers[ticker.ProductId]; !ok {
						bot.tickers[ticker.ProductId] = bot.buySize
					} else {
						bot.tickers[ticker.ProductId] += bot.buySize
					}
				} else {
					bot.logger.Error(
						fmt.Sprintf("can't buy %s by %f: %s", ticker.ProductId, ticker.Candle.Close, string(resp)),
					)
				}

			case indicator.Sell:
				if _, ok := bot.tickers[ticker.ProductId]; ok {
					var min uint64
					if bot.buySize < bot.tickers[ticker.ProductId] {
						min = bot.buySize
					} else {
						min = bot.tickers[ticker.ProductId]
					}

					resp, err := bot.tradingService.SendOrder(tradingdto.Order{
						OrderType:  "ioc",
						Symbol:     ticker.ProductId,
						Side:       "sell",
						Size:       min,
						LimitPrice: float64(ticker.Candle.Close) * 0.9,
					})
					if err != nil {
						bot.logger.Error(err)
					}

					if resp == trading.Cancelled {
						app := domain.Application{
							Ticker: ticker.ProductId,
							Cost:   float64(ticker.Candle.Close),
							Size:   min,
							Type:   domain.Buy,
						}

						bot.applicationsRepository.Add(app)
						bot.telegramBot.SendSubscribedMessage(app.String())
						bot.tickers[ticker.ProductId] -= min
					} else {
						bot.logger.Error(
							fmt.Sprintf("can't sell %s by %f: %s", ticker.ProductId, ticker.Candle.Close, string(resp)),
						)
					}
				}
			}
		}
	}
}
