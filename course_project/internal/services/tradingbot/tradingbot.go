package tradingbot

import (
	"context"
	"fmt"
	"sync"

	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
	"github.com/willsem/tfs-go-hw/course_project/internal/dto"
	"github.com/willsem/tfs-go-hw/course_project/pkg/log"
)

const (
	loggerServiceName = "[TradingBot]"
)

type TradingBot struct {
	subscribeService       SubscribeService
	tradingService         TradingService
	indicatorService       IndicatorService
	applicationsRepository ApplicationsRepository
	telegramBot            TelegramBot
	logger                 log.Logger
	tickers                map[string]uint64
	buySize                uint64
	isWorking              bool
	cancelCtx              func()
	sizeMutex              *sync.RWMutex
	tickersMutex           *sync.RWMutex
	isWorkingMutex         *sync.RWMutex
}

func New(
	subscribeService SubscribeService,
	tradingService TradingService,
	indicatorService IndicatorService,
	applicationsRepository ApplicationsRepository,
	telegramBot TelegramBot,
	logger log.Logger,
) *TradingBot {
	return &TradingBot{
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
		tickersMutex:           &sync.RWMutex{},
		isWorkingMutex:         &sync.RWMutex{},
	}
}

func (bot *TradingBot) Start() error {
	ctx, cancel := context.WithCancel(context.Background())
	bot.cancelCtx = cancel

	go bot.workWithTickers(ctx)
	bot.isWorking = true

	return nil
}

func (bot *TradingBot) Stop() {
	bot.cancelCtx()
}

func (bot *TradingBot) IsWorking() bool {
	bot.isWorkingMutex.RLock()
	defer bot.isWorkingMutex.RUnlock()
	return bot.isWorking
}

func (bot *TradingBot) Continue() error {
	bot.isWorkingMutex.Lock()
	bot.isWorking = true
	bot.isWorkingMutex.Unlock()
	return nil
}

func (bot *TradingBot) Pause() error {
	bot.isWorkingMutex.Lock()
	bot.isWorking = false
	bot.isWorkingMutex.Unlock()
	return nil
}

func (bot *TradingBot) Tickers() []string {
	bot.tickersMutex.RLock()
	result := make([]string, 0, len(bot.tickers))
	for key := range bot.tickers {
		result = append(result, key)
	}
	bot.tickersMutex.RUnlock()

	return result
}

func (bot *TradingBot) AddTicker(ticker string) error {
	err := bot.subscribeService.Subscribe(ticker, domain.Candle1m)
	if err != nil {
		return err
	}

	bot.tickersMutex.Lock()
	bot.tickers[ticker] = 0
	bot.tickersMutex.Unlock()

	return nil
}

func (bot *TradingBot) RemoveTicker(ticker string) error {
	err := bot.subscribeService.Unsubscribe(ticker, domain.Candle1m)
	if err != nil {
		return err
	}

	bot.tickersMutex.Lock()
	if _, ok := bot.tickers[ticker]; ok {
		delete(bot.tickers, ticker)
	}
	bot.tickersMutex.Unlock()

	return nil
}

func (bot *TradingBot) ChangeSize(newSize uint64) {
	bot.sizeMutex.Lock()
	bot.buySize = newSize
	bot.sizeMutex.Unlock()
}

func (bot *TradingBot) workWithTickers(ctx context.Context) {
	tickers := bot.subscribeService.GetChan()

	for {
		select {
		case <-ctx.Done():
			bot.subscribeService.Close()
			return

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

func (bot *TradingBot) buyTicker(ticker domain.TickerInfo) {
	bot.sizeMutex.RLock()
	size := bot.buySize
	bot.sizeMutex.RUnlock()

	resp, err := bot.tradingService.SendOrder(dto.Order{
		OrderType:  "mkt",
		Symbol:     ticker.ProductId,
		Side:       "buy",
		Size:       size,
		LimitPrice: 0,
	})

	if err != nil {
		bot.logger.Error(loggerServiceName, " error after buy a ticker: ", err)
	} else {
		if resp == domain.Placed {
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

			bot.tickersMutex.Lock()
			if _, ok := bot.tickers[ticker.ProductId]; !ok {
				bot.tickers[ticker.ProductId] = size
			} else {
				bot.tickers[ticker.ProductId] += size
			}
			bot.tickersMutex.Unlock()
		} else {
			bot.logger.Error(loggerServiceName,
				fmt.Sprintf("can't buy %s by %f: %s", ticker.ProductId, ticker.Candle.Close, string(resp)),
			)
		}
	}
}

func (bot *TradingBot) sellTicker(ticker domain.TickerInfo) {
	if _, ok := bot.tickers[ticker.ProductId]; ok {
		bot.sizeMutex.RLock()
		size := bot.buySize
		bot.sizeMutex.RUnlock()

		bot.tickersMutex.RLock()
		var min uint64
		if size < bot.tickers[ticker.ProductId] {
			min = size
		} else {
			min = bot.tickers[ticker.ProductId]
		}
		bot.tickersMutex.RUnlock()

		resp, err := bot.tradingService.SendOrder(dto.Order{
			OrderType:  "mkt",
			Symbol:     ticker.ProductId,
			Side:       "sell",
			Size:       min,
			LimitPrice: 0,
		})
		if err != nil {
			bot.logger.Error(loggerServiceName, " error after sell a ticker: ", err)
		} else {
			if resp == domain.Placed {
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

				bot.tickersMutex.Lock()
				bot.tickers[ticker.ProductId] -= min
				bot.tickersMutex.Unlock()
			} else {
				bot.logger.Error(loggerServiceName,
					fmt.Sprintf("can't sell %s by %f: %s", ticker.ProductId, ticker.Candle.Close, string(resp)),
				)
			}
		}
	}
}
