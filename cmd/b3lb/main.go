package main

import (
	"b3lb/pkg/app"
	"fmt"
	"os"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "this is the startup error: %s\n", err)
	}
}

func run() error {
	err := app.NewServer().Run()

	if err != nil {
		return err
	}

	return nil
}
