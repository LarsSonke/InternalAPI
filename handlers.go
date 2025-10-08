package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// Error response structure
type ErrorResponse struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Timestamp int64  `json:"timestamp"`
}

// Album data structure for demo endpoints
type Album struct {
	ID     string  `json:"id"`
	Title  string  `json:"title"`
	Artist string  `json:"artist"`
	Price  float64 `json:"price"`
}

// Send standardized error response
func sendError(c *gin.Context, statusCode int, code, message string) {
	c.JSON(statusCode, ErrorResponse{
		Code:      code,
		Message:   message,
		Timestamp: time.Now().Unix(),
	})
}

// Get albums handler - will be replaced with hotel bookings/rooms
func getAlbums(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	userID, _ := c.Get("userID")
	
	log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"action":     "get_albums",
	}).Info("Processing get albums request")

	// STEP 1: Check permissions with Central Management System
	permissionCheck := map[string]interface{}{
		"userID":   userID,
		"action":   "read_albums",
		"resource": "albums",
	}

	permission, err := callCentralManagement("POST", "/check-permission", permissionCheck)
	if err != nil {
		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    userID,
		}).WithError(err).Error("Error checking permissions with Central Management")
		sendError(c, http.StatusInternalServerError, "PERMISSION_SERVICE_ERROR", "Failed to check permissions")
		return
	}

	if !permission["allowed"].(bool) {
		reason := "Permission denied"
		if msg, exists := permission["reason"]; exists {
			reason = msg.(string)
		}
		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    userID,
			"reason":     reason,
		}).Warn("Permission denied for get albums")
		sendError(c, http.StatusForbidden, "PERMISSION_DENIED", reason)
		return
	}

	// STEP 2: Call API Beheerder to get albums data
	result, err := callAPIBeheerder("GET", "/albums", nil)
	if err != nil {
		log.WithFields(logrus.Fields{
			"request_id": requestID,
			"user_id":    userID,
		}).WithError(err).Error("Error calling API Beheerder")
		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to fetch albums from data service")
		return
	}

	// Return data from API Beheerder
	c.JSON(http.StatusOK, result)
}

// Post albums handler - will be replaced with create booking/room
func postAlbums(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	userID, _ := c.Get("userID")
	
	log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"action":     "create_album",
	}).Info("Processing create album request")

	var newAlbum Album
	if err := c.BindJSON(&newAlbum); err != nil {
		log.WithField("request_id", requestID).WithError(err).Warn("Invalid JSON in request body")
		sendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Check permissions with Central Management
	permissionCheck := map[string]interface{}{
		"userID":   userID,
		"action":   "create_albums",
		"resource": "albums",
	}

	permission, err := callCentralManagement("POST", "/check-permission", permissionCheck)
	if err != nil {
		log.WithError(err).Error("Error checking permissions with Central Management")
		sendError(c, http.StatusInternalServerError, "PERMISSION_SERVICE_ERROR", "Failed to check permissions")
		return
	}

	if !permission["allowed"].(bool) {
		sendError(c, http.StatusForbidden, "PERMISSION_DENIED", "Permission denied")
		return
	}

	// Forward to API Beheerder
	result, err := callAPIBeheerder("POST", "/albums", newAlbum)
	if err != nil {
		log.WithError(err).Error("Error calling API Beheerder")
		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to create album")
		return
	}

	c.JSON(http.StatusCreated, result)
}

// Get album by ID handler - will be replaced with specific booking/room
func getAlbumByID(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	userID, _ := c.Get("userID")
	id := c.Param("id")
	
	log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"album_id":   id,
		"action":     "get_album_by_id",
	}).Info("Processing get album by ID request")

	// Check permissions
	permissionCheck := map[string]interface{}{
		"userID":     userID,
		"action":     "read_albums",
		"resource":   "albums",
		"resourceID": id,
	}

	permission, err := callCentralManagement("POST", "/check-permission", permissionCheck)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "PERMISSION_SERVICE_ERROR", "Failed to check permissions")
		return
	}

	if !permission["allowed"].(bool) {
		sendError(c, http.StatusForbidden, "PERMISSION_DENIED", "Permission denied")
		return
	}

	// Forward to API Beheerder
	result, err := callAPIBeheerder("GET", "/albums/"+id, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to fetch album")
		return
	}

	c.JSON(http.StatusOK, result)
}

// Update album handler - will be replaced with update booking/room
func updateAlbum(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	
	var updatedAlbum Album
	if err := c.BindJSON(&updatedAlbum); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_JSON", "Invalid request body")
		return
	}

	// Check permissions
	permissionCheck := map[string]interface{}{
		"userID":     userID,
		"action":     "update_albums",
		"resource":   "albums",
		"resourceID": id,
	}

	permission, err := callCentralManagement("POST", "/check-permission", permissionCheck)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "PERMISSION_SERVICE_ERROR", "Failed to check permissions")
		return
	}

	if !permission["allowed"].(bool) {
		sendError(c, http.StatusForbidden, "PERMISSION_DENIED", "Permission denied")
		return
	}

	// Forward to API Beheerder
	result, err := callAPIBeheerder("PUT", "/albums/"+id, updatedAlbum)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to update album")
		return
	}

	c.JSON(http.StatusOK, result)
}

// Delete album handler - will be replaced with cancel booking/delete room
func deleteAlbum(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	
	// Check permissions
	permissionCheck := map[string]interface{}{
		"userID":     userID,
		"action":     "delete_albums",
		"resource":   "albums",
		"resourceID": id,
	}

	permission, err := callCentralManagement("POST", "/check-permission", permissionCheck)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "PERMISSION_SERVICE_ERROR", "Failed to check permissions")
		return
	}

	if !permission["allowed"].(bool) {
		sendError(c, http.StatusForbidden, "PERMISSION_DENIED", "Permission denied")
		return
	}

	// Forward to API Beheerder
	_, err = callAPIBeheerder("DELETE", "/albums/"+id, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "DATA_SERVICE_ERROR", "Failed to delete album")
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Album deleted successfully"})
}

// Get system status handler - admin endpoint
func getSystemStatus(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	userID, _ := c.Get("userID")
	
	log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"action":     "get_system_status",
	}).Info("Admin requesting system status")

	// Return system status information
	c.JSON(http.StatusOK, gin.H{
		"status":               "healthy",
		"api_beheerder_status": "connected",
		"central_mgmt_status":  "connected", 
		"uptime":               "2h 15m",
		"active_connections":   42,
		"requests_today":       1337,
	})
}

// Get audit logs handler - admin endpoint  
func getAuditLogs(c *gin.Context) {
	requestID, _ := c.Get("requestID")
	userID, _ := c.Get("userID")
	
	log.WithFields(logrus.Fields{
		"request_id": requestID,
		"user_id":    userID,
		"action":     "get_audit_logs",
	}).Info("Admin requesting audit logs")

	c.JSON(http.StatusOK, gin.H{
		"logs": []gin.H{
			{
				"timestamp": time.Now().Add(-2 * time.Hour).Unix(),
				"action":    "room_status_updated",
				"user":      "worker_001",
				"room":      "101",
				"details":   "Status changed from occupied to cleaning",
			},
			{
				"timestamp": time.Now().Add(-1 * time.Hour).Unix(),
				"action":    "booking_created",
				"user":      "user_portal",
				"booking":   "BK123456",
				"details":   "New booking for 3 nights",
			},
		},
		"total": 2,
	})
}