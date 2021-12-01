package tradingbot

import (
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
	"github.com/willsem/tfs-go-hw/course_project/internal/dto"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/tradingbot/mock_tradingbot"
)

func TestBotNothingDecision(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tickerInfo := domain.TickerInfo{
		Feed:      "candle",
		ProductId: "APPL",
		Candle: domain.Candle{
			Volume: 0,
			Open:   1,
			Close:  1,
			Low:    1,
			High:   1,
		},
	}

	mockTelegramBot := mock_tradingbot.NewMockTelegramBot(ctrl)

	mockSubscribe := mock_tradingbot.NewMockSubscribeService(ctrl)
	channel := make(chan domain.TickerInfo)
	mockSubscribe.EXPECT().GetChan().Return(channel)
	mockSubscribe.EXPECT().Subscribe("APPL", domain.Candle1m).Return(nil)
	mockSubscribe.EXPECT().Unsubscribe("APPL", domain.Candle1m).Return(nil)
	mockSubscribe.EXPECT().Close().Return(nil)

	mockRepo := mock_tradingbot.NewMockApplicationsRepository(ctrl)

	mockIndicator := mock_tradingbot.NewMockIndicatorService(ctrl)
	mockIndicator.EXPECT().MakeDecision(tickerInfo).Return(domain.NothingDecision)

	mockTrading := mock_tradingbot.NewMockTradingService(ctrl)

	mockLogger := mock_tradingbot.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(loggerServiceName, tickerInfo).Return()

	bot := New(mockSubscribe, mockTrading, mockIndicator, mockRepo, mockTelegramBot, mockLogger)
	err := bot.Start()
	bot.AddTicker("APPL")
	channel <- tickerInfo
	bot.RemoveTicker("APPL")
	bot.Stop()

	assert.Nil(t, err)

	time.Sleep(1 * time.Second)
}

func TestBotBuyAndSellDecision(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tickerInfoBuy := domain.TickerInfo{
		Feed:      "candle",
		ProductId: "APPL",
		Candle: domain.Candle{
			Volume: 0,
			Open:   1,
			Close:  1,
			Low:    1,
			High:   1,
		},
	}

	tickerInfoSell := domain.TickerInfo{
		Feed:      "candle",
		ProductId: "APPL",
		Candle: domain.Candle{
			Volume: 1,
			Open:   1,
			Close:  1,
			Low:    1,
			High:   1,
		},
	}

	orderBuy := dto.Order{
		OrderType:  "mkt",
		Symbol:     "APPL",
		Side:       "buy",
		Size:       1,
		LimitPrice: 0,
	}

	orderSell := dto.Order{
		OrderType:  "mkt",
		Symbol:     "APPL",
		Side:       "sell",
		Size:       1,
		LimitPrice: 0,
	}

	appBuy := domain.Application{
		Ticker: "APPL",
		Cost:   1,
		Size:   1,
		Type:   domain.BuyAppType,
	}

	appSell := domain.Application{
		Ticker: "APPL",
		Cost:   1,
		Size:   1,
		Type:   domain.SellAppType,
	}

	mockTelegramBot := mock_tradingbot.NewMockTelegramBot(ctrl)
	mockTelegramBot.EXPECT().SendSubscribedMessage(appBuy.String()).Return()
	mockTelegramBot.EXPECT().SendSubscribedMessage(appSell.String()).Return()

	mockSubscribe := mock_tradingbot.NewMockSubscribeService(ctrl)
	channel := make(chan domain.TickerInfo)
	mockSubscribe.EXPECT().GetChan().Return(channel)
	mockSubscribe.EXPECT().Subscribe("APPL", domain.Candle1m).Return(nil)
	mockSubscribe.EXPECT().Unsubscribe("APPL", domain.Candle1m).Return(nil)
	mockSubscribe.EXPECT().Close().Return(nil)

	mockRepo := mock_tradingbot.NewMockApplicationsRepository(ctrl)
	mockRepo.EXPECT().Add(appBuy).Return(nil)
	mockRepo.EXPECT().Add(appSell).Return(nil)

	mockIndicator := mock_tradingbot.NewMockIndicatorService(ctrl)
	mockIndicator.EXPECT().MakeDecision(tickerInfoBuy).Return(domain.BuyDecision)
	mockIndicator.EXPECT().MakeDecision(tickerInfoSell).Return(domain.SellDecision)

	mockTrading := mock_tradingbot.NewMockTradingService(ctrl)
	mockTrading.EXPECT().SendOrder(orderBuy).Return(domain.Placed, nil)
	mockTrading.EXPECT().SendOrder(orderSell).Return(domain.Placed, nil)

	mockLogger := mock_tradingbot.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(loggerServiceName, tickerInfoBuy).Return()
	mockLogger.EXPECT().Info(loggerServiceName, tickerInfoSell).Return()

	bot := New(mockSubscribe, mockTrading, mockIndicator, mockRepo, mockTelegramBot, mockLogger)
	err := bot.Start()
	bot.AddTicker("APPL")
	channel <- tickerInfoBuy
	channel <- tickerInfoSell
	time.Sleep(100 * time.Microsecond)
	bot.RemoveTicker("APPL")
	bot.Stop()

	assert.Nil(t, err)

	time.Sleep(1 * time.Second)
}

func TestBotSellDecision(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	tickerInfo := domain.TickerInfo{
		Feed:      "candle",
		ProductId: "APPL",
		Candle: domain.Candle{
			Volume: 0,
			Open:   1,
			Close:  1,
			Low:    1,
			High:   1,
		},
	}

	order := dto.Order{
		OrderType:  "mkt",
		Symbol:     "APPL",
		Side:       "sell",
		Size:       0,
		LimitPrice: 0,
	}

	buyError := errors.New("Argument can't be 0")

	mockTelegramBot := mock_tradingbot.NewMockTelegramBot(ctrl)

	mockSubscribe := mock_tradingbot.NewMockSubscribeService(ctrl)
	channel := make(chan domain.TickerInfo)
	mockSubscribe.EXPECT().GetChan().Return(channel)
	mockSubscribe.EXPECT().Subscribe("APPL", domain.Candle1m).Return(nil)
	mockSubscribe.EXPECT().Unsubscribe("APPL", domain.Candle1m).Return(nil)
	mockSubscribe.EXPECT().Close().Return(nil)

	mockRepo := mock_tradingbot.NewMockApplicationsRepository(ctrl)

	mockIndicator := mock_tradingbot.NewMockIndicatorService(ctrl)
	mockIndicator.EXPECT().MakeDecision(tickerInfo).Return(domain.SellDecision)

	mockTrading := mock_tradingbot.NewMockTradingService(ctrl)
	mockTrading.EXPECT().SendOrder(order).Return(domain.Empty, buyError)

	mockLogger := mock_tradingbot.NewMockLogger(ctrl)
	mockLogger.EXPECT().Info(loggerServiceName, tickerInfo).Return()
	mockLogger.EXPECT().Error(loggerServiceName, " error after sell a ticker: ", buyError).Return()

	bot := New(mockSubscribe, mockTrading, mockIndicator, mockRepo, mockTelegramBot, mockLogger)
	err := bot.Start()
	bot.AddTicker("APPL")
	channel <- tickerInfo
	bot.RemoveTicker("APPL")
	bot.Stop()

	assert.Nil(t, err)

	time.Sleep(1 * time.Second)
}

func TestBotSetings(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTelegramBot := mock_tradingbot.NewMockTelegramBot(ctrl)

	mockSubscribe := mock_tradingbot.NewMockSubscribeService(ctrl)
	channel := make(chan domain.TickerInfo)
	mockSubscribe.EXPECT().GetChan().Return(channel)
	mockSubscribe.EXPECT().Close().Return(nil)

	mockRepo := mock_tradingbot.NewMockApplicationsRepository(ctrl)

	mockIndicator := mock_tradingbot.NewMockIndicatorService(ctrl)

	mockTrading := mock_tradingbot.NewMockTradingService(ctrl)

	mockLogger := mock_tradingbot.NewMockLogger(ctrl)

	bot := New(mockSubscribe, mockTrading, mockIndicator, mockRepo, mockTelegramBot, mockLogger)
	err := bot.Start()
	assert.Equal(t, true, bot.IsWorking())
	bot.Pause()
	assert.Equal(t, false, bot.IsWorking())
	bot.Continue()
	assert.Equal(t, true, bot.IsWorking())
	bot.ChangeSize(10)
	assert.Equal(t, uint64(10), bot.buySize)
	bot.Stop()

	assert.Nil(t, err)

	time.Sleep(1 * time.Second)
}

func TestBotGetTickers(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockTelegramBot := mock_tradingbot.NewMockTelegramBot(ctrl)

	mockSubscribe := mock_tradingbot.NewMockSubscribeService(ctrl)
	channel := make(chan domain.TickerInfo)
	mockSubscribe.EXPECT().GetChan().Return(channel)
	mockSubscribe.EXPECT().Subscribe("APPL", domain.Candle1m).Return(nil)
	mockSubscribe.EXPECT().Unsubscribe("APPL", domain.Candle1m).Return(nil)
	mockSubscribe.EXPECT().Close().Return(nil)

	mockRepo := mock_tradingbot.NewMockApplicationsRepository(ctrl)

	mockIndicator := mock_tradingbot.NewMockIndicatorService(ctrl)

	mockTrading := mock_tradingbot.NewMockTradingService(ctrl)

	mockLogger := mock_tradingbot.NewMockLogger(ctrl)

	bot := New(mockSubscribe, mockTrading, mockIndicator, mockRepo, mockTelegramBot, mockLogger)
	err := bot.Start()
	tickers1 := bot.Tickers()
	bot.AddTicker("APPL")
	tickers2 := bot.Tickers()
	bot.RemoveTicker("APPL")
	tickers3 := bot.Tickers()
	bot.Stop()

	assert.Nil(t, err)
	assert.Equal(t, []string{}, tickers1)
	assert.Equal(t, []string{"APPL"}, tickers2)
	assert.Equal(t, []string{}, tickers3)

	time.Sleep(1 * time.Second)
}
