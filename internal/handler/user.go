package handler

import (
	"context"

	"github.com/gofiber/fiber/v2"

	"github.com/MayukhSobo/scaffold/internal/repository/users"
	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"
	"github.com/MayukhSobo/scaffold/pkg/utils"
)

func NewUserHandler(handler *Handler, userService service.UserService, userRepo users.Querier) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
		userRepo:    userRepo,
	}
}

type UserHandler struct {
	*Handler
	userService service.UserService
	userRepo    users.Querier
}

// GetAdminUsers retrieves all users with admin access
func (h *UserHandler) GetAdminUsers(c *fiber.Ctx) error {
	h.GetLogger().Info("GetAdminUsers called")

	ctx := context.Background()
	adminUsers, err := h.userRepo.GetAdminUsers(ctx)
	if err != nil {
		h.GetLogger().Error("Failed to retrieve admin users", log.Error(err))
		return utils.HandleFiberError(c, fiber.StatusInternalServerError, "Failed to retrieve admin users")
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

	ctx := context.Background()
	pendingUsers, err := h.userRepo.GetPendingVerificationUsers(ctx)
	if err != nil {
		h.GetLogger().Error("Failed to retrieve pending verification users", log.Error(err))
		return utils.HandleFiberError(c, fiber.StatusInternalServerError, "Failed to retrieve pending verification users")
	}

	h.GetLogger().Info("Retrieved pending verification users", log.Int("count", len(pendingUsers)))
	return utils.HandleFiberSuccess(c, fiber.Map{
		"users": pendingUsers,
		"count": len(pendingUsers),
	})
}
