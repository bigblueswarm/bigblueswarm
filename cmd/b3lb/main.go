package main

import (
	"b3lb/pkg/app"
	"b3lb/pkg/config"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

const version = "1.0.0-SNAPSHOT"

func main() {
	initLog()
	conf, err := config.Load(configPath())

	if err != nil {
		panic(fmt.Sprintf("Unable to load configuration: %s \n", err))
	}

	if err := run(*conf); err != nil {
		panic(fmt.Sprintf("Server can't start: %s\n", err))
	}
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
}

func configPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yml", "Config file path")
	flag.Parse()

	return configPath
}

func run(conf config.Config) error {
	log.Info(fmt.Sprintf("Starting b3lb server version %s", version))
	err := app.NewServer(&conf).Run()

	if err != nil {
		return err
	}

	return nil
}
