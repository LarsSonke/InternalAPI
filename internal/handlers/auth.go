package handlers

import (
	"net/http"

	"InternalAPI/internal/config"
	"InternalAPI/internal/models"
	"InternalAPI/internal/services"

	"github.com/gin-gonic/gin"
)

// AuthHandlers contains all authentication-related handlers
type AuthHandlers struct {
	externalService *services.ExternalService
}

// NewAuthHandlers creates a new auth handlers instance
func NewAuthHandlers(config *config.Config) *AuthHandlers {
	return &AuthHandlers{
		externalService: services.New(config),
	}
}

// Login handles user login
func (ah *AuthHandlers) Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Call central management service for authentication
	authData := map[string]interface{}{
		"username": req.Username,
		"password": req.Password,
	}

	response, err := ah.externalService.Call("central", "POST", "/auth/login", authData)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "AUTH_SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
func (ah *AuthHandlers) RefreshToken(c *gin.Context) {
	var req models.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	// Call central management service for token refresh
	refreshData := map[string]interface{}{
		"refresh_token": req.RefreshToken,
	}

	response, err := ah.externalService.Call("central", "POST", "/auth/refresh", refreshData)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "AUTH_SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// Logout handles user logout
func (ah *AuthHandlers) Logout(c *gin.Context) {
	token, exists := c.Get("token")
	if !exists {
		sendError(c, http.StatusUnauthorized, "MISSING_TOKEN", "Token not found")
		return
	}

	// Call central management service for logout
	logoutData := map[string]interface{}{
		"token": token,
	}

	_, err := ah.externalService.Call("central", "POST", "/auth/logout", logoutData)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "AUTH_SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Successfully logged out",
	})
}

// GetUserInfo returns current user information
func (ah *AuthHandlers) GetUserInfo(c *gin.Context) {
	user, exists := c.Get("user")
	if !exists {
		sendError(c, http.StatusUnauthorized, "MISSING_USER", "User information not found")
		return
	}

	c.JSON(http.StatusOK, user)
}

// ChangePassword handles password change requests
func (ah *AuthHandlers) ChangePassword(c *gin.Context) {
	var req models.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	user, exists := c.Get("user")
	if !exists {
		sendError(c, http.StatusUnauthorized, "MISSING_USER", "User information not found")
		return
	}

	userInfo := user.(*models.UserInfo)

	// Call central management service for password change
	changeData := map[string]interface{}{
		"user_id":          userInfo.UserID,
		"current_password": req.CurrentPassword,
		"new_password":     req.NewPassword,
	}

	response, err := ah.externalService.Call("central", "PUT", "/auth/change-password", changeData)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "AUTH_SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
