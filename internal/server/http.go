package server

import (
	"golang-di/internal/handler"
	"golang-di/internal/middleware"
	resp "golang-di/pkg/helper"
	"golang-di/pkg/log"

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
		resp.HandleSuccess(ctx, map[string]any{
			"say": "Hi Nunu!",
		})
	})
	r.GET("/user", userHandler.GetUserById)

	return r
}
