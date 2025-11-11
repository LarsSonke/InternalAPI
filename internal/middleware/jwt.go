package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"InternalAPI/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gin-gonic/gin"
)

var (
	// Token blacklist for revoked tokens
	tokenBlacklist = make(map[string]time.Time)
	blacklistMu    sync.RWMutex
	
	// JWT secret key (should come from config)
	jwtSecretKey []byte
)

// InitJWT initializes the JWT secret key
func InitJWT(secret string) {
	jwtSecretKey = []byte(secret)
	
	// Start cleanup routine for expired blacklisted tokens
	go cleanupBlacklist()
}

// Claims represents JWT claims
type Claims struct {
	UserID   string   `json:"user_id"`
	Username string   `json:"username"`
	Email    string   `json:"email"`
	Roles    []string `json:"roles"`
	jwt.RegisteredClaims
}

// ValidateJWT validates a JWT token and returns the claims
func ValidateJWT(tokenString string) (*Claims, error) {
	if len(jwtSecretKey) == 0 {
		return nil, errors.New("JWT secret not initialized")
	}

	// Check if token is blacklisted
	if isBlacklisted(tokenString) {
		return nil, errors.New("token has been revoked")
	}

	// Parse and validate token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return jwtSecretKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Extract and validate claims
	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, errors.New("invalid token claims")
	}

	// Check expiration
	if claims.ExpiresAt != nil && claims.ExpiresAt.Time.Before(time.Now()) {
		return nil, errors.New("token has expired")
	}

	return claims, nil
}

// BlacklistToken adds a token to the blacklist
func BlacklistToken(tokenString string, expiresAt time.Time) {
	blacklistMu.Lock()
	defer blacklistMu.Unlock()
	tokenBlacklist[tokenString] = expiresAt
}

// isBlacklisted checks if a token is in the blacklist
func isBlacklisted(tokenString string) bool {
	blacklistMu.RLock()
	defer blacklistMu.RUnlock()
	_, exists := tokenBlacklist[tokenString]
	return exists
}

// cleanupBlacklist removes expired tokens from blacklist
func cleanupBlacklist() {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	for range ticker.C {
		blacklistMu.Lock()
		now := time.Now()
		for token, expiresAt := range tokenBlacklist {
			if expiresAt.Before(now) {
				delete(tokenBlacklist, token)
			}
		}
		blacklistMu.Unlock()
	}
}

// JWTAuthMiddleware validates JWT authentication for protected routes
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			sendError(c, http.StatusUnauthorized, "MISSING_AUTH", "Authorization header is required")
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>" format
		tokenString := extractToken(authHeader)
		if tokenString == "" {
			sendError(c, http.StatusUnauthorized, "INVALID_AUTH_FORMAT", "Authorization header must be in format 'Bearer <token>'")
			c.Abort()
			return
		}

		// Validate token
		claims, err := ValidateJWT(tokenString)
		if err != nil {
			sendError(c, http.StatusUnauthorized, "INVALID_TOKEN", fmt.Sprintf("Token validation failed: %v", err))
			c.Abort()
			return
		}

		// Store user info in context
		userInfo := &models.UserInfo{
			UserID:   claims.UserID,
			Username: claims.Username,
			Email:    claims.Email,
			Roles:    claims.Roles,
			Exp:      claims.ExpiresAt.Unix(),
		}

		c.Set("user", userInfo)
		c.Set("userID", userInfo.UserID)
		c.Set("token", tokenString)
		c.Next()
	}
}

// extractToken extracts the token from Authorization header
func extractToken(authHeader string) string {
	const bearerPrefix = "Bearer "
	if len(authHeader) > len(bearerPrefix) && authHeader[:len(bearerPrefix)] == bearerPrefix {
		return authHeader[len(bearerPrefix):]
	}
	return ""
}
