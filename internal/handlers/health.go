package handlers

import (
	"net/http"
	"time"

	"InternalAPI/internal/circuitbreaker"
	"InternalAPI/internal/models"

	"github.com/gin-gonic/gin"
)

// HealthHandler handles health check requests
func HealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "healthy",
		"service":   "internal-api",
		"timestamp": time.Now().Unix(),
	})
}

// GetCircuitBreakerStatusHandler returns the status of all circuit breakers
func GetCircuitBreakerStatusHandler(c *gin.Context) {
	status := circuitbreaker.GetAllStatus()

	c.JSON(http.StatusOK, gin.H{
		"circuit_breakers": status,
		"timestamp":        time.Now().Unix(),
	})
}

// ResetCircuitBreakerHandler resets a specific circuit breaker
func ResetCircuitBreakerHandler(c *gin.Context) {
	serviceName := c.Param("service")

	err := circuitbreaker.ResetByName(serviceName)
	if err != nil {
		sendError(c, http.StatusNotFound, "SERVICE_NOT_FOUND", "Circuit breaker for service not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Circuit breaker for " + serviceName + " has been reset",
	})
}

// sendError sends an error response
func sendError(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, models.ErrorResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}
