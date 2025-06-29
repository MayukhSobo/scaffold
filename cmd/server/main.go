package main

import (
	"fmt"

	"github.com/MayukhSobo/scaffold/pkg/config"
	"github.com/MayukhSobo/scaffold/pkg/log"

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
	var err error
	logger, err = log.CreateLoggerFromConfig(conf)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}
}

func main() {
	logger.Info("Starting application...")
}
