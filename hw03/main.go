package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/willsem/tfs-go-hw/hw03/domain"

	log "github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/hw03/generator"
)

var (
	tickers = []string{"AAPL", "SBER", "NVDA", "TSLA"}
	debug   bool
)

func init() {
	flag.BoolVar(&debug, "debug", false, "logging level")
}

func main() {
	flag.Parse()

	logger := log.New()
	if debug {
		logger.Level = log.DebugLevel
	}

	ctx, cancel := context.WithCancel(context.Background())

	pg := generator.NewPricesGenerator(generator.Config{
		Factor:  10,
		Delay:   time.Millisecond * 500,
		Tickers: tickers,
	})

	logger.Info("start prices generator...")
	prices := pg.Prices(ctx)

	pipeline := domain.NewPipeline(logger)
	wg := &sync.WaitGroup{}

	wg.Add(1)
	go pipeline.Start(prices, wg)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		for _ = range c {
			logger.Info("Programm is closing")
			cancel()
			wg.Done()
		}
	}(wg)

	wg.Wait()
}
