package main

import (
	"fmt"

	"github.com/MayukhSobo/scaffold/internal/server"
	"github.com/MayukhSobo/scaffold/pkg/config"
	"github.com/MayukhSobo/scaffold/pkg/log"
	"github.com/gofiber/fiber/v2"
)

func main() {
	// Display startup banner
	fmt.Println("Starting Scaffold API Server...")

	// Load configuration
	conf := config.NewConfig()

	// Initialize logger
	logger, err := log.CreateLoggerFromConfig(conf)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}

	logger.Info("Starting application...")

	// Start server with configuration
	server.RunServer(conf, logger)
}

// Example with custom setup
func mainWithCustomSetup() {
	// Load configuration
	conf := config.NewConfig()

	// Initialize logger
	logger, err := log.CreateLoggerFromConfig(conf)
	if err != nil {
		panic(fmt.Sprintf("failed to create logger: %v", err))
	}

	logger.Info("Starting application with custom setup...")

	// Start server with custom setup
	server.RunWithCustomSetup(conf, logger, func(s *server.FiberServer) {
		// Add custom routes
		s.AddRoutes(func(app *fiber.App) {
			app.Get("/custom", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{
					"message": "Custom route",
					"custom":  true,
				})
			})
		})

		// Add API v1 group
		s.AddGroup("/api/v1", func(router fiber.Router) {
			router.Get("/users", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{
					"users": []string{"user1", "user2"},
				})
			})

			router.Get("/posts", func(c *fiber.Ctx) error {
				return c.JSON(fiber.Map{
					"posts": []string{"post1", "post2"},
				})
			})
		})
	})
}
