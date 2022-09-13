package main

import (
	"fmt"

	"github.com/SLedunois/b3lb/v2/pkg/app"
	"github.com/SLedunois/b3lb/v2/pkg/config"

	log "github.com/sirupsen/logrus"
)

const version = "2.1.0"

func main() {
	initLog()
	configPath, err := config.FormalizeConfigPath(config.Path())
	if err != nil {
		panic(fmt.Errorf("unable to parse configuration: %s", err.Error()))
	}

	conf, err := config.Load(configPath)

	if err != nil {
		panic(fmt.Sprintf("Unable to load configuration: %s \n", err))
	}

	if err := run(conf); err != nil {
		panic(fmt.Sprintf("Server can't start: %s\n", err))
	}
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
}

func run(conf *config.Config) error {
	log.Info(fmt.Sprintf("Starting b3lb server version %s", version))
	err := app.NewServer(conf).Run()

	if err != nil {
		return err
	}

	return nil
}
