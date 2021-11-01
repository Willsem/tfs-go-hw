package main

import (
	"flag"

	"github.com/BurntSushi/toml"
	log "github.com/sirupsen/logrus"
	"github.com/willsem/tfs-go-hw/course_project/internal/config"
)

var configPath string

func init() {
	flag.StringVar(
		&configPath,
		"config",
		"configs/config.toml",
		"path to config file",
	)
}

func main() {
	flag.Parse()
	log.Info("Path of configuration file: ", configPath)

	var config config.Toml
	if _, err := toml.DecodeFile(configPath, &config); err != nil {
		log.Fatal(err)
	}
}
