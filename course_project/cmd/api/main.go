package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"

	"github.com/willsem/tfs-go-hw/course_project/internal/config"
	"github.com/willsem/tfs-go-hw/course_project/internal/handlers"
	"github.com/willsem/tfs-go-hw/course_project/internal/repositories"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/indicator"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/subscribe"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/telegram"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/trading"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/tradingbot"
	"github.com/willsem/tfs-go-hw/course_project/pkg/postgres"
)

var (
	configPath string
	debug      bool
)

func init() {
	flag.StringVar(&configPath, "config", "configs/config.toml", "path to config file")
	flag.BoolVar(&debug, "debug", false, "debug mode of application")
}

func main() {
	flag.Parse()

	logger := log.New()
	if debug {
		logger.SetLevel(log.DebugLevel)
	} else {
		logFile, err := os.OpenFile("./logs/traider.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		if err != nil {
			logger.Fatal(err)
		}
		defer logFile.Close()

		logger.SetFormatter(&log.JSONFormatter{})
		logger.SetOutput(logFile)
		logger.SetLevel(log.InfoLevel)
	}

	logger.Info("Path of configuration file: ", configPath)

	var parsedConfig config.Toml
	if _, err := toml.DecodeFile(configPath, &parsedConfig); err != nil {
		logger.Fatal(err)
	}

	pgsqlPool, err := postgres.NewPool(parsedConfig.Database.ConnectionString)
	if err != nil {
		logger.Fatal(err)
	}

	appRepo := repositories.NewApplicaitionsRepository(pgsqlPool)

	telegramBot, err := telegram.NewBot(appRepo, logger, parsedConfig.Telegram)
	if err != nil {
		logger.Fatal(err)
	}

	telegramBot.Start()

	subscribeService, err := subscribe.NewKrakenSubscribeService(parsedConfig.Kraken, logger)
	if err != nil {
		logger.Fatal(err)
	}

	tradingService := trading.NewKrakenTradingService(parsedConfig.Kraken)

	// indicatorService := indicator.NewTripleCandlesTemplate()
	indicatorService := indicator.NewOneCandleTemplate()

	tradingBot := tradingbot.New(subscribeService, tradingService, indicatorService, appRepo, telegramBot, logger)
	err = tradingBot.Start()
	if err != nil {
		logger.Fatal(err)
	}
	defer tradingBot.Stop()

	r := chi.NewRouter()

	tradingBotHandler := handlers.NewTradingBotHandler(tradingBot, logger)
	r.Mount("/trading", tradingBotHandler.Routes())

	logger.Info("listen " + parsedConfig.Server.ListenAddr)

	if debug {
		err = http.ListenAndServe(parsedConfig.Server.ListenAddr, r)
	} else {
		conf := parsedConfig.Server
		err = http.ListenAndServeTLS(conf.ListenAddr, conf.CertFile, conf.KeyFile, r)
	}

	if err != nil {
		logger.Fatal(err)
	}
}
