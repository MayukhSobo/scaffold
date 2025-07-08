package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/viper"

	"github.com/MayukhSobo/scaffold/internal/handler"
	"github.com/MayukhSobo/scaffold/internal/repository/users"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"
)

// RouteConfig holds the dependencies needed for route registration
type RouteConfig struct {
	App         *fiber.App
	Config      *viper.Viper
	Logger      log.Logger
	UserService service.UserService
	UserRepo    users.Querier
}

// RegisterRoutes sets up all application routes
func RegisterRoutes(rc *RouteConfig) {
	// Create base handler
	baseHandler := handler.NewHandler(rc.Logger)

	// Register API routes group
	api := rc.App.Group("/api")

	// Register v1 routes
	v1 := api.Group("/v1")

	// Register user routes
	RegisterUserRoutes(v1, baseHandler, rc.UserService, rc.UserRepo)
}
