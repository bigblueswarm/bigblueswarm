// BigBlueSwarm is a metrics based load balancer for BigBlueButton service
package main

import (
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/app"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"

	log "github.com/sirupsen/logrus"
)

var version = ""
var buildTime = ""
var commitHash = ""

func main() {
	displayStartup()
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

func displayStartup() {
	fmt.Println("----------------------------------------------------------")
	fmt.Println("BigBlueSwarm load balancer")
	fmt.Println("Version:     \t" + version)
	fmt.Println("Build date @:\t" + buildTime)
	fmt.Println("Commit:      \t" + commitHash)
	fmt.Println("----------------------------------------------------------")
}

func initLog() {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetReportCaller(true)
}

func run(conf *config.Config) error {
	log.Info("Starting BigBlueSwarm server")
	err := app.NewServer(conf).Run()

	if err != nil {
		return err
	}

	return nil
}

// go build -ldflags="-X 'main.version=v1.0.0' -X 'main.buildTime=$(date)' -X 'main.commitHash=$(git rev-parse HEAD)'" -o main ./cmd/bigblueswarm/main.go
