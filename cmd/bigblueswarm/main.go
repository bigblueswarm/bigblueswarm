// BigBlueSwarm is a metrics based load balancer for BigBlueButton service
package main

import (
	"flag"
	"fmt"
	"os"

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
	logPath    = ""
)

func main() {
	parseFlags()
	f, err := initLog()
	if err != nil {
		panic(err)
	}

	defer f.Close()

	displayStartup()

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

func initLog() (*os.File, error) {
	var file *os.File
	var err error
	disableColors := false

	if logPath != "" {
		disableColors = true
		file, err = os.OpenFile(logPath, os.O_WRONLY|os.O_CREATE, 0755)
		if err != nil {
			return nil, fmt.Errorf("unable to create or open %s log file: %s", logPath, err)
		}

		log.SetOutput(file)
		gin.DefaultWriter = file
		os.Stdout = file
		log.WithField("path", logPath).Infoln("writing logs in configured path")
	}

	if lvl, err := log.ParseLevel(logLevel); err == nil {
		log.Infoln("Setting up BigBlueSearm log level as", lvl.String())
		log.SetLevel(lvl)
	}

	log.SetReportCaller(true)

	log.SetFormatter(&log.TextFormatter{
		DisableColors: disableColors,
		FullTimestamp: true,
	})

	return file, err
}

func parseFlags() {
	flag.StringVar(&configPath, "config", config.DefaultConfigPath(), "Config file path")
	flag.StringVar(&logLevel, "log.level", log.InfoLevel.String(), "Log level. Default is debug for development")
	flag.StringVar(&logPath, "log.path", "", "Log path. Specify a path to write into a file. By default BigBlueSwarm prints log in stdout")
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
