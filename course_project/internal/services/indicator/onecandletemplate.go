package indicator

import (
	"time"

	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
)

type OneCandleTemplate struct {
	candles map[string][]domain.Candle
}

func NewOneCandleTemplate() *OneCandleTemplate {
	return &OneCandleTemplate{
		candles: make(map[string][]domain.Candle),
	}
}

func (service *OneCandleTemplate) MakeDecision(ticker domain.TickerInfo) domain.Decision {
	if _, ok := service.candles[ticker.ProductId]; !ok {
		service.candles[ticker.ProductId] = make([]domain.Candle, 0, initialSize)
		return domain.NothingDecision
	}

	candles := service.candles[ticker.ProductId]
	n := len(candles)
	if n > 0 && time.Time(candles[n-1].Time).Equal(time.Time(ticker.Candle.Time)) {
		return domain.NothingDecision
	}

	service.candles[ticker.ProductId] = append(candles, ticker.Candle)
	return service.findTemplates(service.candles[ticker.ProductId])
}

func (service *OneCandleTemplate) findTemplates(candles []domain.Candle) domain.Decision {
	candle := candles[len(candles)-1]

	if candle.Open < candle.Close {
		return domain.SellDecision
	}

	if candle.Open > candle.Close {
		return domain.BuyDecision
	}

	return domain.NothingDecision
}
