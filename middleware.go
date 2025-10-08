package main

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Request ID middleware - adds correlation ID to each request
func requestIDMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Check if request ID already exists in headers
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		
		// Set request ID in context and response header
		c.Set("requestID", requestID)
		c.Header("X-Request-ID", requestID)
		
		// Add to logger context
		log.WithField("request_id", requestID).Debug("Processing request")
		
		c.Next()
	}
}

// Metrics middleware - tracks request metrics
func metricsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}
		
		c.Next()
		
		duration := time.Since(start)
		statusCode := fmt.Sprintf("%d", c.Writer.Status())
		
		// Record metrics
		requestsTotal.WithLabelValues(c.Request.Method, path, statusCode).Inc()
		requestDuration.WithLabelValues(c.Request.Method, path).Observe(duration.Seconds())
		
		// Log request completion
		requestID, _ := c.Get("requestID")
		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"method":     c.Request.Method,
			"path":       path,
			"status":     c.Writer.Status(),
			"duration":   duration.Milliseconds(),
			"ip":         c.ClientIP(),
		}).Info("Request completed")
	}
}

// User Portal authentication middleware - only accepts JWT tokens from User Portal
func userPortalAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID, _ := c.Get("requestID")
		
		// Check for Authorization header (JWT token from user login)
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.WithField("request_id", requestID).Warn("No JWT authorization header provided")
			sendError(c, http.StatusUnauthorized, "MISSING_TOKEN", "Authorization token required from User Portal")
			c.Abort()
			return
		}

		// Extract Bearer token
		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.WithField("request_id", requestID).Warn("Invalid JWT token format")
			sendError(c, http.StatusUnauthorized, "INVALID_TOKEN_FORMAT", "Token must be in format: Bearer <token>")
			c.Abort()
			return
		}

		token := tokenParts[1]

		// Validate JWT token
		if !isValidJWTToken(token) {
			log.WithField("request_id", requestID).Warn("JWT validation failed")
			sendError(c, http.StatusUnauthorized, "INVALID_TOKEN", "Invalid or expired token")
			c.Abort()
			return
		}

		// Store user info in context for use in handlers
		userID := extractUserIDFromToken(token)
		c.Set("authType", "user")
		c.Set("userID", userID)
		c.Set("token", token)
		
		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"auth_type":  "user",
			"user_id":    userID,
		}).Debug("User Portal authentication successful")

		c.Next()
	}
}