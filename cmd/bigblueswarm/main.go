// BigBlueSwarm is a metrics based load balancer for BigBlueButton service
package main

import (
	"flag"
	"fmt"

	"github.com/bigblueswarm/bigblueswarm/v2/pkg/app"
	"github.com/bigblueswarm/bigblueswarm/v2/pkg/config"
	"github.com/gin-gonic/gin"

	log "github.com/sirupsen/logrus"
)

var (
	version    = ""
	buildTime  = ""
	commitHash = ""
)

var (
	configPath = ""
	logLevel   = ""
)

func main() {
	displayStartup()
	parseFlags()
	initLog()
	configPath, err := config.FormalizeConfigPath(configPath)
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

	if lvl, err := log.ParseLevel(logLevel); err == nil {
		log.Infoln("Setting up BigBlueSearm log level as", lvl.String())
		log.SetLevel(lvl)
	}

	log.SetReportCaller(true)
}

func parseFlags() {
	flag.StringVar(&configPath, "config", config.DefaultConfigPath(), "Config file path")
	flag.StringVar(&logLevel, "log.level", log.DebugLevel.String(), "Log level. Default is debug for development")
	flag.Parse()
}

func isDebugMode() bool {
	return log.GetLevel() == log.DebugLevel
}

func run(conf *config.Config) error {
	if !isDebugMode() {
		gin.SetMode(gin.ReleaseMode)
	}

	err := app.NewServer(conf).Run()

	if err != nil {
		return err
	}

	return nil
}
