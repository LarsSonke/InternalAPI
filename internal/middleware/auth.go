package middleware

import (
	"net/http"
	"strings"
	"time"

	"InternalAPI/internal/models"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware validates authentication for protected routes
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sendError(c, http.StatusUnauthorized, "MISSING_AUTH", "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || parts[0] != "Bearer" {
			sendError(c, http.StatusUnauthorized, "INVALID_AUTH_FORMAT", "Authorization header must be in format 'Bearer <token>'")
			c.Abort()
			return
		}

		token := parts[1]

		// Here you would validate the token with your auth service
		// For now, we'll simulate user info extraction from JWT
		userInfo := &models.UserInfo{
			UserID:   "user123",
			Username: "testuser",
			Email:    "test@example.com",
			Roles:    []string{"user"},
			Exp:      time.Now().Add(time.Hour).Unix(),
		}

		// Store user info in context for use in handlers
		c.Set("user", userInfo)
		c.Set("userID", userInfo.UserID) // For backward compatibility
		c.Set("token", token)
		c.Next()
	}
}

// RequireRoles creates middleware that requires specific roles
func RequireRoles(requiredRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userInterface, exists := c.Get("user")
		if !exists {
			sendError(c, http.StatusUnauthorized, "MISSING_USER", "User information not found in context")
			c.Abort()
			return
		}

		user, ok := userInterface.(*models.UserInfo)
		if !ok {
			sendError(c, http.StatusInternalServerError, "INVALID_USER_TYPE", "Invalid user information type")
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, userRole := range user.Roles {
			for _, requiredRole := range requiredRoles {
				if userRole == requiredRole {
					hasRole = true
					break
				}
			}
			if hasRole {
				break
			}
		}

		if !hasRole {
			sendError(c, http.StatusForbidden, "INSUFFICIENT_PERMISSIONS", "User does not have required permissions")
			c.Abort()
			return
		}

		c.Next()
	}
}

// AdminOnly is a convenience middleware for admin-only routes
func AdminOnly() gin.HandlerFunc {
	return RequireRoles("admin", "super_admin")
}

// sendError sends an error response using the structured models
func sendError(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, models.ErrorResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}