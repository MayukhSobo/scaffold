package routes

import (
	"github.com/gofiber/fiber/v2"

	"github.com/MayukhSobo/scaffold/internal/handler"
	"github.com/MayukhSobo/scaffold/internal/service"
)

// RegisterUserRoutes sets up the user-related routes requested by the user
func RegisterUserRoutes(router fiber.Router, baseHandler *handler.Handler, userService service.UserService) {
	// Create user handler
	userHandler := handler.NewUserHandler(baseHandler, userService)

	// User routes group
	users := router.Group("/users")

	// Admin-specific user routes
	users.Get("/admin", userHandler.GetAdminUsers) // GET /api/v1/users/admin

	// Verification-specific user routes
	users.Get("/pending-verification", userHandler.GetPendingVerificationUsers) // GET /api/v1/users/pending-verification
}
