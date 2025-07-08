package main

import (
	"fmt"

	"github.com/MayukhSobo/scaffold/internal/repository"
	"github.com/MayukhSobo/scaffold/internal/server"
	"github.com/MayukhSobo/scaffold/internal/service"
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

	// Create dependencies
	logger.Info("Initializing dependencies...")

	// Create database connection (currently mock, but this is where you'd connect to real DB)
	db := repository.NewDb()
	logger.Info("Database initialized")

	// Create repository layer
	repo := repository.NewRepository(logger, db)
	userRepo := repository.NewUserRepository(repo)
	logger.Info("Repository layer initialized")

	// Create service layer
	baseService := service.NewService(logger)
	userService := service.NewUserService(baseService, userRepo)
	logger.Info("Service layer initialized")

	// Start server with custom setup to connect business routes
	logger.Info("Starting server with business routes...")
	server.RunWithCustomSetup(conf, logger, func(s *server.FiberServer) {
		// Setup business routes with dependencies
		s.SetupBusinessRoutes(userService)
		logger.Info("Business routes registered successfully")
	})
}
