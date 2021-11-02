package main

import (
	"b3lb/pkg/app"
	"b3lb/pkg/config"
	"flag"
	"fmt"
)

func main() {
	conf, err := config.Load(configPath())

	if err != nil {
		panic(fmt.Sprintf("Unable to load configuration: %s \n", err))
	}

	if err := run(*conf); err != nil {
		panic(fmt.Sprintf("Server can't start: %s\n", err))
	}
}

func configPath() string {
	var configPath string

	flag.StringVar(&configPath, "config", "config.yml", "Config file path")
	flag.Parse()

	return configPath
}

func run(conf config.Config) error {
	err := app.NewServer(&conf).Run()

	if err != nil {
		return err
	}

	return nil
}
