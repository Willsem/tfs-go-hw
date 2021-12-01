package indicator

import (
	"time"

	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
)

const (
	initialSize = 100
)

type TripleCandlesTemplate struct {
	candles map[string][]domain.Candle
}

func NewTripleCandlesTemplate() *TripleCandlesTemplate {
	return &TripleCandlesTemplate{
		candles: make(map[string][]domain.Candle),
	}
}

func (service *TripleCandlesTemplate) MakeDecision(ticker domain.TickerInfo) domain.Decision {
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

	if len(service.candles[ticker.ProductId]) < 3 {
		return domain.NothingDecision
	}

	return service.findTemplates(service.candles[ticker.ProductId])
}

func (service *TripleCandlesTemplate) findTemplates(candles []domain.Candle) domain.Decision {
	red := 0
	green := 0
	for _, candle := range candles[len(candles)-3:] {
		if candle.Open < candle.Close {
			green++
		} else if candle.Close < candle.Open {
			red++
		}
	}

	if red == 3 {
		return domain.BuyDecision
	}

	if green == 3 {
		return domain.SellDecision
	}

	return domain.NothingDecision
}
