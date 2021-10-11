package domain

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
)

type CandleFilename string

const (
	CandleFile1m  CandleFilename = "candles_1m.csv"
	CandleFile2m  CandleFilename = "candles_2m.csv"
	CandleFile10m CandleFilename = "candles_10m.csv"

	LenOfBuffer int = 100
)

type Pipeline struct {
	logger *log.Logger
	ctx    context.Context
}

func NewPipeline(logger *log.Logger) Pipeline {
	return Pipeline{
		logger: logger,
	}
}

func (p Pipeline) Start(in <-chan Price, ctx context.Context) error {
	ch := p.calculate1mCandles(in)
	ch, err := p.writeCandlesToFile(ch, CandleFile1m)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	ch = p.calculateFewMCandles(ch, CandlePeriod2m)
	ch, err = p.writeCandlesToFile(ch, CandleFile2m)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	ch = p.calculateFewMCandles(ch, CandlePeriod10m)
	ch, err = p.writeCandlesToFile(ch, CandleFile10m)
	if err != nil {
		p.logger.Error(err)
		return err
	}

	for range ctx.Done() {
		break
	}

	p.logger.Info("End writing to files")
	return nil
}

func (p Pipeline) writeCandlesToFile(in <-chan Candle, filename CandleFilename) (<-chan Candle, error) {
	out := make(chan Candle, LenOfBuffer)
	file, err := os.Create(string(filename))
	if err != nil {
		return nil, err
	}

	go func() {
		defer close(out)
		defer file.Close()

		for candle := range in {
			out <- candle
			p.logger.Infof("write to file %s: %+v", filename, candle)

			_, err := fmt.Fprintf(
				file,
				"%s,%s,%.6f,%.6f,%.6f,%.6f\n",
				candle.Ticker,
				candle.TS.Format(time.RFC3339),
				candle.Open,
				candle.Close,
				candle.High,
				candle.Low,
			)

			if err != nil {
				p.logger.Error(err)
				return
			}
		}
	}()

	return out, nil
}

func (p Pipeline) calculate1mCandles(in <-chan Price) <-chan Candle {
	out := make(chan Candle, LenOfBuffer)

	go func() {
		defer close(out)

		companies := make(map[string][]Price)
		for price := range in {
			p.logger.Debugf("was read price: %+v", price)

			current, ok := companies[price.Ticker]
			if !ok {
				current = make([]Price, 0, 1)
				current = append(current, price)
				companies[price.Ticker] = current
			} else {
				prevTS, _ := PeriodTS(CandlePeriod1m, current[len(current)-1].TS)
				newTS, _ := PeriodTS(CandlePeriod1m, price.TS)

				if prevTS == newTS {
					current = append(current, price)
					companies[price.Ticker] = current
				} else {
					out <- makeCandleFromPrices(current, prevTS)

					current = make([]Price, 0, 1)
					current = append(current, price)
					companies[price.Ticker] = current
				}
			}
		}

		p.logger.Info("end reading prices")

		for _, prices := range companies {
			ts, _ := PeriodTS(CandlePeriod1m, prices[0].TS)
			out <- makeCandleFromPrices(prices, ts)
		}

		p.logger.Infof("end sending from 1m")
	}()

	return out
}

func makeCandleFromPrices(prices []Price, ts time.Time) Candle {
	candle := Candle{
		Ticker: prices[0].Ticker,
		Period: CandlePeriod1m,
		TS:     ts,
		Open:   prices[0].Value,
		Close:  prices[len(prices)-1].Value,
		High:   prices[0].Value,
		Low:    prices[0].Value,
	}

	for _, p := range prices {
		if p.Value > candle.High {
			candle.High = p.Value
		}
		if p.Value < candle.Low {
			candle.Low = p.Value
		}
	}

	return candle
}

func (p Pipeline) calculateFewMCandles(in <-chan Candle, period CandlePeriod) <-chan Candle {
	out := make(chan Candle, LenOfBuffer)

	go func() {
		defer close(out)

		companies := make(map[string][]Candle)
		for candle := range in {
			p.logger.Debugf("was read candle for %s: %+v", period, candle)

			current, ok := companies[candle.Ticker]
			if !ok {
				current = make([]Candle, 0, 1)
				current = append(current, candle)
				companies[candle.Ticker] = current
			} else {
				prevTS, _ := PeriodTS(period, current[len(current)-1].TS)
				newTS, _ := PeriodTS(period, candle.TS)

				if prevTS == newTS {
					current = append(current, candle)
					companies[candle.Ticker] = current
				} else {
					out <- makeCandleFromCandles(current, period, prevTS)

					current = make([]Candle, 0, 1)
					current = append(current, candle)
					companies[candle.Ticker] = current
				}
			}
		}

		p.logger.Infof("end reading candles for: %s", period)

		for _, candles := range companies {
			ts, _ := PeriodTS(period, candles[0].TS)
			out <- makeCandleFromCandles(candles, period, ts)
		}

		p.logger.Infof("end sending from %s", period)
	}()

	return out
}

func makeCandleFromCandles(candles []Candle, period CandlePeriod, ts time.Time) Candle {
	candleToOut := Candle{
		Ticker: candles[0].Ticker,
		Period: period,
		TS:     ts,
		Open:   candles[0].Open,
		Close:  candles[len(candles)-1].Close,
		High:   candles[0].High,
		Low:    candles[0].Low,
	}

	for _, c := range candles {
		if c.High > candleToOut.High {
			candleToOut.High = c.High
		}
		if c.Low < candleToOut.Low {
			candleToOut.Low = c.Low
		}
	}

	return candleToOut
}
