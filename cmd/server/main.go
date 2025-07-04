package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/MayukhSobo/scaffold/pkg/config"
	"github.com/MayukhSobo/scaffold/pkg/log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
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

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName: "Scaffold v1.0.0",
	})

	// Add middleware
	app.Use(recover.New())
	app.Use(cors.New())

	// Routes
	app.Get("/ping", func(c *fiber.Ctx) error {
		logger.Info("Ping endpoint called")
		return c.JSON(fiber.Map{
			"message": "pong",
			"status":  "ok",
		})
	})

	app.Get("/health", func(c *fiber.Ctx) error {
		logger.Info("Health endpoint called")
		return c.JSON(fiber.Map{
			"status": "healthy",
			"env":    conf.GetString("env"),
		})
	})

	// Get port from config
	port := conf.GetString("http.port")
	if port == "" {
		port = "8000"
	}

	// Start server in a goroutine
	go func() {
		logger.Info(fmt.Sprintf("Server starting on port %s", port))
		if err := app.Listen(":" + port); err != nil {
			logger.Error("Server startup failed", log.Error(err))
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Gracefully shutdown the server
	if err := app.Shutdown(); err != nil {
		logger.Error("Server forced to shutdown", log.Error(err))
		os.Exit(1)
	}

	logger.Info("Server exited")
}
