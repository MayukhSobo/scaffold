package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gofiber/fiber/v2"
)

// Response represents the standard API response structure
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HandleSuccess sends a successful response
func HandleSuccess(ctx *gin.Context, data interface{}) {
	response := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	ctx.JSON(http.StatusOK, response)
}

// HandleError sends an error response
func HandleError(ctx *gin.Context, statusCode int, errorCode int, message string, data interface{}) {
	response := Response{
		Code:    errorCode,
		Message: message,
		Data:    data,
	}
	ctx.JSON(statusCode, response)
}

// HandleErrorWithStatus sends an error response with custom status and message
func HandleErrorWithStatus(ctx *gin.Context, statusCode int, message string) {
	HandleError(ctx, statusCode, statusCode, message, nil)
}

// HandleBadRequest sends a 400 Bad Request response
func HandleBadRequest(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusBadRequest, http.StatusBadRequest, message, nil)
}

// HandleInternalError sends a 500 Internal Server Error response
func HandleInternalError(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusInternalServerError, http.StatusInternalServerError, message, nil)
}

// HandleNotFound sends a 404 Not Found response
func HandleNotFound(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusNotFound, http.StatusNotFound, message, nil)
}

// HandleUnauthorized sends a 401 Unauthorized response
func HandleUnauthorized(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusUnauthorized, http.StatusUnauthorized, message, nil)
}

// HandleForbidden sends a 403 Forbidden response
func HandleForbidden(ctx *gin.Context, message string) {
	HandleError(ctx, http.StatusForbidden, http.StatusForbidden, message, nil)
}

// Fiber-specific response utilities

// HandleFiberSuccess sends a successful response for Fiber
func HandleFiberSuccess(c *fiber.Ctx, data interface{}) error {
	response := Response{
		Code:    0,
		Message: "success",
		Data:    data,
	}
	return c.Status(fiber.StatusOK).JSON(response)
}

// HandleFiberError sends an error response for Fiber
func HandleFiberError(c *fiber.Ctx, statusCode int, message string) error {
	response := Response{
		Code:    statusCode,
		Message: message,
		Data:    nil,
	}
	return c.Status(statusCode).JSON(response)
}

// HandleFiberBadRequest sends a 400 Bad Request response for Fiber
func HandleFiberBadRequest(c *fiber.Ctx, message string) error {
	return HandleFiberError(c, fiber.StatusBadRequest, message)
}

// HandleFiberInternalError sends a 500 Internal Server Error response for Fiber
func HandleFiberInternalError(c *fiber.Ctx, message string) error {
	return HandleFiberError(c, fiber.StatusInternalServerError, message)
}

// HandleFiberNotFound sends a 404 Not Found response for Fiber
func HandleFiberNotFound(c *fiber.Ctx, message string) error {
	return HandleFiberError(c, fiber.StatusNotFound, message)
}

// HandleFiberUnauthorized sends a 401 Unauthorized response for Fiber
func HandleFiberUnauthorized(c *fiber.Ctx, message string) error {
	return HandleFiberError(c, fiber.StatusUnauthorized, message)
}

// HandleFiberForbidden sends a 403 Forbidden response for Fiber
func HandleFiberForbidden(c *fiber.Ctx, message string) error {
	return HandleFiberError(c, fiber.StatusForbidden, message)
}
