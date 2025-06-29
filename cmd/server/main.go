package main

import (
	"fmt"
	"golang-di/pkg/config"
	"golang-di/pkg/log"

	"github.com/spf13/viper"
)

var (
	conf   *viper.Viper
	logger log.Logger
)

func init() {
	// Display startup banner
	fmt.Println(DisplayBanner())
	conf = config.NewConfig()
	logger = config.CreateLoggerFromConfig(conf)
}

func main() {
	logger.Info("Starting application...")
}
