package main

import (
	"fmt"

	"github.com/MayukhSobo/scaffold/internal/server"
	"github.com/MayukhSobo/scaffold/pkg/config"
	"github.com/MayukhSobo/scaffold/pkg/container"
	"github.com/MayukhSobo/scaffold/pkg/db"
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
	logger.Info("Starting application with container pattern...")

	// Create dependencies
	logger.Info("Initializing dependencies...")

	// Create database connection using the db package
	database := db.MustConnect(conf, logger)

	// Create dependency container - this handles ALL dependencies
	// When you add new services/repositories, just add them to the container
	appContainer := container.NewTypedContainer(conf, logger, database)
	logger.Info("Dependency container initialized with all services and repositories")

	// Start server with container-based setup
	logger.Info("Starting server with container-based routes...")
	server.RunWithCustomSetup(conf, logger, func(s *server.FiberServer) {
		// Setup business routes using container - scales to any number of services
		s.SetupBusinessRoutesWithContainer(appContainer)
		logger.Info("All business routes registered successfully via container")
	})
}
