package main

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/course_project/internal/config"
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
}
