package handlers

import (
	"net/http"

	"InternalAPI/internal/config"
	"InternalAPI/internal/models"
	"InternalAPI/internal/services"

	"github.com/gin-gonic/gin"
)

// AdminHandlers contains all admin-related handlers
type AdminHandlers struct {
	externalService *services.ExternalService
}

// NewAdminHandlers creates a new admin handlers instance
func NewAdminHandlers(config *config.Config) *AdminHandlers {
	return &AdminHandlers{
		externalService: services.New(config),
	}
}

// GetUsers retrieves all users
func (ah *AdminHandlers) GetUsers(c *gin.Context) {
	response, err := ah.externalService.Call("central", "GET", "/admin/users", nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetUserByID retrieves a specific user by ID
func (ah *AdminHandlers) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	endpoint := "/admin/users/" + id

	response, err := ah.externalService.Call("central", "GET", endpoint, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateUser creates a new user
func (ah *AdminHandlers) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	response, err := ah.externalService.Call("central", "POST", "/admin/users", req)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateUser updates an existing user
func (ah *AdminHandlers) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	endpoint := "/admin/users/" + id

	var req models.UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	response, err := ah.externalService.Call("central", "PUT", endpoint, req)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteUser deletes a user
func (ah *AdminHandlers) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	endpoint := "/admin/users/" + id

	response, err := ah.externalService.Call("central", "DELETE", endpoint, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetRoles retrieves all roles
func (ah *AdminHandlers) GetRoles(c *gin.Context) {
	response, err := ah.externalService.Call("central", "GET", "/admin/roles", nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// AssignRole assigns a role to a user
func (ah *AdminHandlers) AssignRole(c *gin.Context) {
	id := c.Param("id")
	endpoint := "/admin/users/" + id + "/roles"

	var req models.AssignRoleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	response, err := ah.externalService.Call("central", "POST", endpoint, req)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// RemoveRole removes a role from a user
func (ah *AdminHandlers) RemoveRole(c *gin.Context) {
	id := c.Param("id")
	role := c.Param("role")
	endpoint := "/admin/users/" + id + "/roles/" + role

	response, err := ah.externalService.Call("central", "DELETE", endpoint, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetSystemStats retrieves system statistics
func (ah *AdminHandlers) GetSystemStats(c *gin.Context) {
	response, err := ah.externalService.Call("central", "GET", "/admin/system/stats", nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAuditLogs retrieves audit logs
func (ah *AdminHandlers) GetAuditLogs(c *gin.Context) {
	response, err := ah.externalService.Call("central", "GET", "/admin/audit-logs", nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
