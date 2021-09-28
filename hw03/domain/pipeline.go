package domain

import (
	"context"
	"fmt"
	"os"
	"sync"
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

func NewPipeline(logger *log.Logger, ctx context.Context) Pipeline {
	return Pipeline{
		logger: logger,
		ctx:    ctx,
	}
}

func (p Pipeline) Start(in <-chan Price, wg *sync.WaitGroup) {
	defer wg.Done()

	wgFiles := &sync.WaitGroup{}

	minuteOut := p.calculate1mCandles(in)

	file1mChan := make(chan Candle, LenOfBuffer)
	in2mChan := make(chan Candle, LenOfBuffer)

	go func() {
		for candle := range minuteOut {
			file1mChan <- candle
			in2mChan <- candle
		}

		close(file1mChan)
		close(in2mChan)
	}()

	wgFiles.Add(1)
	go p.writeCandlesToFile(file1mChan, CandleFile1m, wgFiles)

	twoMinuteOut := p.calculateFewMCandles(in2mChan, CandlePeriod2m)

	file2mChan := make(chan Candle, LenOfBuffer)
	in10mChan := make(chan Candle, LenOfBuffer)

	go func() {
		for candle := range twoMinuteOut {
			file2mChan <- candle
			in10mChan <- candle
		}

		close(file2mChan)
		close(in10mChan)
	}()

	wgFiles.Add(1)
	go p.writeCandlesToFile(file2mChan, CandleFile2m, wgFiles)
	tenMinuteOut := p.calculateFewMCandles(in10mChan, CandlePeriod10m)

	wgFiles.Add(1)
	go p.writeCandlesToFile(tenMinuteOut, CandleFile10m, wgFiles)

	wgFiles.Wait()
	p.logger.Info("End writing to files")
}

func (p Pipeline) writeCandlesToFile(in <-chan Candle, filename CandleFilename, wg *sync.WaitGroup) {
	defer wg.Done()

	file, err := os.Create(string(filename))
	if err != nil {
		p.logger.Error(err)
		return
	}
	defer file.Close()

	for candle := range in {
		p.logger.Debugf("write to file %s: %+v", filename, candle)

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
}

func (p Pipeline) calculate1mCandles(in <-chan Price) <-chan Candle {
	out := make(chan Candle)

	go func() {
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
					candle := Candle{
						Ticker: price.Ticker,
						Period: CandlePeriod1m,
						TS:     prevTS,
						Open:   current[0].Value,
						Close:  current[len(current)-1].Value,
						High:   current[0].Value,
						Low:    current[0].Value,
					}

					for _, p := range current {
						if p.Value > candle.High {
							candle.High = p.Value
						}
						if p.Value < candle.Low {
							candle.Low = p.Value
						}
					}

					out <- candle

					current = make([]Price, 0, 1)
					current = append(current, price)
					companies[price.Ticker] = current
				}
			}
		}

		p.logger.Info("end reading prices")
		close(out)
	}()

	return out
}

func (p Pipeline) calculateFewMCandles(in <-chan Candle, period CandlePeriod) <-chan Candle {
	out := make(chan Candle)

	go func() {
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
					candleToOut := Candle{
						Ticker: candle.Ticker,
						Period: period,
						TS:     prevTS,
						Open:   current[0].Open,
						Close:  current[len(current)-1].Close,
						High:   current[0].High,
						Low:    current[0].Low,
					}

					for _, c := range current {
						if c.High > candleToOut.High {
							candleToOut.High = c.High
						}
						if c.Low < candleToOut.Low {
							candleToOut.Low = c.Low
						}
					}

					out <- candleToOut

					current = make([]Candle, 0, 1)
					current = append(current, candle)
					companies[candle.Ticker] = current
				}
			}
		}

		p.logger.Infof("end reading candles for: %s", period)
		close(out)
	}()

	return out
}
