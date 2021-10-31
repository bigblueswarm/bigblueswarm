package main

import (
	"b3lb/pkg/app"
	"b3lb/pkg/config"
	"fmt"
)

func main() {
	conf, err := config.Init()

	if err != nil {
		panic(fmt.Sprintf("Unable to load configuration: %s \n", err))
	}

	if err := run(*conf); err != nil {
		panic(fmt.Sprintf("Server can't start: %s\n", err))
	}
}

func run(conf config.Config) error {
	err := app.NewServer(&conf).Run()

	if err != nil {
		return err
	}

	return nil
}
