package server

import (
	"github.com/MayukhSobo/scaffold/internal/handler"
	"github.com/MayukhSobo/scaffold/internal/middleware"
	"github.com/MayukhSobo/scaffold/pkg/log"
	"github.com/MayukhSobo/scaffold/pkg/utils"

	"github.com/gin-gonic/gin"
)

func NewServerHTTP(
	logger log.Logger,
	userHandler *handler.UserHandler,
) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.Use(
		middleware.CORSMiddleware(),
	)
	r.GET("/", func(ctx *gin.Context) {
		logger.Info("Root endpoint called")
		utils.HandleSuccess(ctx, map[string]any{
			"say": "Hi Nunu!",
		})
	})
	r.GET("/user", userHandler.GetUserById)

	return r
}
