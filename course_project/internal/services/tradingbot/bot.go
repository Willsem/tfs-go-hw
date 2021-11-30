package tradingbot

import (
	"context"
	"fmt"
	"sync"

	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/subscribe"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/trading"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/trading/tradingdto"
	"github.com/willsem/tfs-go-hw/course_project/pkg/log"
)

const (
	loggerServiceName = "[TradingBot]"
)

type TradingBotImpl struct {
	subscribeService       subscribe.SubscribeService
	tradingService         trading.TradingService
	indicatorService       IndicatorService
	applicationsRepository ApplicationsRepository
	telegramBot            TelegramBot
	logger                 log.Logger
	tickers                map[string]uint64
	buySize                uint64
	isWorking              bool
	cancelCtx              func()
	sizeMutex              *sync.RWMutex
	isWorkingMutex         *sync.RWMutex
}

func New(
	subscribeService subscribe.SubscribeService,
	tradingService trading.TradingService,
	indicatorService IndicatorService,
	applicationsRepository ApplicationsRepository,
	telegramBot TelegramBot,
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
		sizeMutex:              &sync.RWMutex{},
		isWorkingMutex:         &sync.RWMutex{},
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
	bot.isWorkingMutex.RLock()
	defer bot.isWorkingMutex.RUnlock()
	return bot.isWorking
}

func (bot *TradingBotImpl) Continue() error {
	bot.isWorkingMutex.Lock()
	bot.isWorking = true
	bot.isWorkingMutex.Unlock()
	return nil
}

func (bot *TradingBotImpl) Pause() error {
	bot.isWorkingMutex.Lock()
	bot.isWorking = false
	bot.isWorkingMutex.Unlock()
	return nil
}

func (bot *TradingBotImpl) Tickers() []string {
	result := make([]string, 0, len(bot.tickers))
	for key := range bot.tickers {
		result = append(result, key)
	}
	return result
}

func (bot *TradingBotImpl) AddTicker(ticker string) error {
	err := bot.subscribeService.Subscribe(ticker, subscribe.Candle1m)
	if err != nil {
		return err
	}

	bot.tickers[ticker] = 0

	return nil
}

func (bot *TradingBotImpl) RemoveTicker(ticker string) error {
	err := bot.subscribeService.Unsubscribe(ticker, subscribe.Candle1m)
	if err != nil {
		return err
	}

	if _, ok := bot.tickers[ticker]; ok {
		delete(bot.tickers, ticker)
	}

	return nil
}

func (bot *TradingBotImpl) ChangeSize(newSize uint64) {
	bot.sizeMutex.Lock()
	bot.buySize = newSize
	bot.sizeMutex.Unlock()
}

func (bot *TradingBotImpl) workWithTickers(ctx context.Context) {
	tickers := bot.subscribeService.GetChan()

	for {
		select {
		case <-ctx.Done():
			bot.subscribeService.Close()

		case ticker := <-tickers:
			bot.isWorkingMutex.RLock()
			isWorking := bot.isWorking
			bot.isWorkingMutex.RUnlock()

			bot.logger.Info(loggerServiceName, ticker)

			decision := bot.indicatorService.MakeDecision(ticker)

			if isWorking {
				switch decision {
				case domain.BuyDecision:
					bot.buyTicker(ticker)

				case domain.SellDecision:
					bot.sellTicker(ticker)
				}
			}
		}
	}
}

func (bot *TradingBotImpl) buyTicker(ticker domain.TickerInfo) {
	bot.sizeMutex.RLock()
	size := bot.buySize
	bot.sizeMutex.RUnlock()

	resp, err := bot.tradingService.SendOrder(tradingdto.Order{
		OrderType:  "mkt",
		Symbol:     ticker.ProductId,
		Side:       "buy",
		Size:       size,
		LimitPrice: 0,
	})

	if err != nil {
		bot.logger.Error(loggerServiceName, " error after buy a ticker: ", err)
	} else {
		if resp == trading.Placed {
			app := domain.Application{
				Ticker: ticker.ProductId,
				Cost:   float64(ticker.Candle.Close),
				Size:   size,
				Type:   domain.BuyAppType,
			}

			err = bot.applicationsRepository.Add(app)
			if err != nil {
				bot.logger.Error(loggerServiceName, " applications repository add: ", err)
			}

			bot.telegramBot.SendSubscribedMessage(app.String())

			if _, ok := bot.tickers[ticker.ProductId]; !ok {
				bot.tickers[ticker.ProductId] = size
			} else {
				bot.tickers[ticker.ProductId] += size
			}
		} else {
			bot.logger.Error(loggerServiceName,
				fmt.Sprintf("can't buy %s by %f: %s", ticker.ProductId, ticker.Candle.Close, string(resp)),
			)
		}
	}
}

func (bot *TradingBotImpl) sellTicker(ticker domain.TickerInfo) {
	if _, ok := bot.tickers[ticker.ProductId]; ok {
		bot.sizeMutex.RLock()
		size := bot.buySize
		bot.sizeMutex.RUnlock()

		var min uint64
		if size < bot.tickers[ticker.ProductId] {
			min = size
		} else {
			min = bot.tickers[ticker.ProductId]
		}

		resp, err := bot.tradingService.SendOrder(tradingdto.Order{
			OrderType:  "mkt",
			Symbol:     ticker.ProductId,
			Side:       "sell",
			Size:       min,
			LimitPrice: 0,
		})
		if err != nil {
			bot.logger.Error(loggerServiceName, " error after sell a ticker: ", err)
		} else {
			if resp == trading.Placed {
				app := domain.Application{
					Ticker: ticker.ProductId,
					Cost:   float64(ticker.Candle.Close),
					Size:   min,
					Type:   domain.SellAppType,
				}

				err = bot.applicationsRepository.Add(app)
				if err != nil {
					bot.logger.Error(loggerServiceName, " applications repository add: ", err)
				}

				bot.telegramBot.SendSubscribedMessage(app.String())
				bot.tickers[ticker.ProductId] -= min
			} else {
				bot.logger.Error(loggerServiceName,
					fmt.Sprintf("can't sell %s by %f: %s", ticker.ProductId, ticker.Candle.Close, string(resp)),
				)
			}
		}
	}
}
