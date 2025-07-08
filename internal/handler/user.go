package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/http"
	"github.com/MayukhSobo/scaffold/pkg/log"
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

	ctx := context.Background()
	adminUsers, err := h.userService.GetAdminUsers(ctx)
	if err != nil {
		h.GetLogger().Error("Failed to retrieve admin users", log.Error(err))
		return http.HandleFiberError(c, fiber.StatusInternalServerError, "Failed to retrieve admin users")
	}

	// Convert to response models (excludes password_hash)
	userResponses := ToUserResponses(adminUsers)

	h.GetLogger().Info("Retrieved admin users", log.Int("count", len(adminUsers)))
	return http.HandleFiberSuccess(c, fiber.Map{
		"users": userResponses,
		"count": len(userResponses),
	})
}

// GetPendingVerificationUsers retrieves all users with pending verification status
func (h *UserHandler) GetPendingVerificationUsers(c *fiber.Ctx) error {
	h.GetLogger().Info("GetPendingVerificationUsers called")

	ctx := context.Background()
	pendingUsers, err := h.userService.GetPendingVerificationUsers(ctx)
	if err != nil {
		h.GetLogger().Error("Failed to retrieve pending verification users", log.Error(err))
		return http.HandleFiberError(c, fiber.StatusInternalServerError, "Failed to retrieve pending verification users")
	}

	// Convert to response models (excludes password_hash)
	userResponses := ToUserResponses(pendingUsers)

	h.GetLogger().Info("Retrieved pending verification users", log.Int("count", len(pendingUsers)))
	return http.HandleFiberSuccess(c, fiber.Map{
		"users": userResponses,
		"count": len(userResponses),
	})
}
