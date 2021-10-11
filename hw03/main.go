package main

import (
	"context"
	"flag"
	"os"
	"os/signal"
	"time"

	"github.com/willsem/tfs-go-hw/hw03/domain"
	"golang.org/x/sync/errgroup"

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
	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error {
		return pipeline.Start(prices, groupCtx)
	})

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	group.Go(func() error {
		for range c {
			logger.Info("Programm is closing")
			cancel()
			break
		}
		return nil
	})

	err := group.Wait()
	if err != nil {
		logger.Error(err)
	}
}
