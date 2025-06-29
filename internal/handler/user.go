package handler

import (
	"net/http"

	"github.com/MayukhSobo/scaffold/internal/service"
	"github.com/MayukhSobo/scaffold/pkg/log"
	"github.com/MayukhSobo/scaffold/pkg/utils"

	"github.com/gin-gonic/gin"
)

func NewUserHandler(handler *Handler,
	userService service.UserService,
) *UserHandler {
	return &UserHandler{
		Handler:     handler,
		userService: userService,
	}
}

type UserHandler struct {
	*Handler
	userService service.UserService
}

func (h *UserHandler) GetUserById(ctx *gin.Context) {
	var params struct {
		Id int64 `form:"id" binding:"required"`
	}
	if err := ctx.ShouldBind(&params); err != nil {
		utils.HandleError(ctx, http.StatusBadRequest, 1, err.Error(), nil)
		return
	}

	user, err := h.userService.GetUserById(params.Id)

	// Use clean logger interface without framework-specific methods
	h.GetLogger().Info("GetUserByID", log.Any("user", user))

	if err != nil {
		h.GetLogger().Error("Failed to get user", log.Error(err), log.Int64("user_id", params.Id))
		utils.HandleError(ctx, http.StatusInternalServerError, 1, err.Error(), nil)
		return
	}
	utils.HandleSuccess(ctx, user)
}

func (h *UserHandler) UpdateUser(ctx *gin.Context) {
	h.GetLogger().Info("UpdateUser called")
	utils.HandleSuccess(ctx, nil)
}
