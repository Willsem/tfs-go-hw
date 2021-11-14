package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/course_project/internal/config"
	"github.com/willsem/tfs-go-hw/course_project/internal/domain"
	"github.com/willsem/tfs-go-hw/course_project/internal/repositories/applications"
	"github.com/willsem/tfs-go-hw/course_project/internal/services/telegram"
	postgres "github.com/willsem/tfs-go-hw/course_project/pkg/postres"
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

	var config config.Toml
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		logger.Fatal(err)
	}

	pool, err := postgres.NewPool(config.Database.ConnectionString)
	if err != nil {
		logger.Fatal(err)
	}
	defer pool.Close()

	repository := applications.New(pool)

	app := domain.Application{
		Ticker: "APPL",
		Cost:   150,
	}
	err = repository.Add(app)
	if err != nil {
		logger.Fatal(err)
	}

	app.Cost++
	err = repository.Add(app)
	if err != nil {
		logger.Fatal(err)
	}

	app.Ticker = "AMD"
	err = repository.Add(app)
	if err != nil {
		logger.Fatal(err)
	}

	fmt.Println(repository.GetAll())
	fmt.Println(repository.GetByTicker("APPL"))
	fmt.Println(repository.GetByTicker("AMD"))

	bot, err := telegram.NewBot(repository, logger, config.Telegram)
	if err != nil {
		logger.Fatal(err)
	}

	bot.Start()

	bufio.NewReader(os.Stdin).ReadBytes('\n')
}
