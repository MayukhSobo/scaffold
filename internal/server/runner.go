package server

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/pkg/log"
)

// RunServer starts the Fiber server with graceful shutdown
func RunServer(config *viper.Viper, logger log.Logger) {
	// Create the server
	server := NewFiberServer(config, logger)

	// Get the Fiber app
	app := server.GetApp()

	// Run the server
	RunFiberApp(app, config, logger)
}

// RunFiberApp runs a Fiber app with graceful shutdown
func RunFiberApp(app *fiber.App, config *viper.Viper, logger log.Logger) {
	// Get port from config
	port := config.GetString("http.port")
	if port == "" {
		port = "8000"
	}

	// Create a channel to listen for interrupt signals
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		logger.Infof("Server starting on port %s", port)
		if err := app.Listen(":" + port); err != nil {
			logger.Errorf("Server startup failed: %v", err)
			os.Exit(1)
		}
	}()

	// Wait for interrupt signal
	<-quit
	logger.Info("Shutting down server...")

	// Get shutdown timeout from config
	shutdownTimeout := config.GetDuration("server.shutdown_timeout")
	if shutdownTimeout == 0 {
		shutdownTimeout = 30 * time.Second
	}

	// Create a context with timeout for graceful shutdown
	ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()

	// Shutdown server
	if err := app.ShutdownWithContext(ctx); err != nil {
		logger.Errorf("Server forced to shutdown: %v", err)
		os.Exit(1)
	}

	logger.Info("Server exited")
}

// RunWithCustomSetup allows custom setup before starting the server
func RunWithCustomSetup(config *viper.Viper, logger log.Logger, setupFunc func(*FiberServer)) {
	// Create the server
	server := NewFiberServer(config, logger)

	// Apply custom setup
	if setupFunc != nil {
		setupFunc(server)
	}

	// Get the Fiber app
	app := server.GetApp()

	// Run the server
	RunFiberApp(app, config, logger)
}
