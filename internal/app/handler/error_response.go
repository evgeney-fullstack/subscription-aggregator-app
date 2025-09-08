package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// errorResponse represents a standardized error response structure
// Used to maintain consistent error formatting across all API endpoints
type errorResponse struct {
	Message string `json:"message"`
}

// statusResponse represents a standardized success response structure
// Used for operations that don't return data but need confirmation
type statusResponse struct {
	Status string `json:"status"`
}

// newErrorResponse logs an error and returns a standardized error response
// This ensures consistent error handling and logging across all handlers
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Error(message)
	c.AbortWithStatusJSON(statusCode, errorResponse{Message: message})
}
