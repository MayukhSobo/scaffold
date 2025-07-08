package server

import (
	"fmt"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/internal/routes"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/container"
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

	// Custom logger middleware using our structured logger
	if s.config.GetBool("server.middleware.logger") {
		s.app.Use(s.createLoggerMiddleware())
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

// createLoggerMiddleware creates a custom logger middleware using our structured logger
func (s *FiberServer) createLoggerMiddleware() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		// Process request
		err := c.Next()

		// Calculate latency
		latency := time.Since(start)

		// Build fields dynamically, only including meaningful values
		fields := []log.Field{
			log.String("method", c.Method()),
			log.String("path", c.Path()),
			log.Int("status", c.Response().StatusCode()),
		}

		// Only add query if it exists
		if query := c.Request().URI().QueryArgs().String(); query != "" {
			fields = append(fields, log.String("query", query))
		}

		// Only add IP if it's not localhost
		if ip := c.IP(); ip != "127.0.0.1" && ip != "::1" {
			fields = append(fields, log.String("ip", ip))
		}

		// Only add user agent if it's not a common development tool
		if userAgent := c.Get("User-Agent"); userAgent != "" &&
			!strings.Contains(strings.ToLower(userAgent), "insomnia") &&
			!strings.Contains(strings.ToLower(userAgent), "postman") &&
			!strings.Contains(strings.ToLower(userAgent), "curl") {
			fields = append(fields, log.String("user_agent", userAgent))
		}

		// Human-readable latency
		fields = append(fields, log.String("latency", s.formatLatency(latency)))

		// Human-readable bytes sent
		fields = append(fields, log.String("bytes_sent", s.formatBytes(len(c.Response().Body()))))

		// Add request ID if available
		if requestID := c.Get("X-Request-ID"); requestID != "" {
			fields = append(fields, log.String("request_id", requestID))
		} else if rid := c.Locals("requestid"); rid != nil {
			fields = append(fields, log.String("request_id", rid.(string)))
		}

		// Log based on status code
		status := c.Response().StatusCode()
		switch {
		case status >= 500:
			s.logger.Error("HTTP Request", fields...)
		case status >= 400:
			s.logger.Warn("HTTP Request", fields...)
		default:
			s.logger.Info("HTTP Request", fields...)
		}

		return err
	}
}

// formatLatency formats duration in a human-readable way
func (s *FiberServer) formatLatency(d time.Duration) string {
	if d < time.Microsecond {
		return d.String()
	}
	if d < time.Millisecond {
		return fmt.Sprintf("%.0fµs", float64(d.Nanoseconds())/1000)
	}
	if d < time.Second {
		return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1000000)
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

// formatBytes formats byte count in a human-readable way
func (s *FiberServer) formatBytes(bytes int) string {
	if bytes < 1024 {
		return fmt.Sprintf("%dB", bytes)
	}
	if bytes < 1024*1024 {
		return fmt.Sprintf("%.1fKB", float64(bytes)/1024)
	}
	if bytes < 1024*1024*1024 {
		return fmt.Sprintf("%.1fMB", float64(bytes)/(1024*1024))
	}
	return fmt.Sprintf("%.1fGB", float64(bytes)/(1024*1024*1024))
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

// SetupBusinessRoutesWithContainer configures business logic routes using the container pattern
// This is the new, scalable approach that handles multiple services and repositories
func (s *FiberServer) SetupBusinessRoutesWithContainer(container *container.TypedContainer) {
	// Create route config using container
	routeConfig := &routes.ContainerRouteConfig{
		App:       s.app,
		Container: container,
	}

	// Register business routes using container pattern
	routes.RegisterRoutesWithContainer(routeConfig)
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
