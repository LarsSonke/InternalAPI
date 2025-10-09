package handlers

import (
	"net/http"

	"InternalAPI/internal/config"
	"InternalAPI/internal/models"
	"InternalAPI/internal/services"

	"github.com/gin-gonic/gin"
)

// AlbumHandlers contains all album-related handlers
type AlbumHandlers struct {
	externalService *services.ExternalService
}

// NewAlbumHandlers creates a new album handlers instance
func NewAlbumHandlers(config *config.Config) *AlbumHandlers {
	return &AlbumHandlers{
		externalService: services.New(config),
	}
}

// GetAlbums retrieves all albums
func (ah *AlbumHandlers) GetAlbums(c *gin.Context) {
	response, err := ah.externalService.Call("beheerder", "GET", "/albums", nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// GetAlbumByID retrieves a specific album by ID
func (ah *AlbumHandlers) GetAlbumByID(c *gin.Context) {
	id := c.Param("id")
	endpoint := "/albums/" + id

	response, err := ah.externalService.Call("beheerder", "GET", endpoint, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// CreateAlbum creates a new album
func (ah *AlbumHandlers) CreateAlbum(c *gin.Context) {
	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	response, err := ah.externalService.Call("beheerder", "POST", "/albums", album)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusCreated, response)
}

// UpdateAlbum updates an existing album
func (ah *AlbumHandlers) UpdateAlbum(c *gin.Context) {
	id := c.Param("id")
	endpoint := "/albums/" + id

	var album models.Album
	if err := c.ShouldBindJSON(&album); err != nil {
		sendError(c, http.StatusBadRequest, "INVALID_REQUEST", err.Error())
		return
	}

	response, err := ah.externalService.Call("beheerder", "PUT", endpoint, album)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}

// DeleteAlbum deletes an album
func (ah *AlbumHandlers) DeleteAlbum(c *gin.Context) {
	id := c.Param("id")
	endpoint := "/albums/" + id

	response, err := ah.externalService.Call("beheerder", "DELETE", endpoint, nil)
	if err != nil {
		sendError(c, http.StatusInternalServerError, "SERVICE_ERROR", err.Error())
		return
	}

	c.JSON(http.StatusOK, response)
}
