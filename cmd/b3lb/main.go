package main

import (
	"flag"
	"fmt"

	"github.com/SLedunois/b3lb/pkg/app"
	"github.com/SLedunois/b3lb/pkg/config"

	log "github.com/sirupsen/logrus"
)

const version = "1.5.0"

func main() {
	initLog()
	configPath, err := config.FormalizeConfigPath(configPath())
	if err != nil {
		panic(fmt.Errorf("unable to parse configuration: %s", err.Error()))
	}

	conf, err := config.Load(configPath)

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

	flag.StringVar(&configPath, "config", config.DefaultConfigPath(), "Config file path")
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
