package server

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/internal/routes"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"
)

// FiberServer wraps the Fiber app with configuration
type FiberServer struct {
	app    *fiber.App
	config *viper.Viper
	logger log.Logger
}

// NewFiberServer creates a new Fiber server with the given configuration
func NewFiberServer(config *viper.Viper, logger log.Logger) *FiberServer {
	// Create Fiber app with config
	app := fiber.New(fiber.Config{
		AppName:      config.GetString("app.name"),
		ServerHeader: config.GetString("app.name") + " " + config.GetString("app.version"),
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Log the error
			logger.Error("Server error", log.Error(err), log.String("path", c.Path()))

			// Handle Fiber errors
			if e, ok := err.(*fiber.Error); ok {
				return c.Status(e.Code).JSON(fiber.Map{
					"error":   true,
					"message": e.Message,
					"code":    e.Code,
				})
			}

			// Handle generic errors
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error":   true,
				"message": "Internal server error",
				"code":    fiber.StatusInternalServerError,
			})
		},
	})

	server := &FiberServer{
		app:    app,
		config: config,
		logger: logger,
	}

	// Setup middleware
	server.setupMiddleware()

	// Setup routes
	server.setupRoutes()

	return server
}

// setupMiddleware configures all middleware
func (s *FiberServer) setupMiddleware() {
	// Recovery middleware
	if s.config.GetBool("server.middleware.recover") {
		s.app.Use(recover.New())
	}

	// Request ID middleware
	if s.config.GetBool("server.middleware.request_id") {
		s.app.Use(requestid.New())
	}

	// Logger middleware
	if s.config.GetBool("server.middleware.logger") {
		s.app.Use(logger.New(logger.Config{
			Format: s.config.GetString("server.middleware.logger_format"),
		}))
	}

	// CORS middleware
	if s.config.GetBool("server.middleware.cors") {
		s.app.Use(cors.New(cors.Config{
			AllowOrigins:     s.config.GetString("server.cors.allow_origins"),
			AllowMethods:     s.config.GetString("server.cors.allow_methods"),
			AllowHeaders:     s.config.GetString("server.cors.allow_headers"),
			AllowCredentials: s.config.GetBool("server.cors.allow_credentials"),
			MaxAge:           s.config.GetInt("server.cors.max_age"),
		}))
	}
}

// setupRoutes configures basic routes
func (s *FiberServer) setupRoutes() {
	// Health check endpoint
	s.app.Get("/health", func(c *fiber.Ctx) error {
		s.logger.Info("Health endpoint called")
		return c.JSON(fiber.Map{
			"status": "healthy",
			"env":    s.config.GetString("env"),
		})
	})

	// Ping endpoint
	s.app.Get("/ping", func(c *fiber.Ctx) error {
		s.logger.Info("Ping endpoint called")
		return c.JSON(fiber.Map{
			"message": "pong",
			"status":  "ok",
		})
	})

	// Root endpoint
	s.app.Get("/", func(c *fiber.Ctx) error {
		s.logger.Info("Root endpoint called")
		return c.JSON(fiber.Map{
			"message": "Welcome to " + s.config.GetString("app.name"),
			"version": s.config.GetString("app.version"),
			"status":  "running",
		})
	})
}

// SetupBusinessRoutes configures business logic routes with dependencies
func (s *FiberServer) SetupBusinessRoutes(userService service.UserService) {
	// Create route config
	routeConfig := &routes.RouteConfig{
		App:         s.app,
		Config:      s.config,
		Logger:      s.logger,
		UserService: userService,
	}

	// Register business routes
	routes.RegisterRoutes(routeConfig)
}

// GetApp returns the underlying Fiber app
func (s *FiberServer) GetApp() *fiber.App {
	return s.app
}

// AddRoutes allows adding additional routes to the server
func (s *FiberServer) AddRoutes(setupFunc func(*fiber.App)) {
	setupFunc(s.app)
}

// AddMiddleware allows adding additional middleware
func (s *FiberServer) AddMiddleware(middleware ...fiber.Handler) {
	for _, m := range middleware {
		s.app.Use(m)
	}
}

// AddGroup creates a new route group
func (s *FiberServer) AddGroup(prefix string, setupFunc func(fiber.Router)) {
	group := s.app.Group(prefix)
	setupFunc(group)
}
