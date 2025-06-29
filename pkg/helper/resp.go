package resp

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
