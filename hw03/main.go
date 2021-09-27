package main

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw03/generator"
)

var tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}

func main() {
	logger := log.New()
	ctx, cancel := context.WithCancel(context.Background())

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	for i := 0; i <= 10; i++ {
		logger.Infof("prices %d: %+v", i, <-prices)
	}

	cancel()
}
