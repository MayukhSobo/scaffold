package handler

import (
	"github.com/gofiber/fiber/v2"

	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"
	"github.com/MayukhSobo/scaffold/pkg/utils"
)

func NewUserHandler(handler *Handler, userService service.UserService) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

type UserHandler struct {
	*Handler
	userService service.UserService
}

// GetAdminUsers retrieves all users with admin access
func (h *UserHandler) GetAdminUsers(c *fiber.Ctx) error {
	h.GetLogger().Info("GetAdminUsers called")

	// TODO: Implement actual admin user retrieval logic
	// For now, return mock data to demonstrate the structure
	adminUsers := []map[string]interface{}{
		{
			"id":       1,
			"username": "admin",
			"role":     "admin",
			"status":   "active",
		},
		{
			"id":       2,
			"username": "superadmin",
			"role":     "super_admin",
			"status":   "active",
		},
	}

	h.GetLogger().Info("Retrieved admin users", log.Int("count", len(adminUsers)))
	return utils.HandleFiberSuccess(c, fiber.Map{
		"users": adminUsers,
		"count": len(adminUsers),
	})
}

// GetPendingVerificationUsers retrieves all users with pending verification status
func (h *UserHandler) GetPendingVerificationUsers(c *fiber.Ctx) error {
	h.GetLogger().Info("GetPendingVerificationUsers called")

	// TODO: Implement actual pending verification user retrieval logic
	// For now, return mock data to demonstrate the structure
	pendingUsers := []map[string]interface{}{
		{
			"id":                 3,
			"username":           "user1",
			"email":              "user1@example.com",
			"status":             "pending_verification",
			"created_at":         "2024-01-01T00:00:00Z",
			"verification_token": "abc123",
		},
		{
			"id":                 4,
			"username":           "user2",
			"email":              "user2@example.com",
			"status":             "pending_verification",
			"created_at":         "2024-01-02T00:00:00Z",
			"verification_token": "def456",
		},
	}

	h.GetLogger().Info("Retrieved pending verification users", log.Int("count", len(pendingUsers)))
	return utils.HandleFiberSuccess(c, fiber.Map{
		"users": pendingUsers,
		"count": len(pendingUsers),
	})
}
