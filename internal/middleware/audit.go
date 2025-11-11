package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

var auditLog *logrus.Logger

func init() {
	auditLog = logrus.New()
	auditLog.SetFormatter(&logrus.JSONFormatter{})
	auditLog.SetLevel(logrus.InfoLevel)
}

// responseWriter wraps gin.ResponseWriter to capture response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// AuditLogger logs all requests and responses for security audit trail
func AuditLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		
		// Capture request body (for non-GET requests)
		var requestBody []byte
		if c.Request.Method != "GET" && c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// Restore the body for the next handler
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Wrap response writer to capture response
		blw := &responseWriter{
			ResponseWriter: c.Writer,
			body:          bytes.NewBufferString(""),
		}
		c.Writer = blw

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user info if authenticated
		userID := "anonymous"
		if uid, exists := c.Get("userID"); exists {
			userID = uid.(string)
		}

		// Get request ID
		requestID := ""
		if rid, exists := c.Get("request_id"); exists {
			requestID = rid.(string)
		}

		// Log the request/response
		fields := logrus.Fields{
			"request_id":   requestID,
			"timestamp":    start.Unix(),
			"method":       c.Request.Method,
			"path":         c.Request.URL.Path,
			"query":        c.Request.URL.RawQuery,
			"ip":           c.ClientIP(),
			"user_agent":   c.Request.UserAgent(),
			"user_id":      userID,
			"status":       c.Writer.Status(),
			"duration_ms":  duration.Milliseconds(),
			"request_size": c.Request.ContentLength,
			"response_size": blw.body.Len(),
		}

		// Log request body for sensitive operations (excluding passwords)
		if c.Request.Method != "GET" && len(requestBody) > 0 && len(requestBody) < 1024 {
			// Don't log passwords or sensitive data
			if c.Request.URL.Path != "/auth/login" && 
			   c.Request.URL.Path != "/auth/change-password" &&
			   c.Request.URL.Path != "/admin/users" {
				fields["request_body"] = string(requestBody)
			}
		}

		// Log at different levels based on status
		if c.Writer.Status() >= 500 {
			auditLog.WithFields(fields).Error("Server error")
		} else if c.Writer.Status() >= 400 {
			auditLog.WithFields(fields).Warn("Client error")
		} else {
			auditLog.WithFields(fields).Info("Request completed")
		}
	}
}
