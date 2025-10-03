package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Mock API Beheerder for testing your Internal API
// This simulates what the real API Beheerder would do

type Album struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Artist    string  `json:"artist"`
	Price     float64 `json:"price"`
	CreatedBy string  `json:"createdBy,omitempty"`
	CreatedAt int64   `json:"createdAt,omitempty"`
	UpdatedBy string  `json:"updatedBy,omitempty"`
	UpdatedAt int64   `json:"updatedAt,omitempty"`
}

// Mock database
var mockAlbums = []Album{
	{ID: "1", Title: "Blue Train", Artist: "John Coltrane", Price: 56.99, CreatedAt: time.Now().Unix()},
	{ID: "2", Title: "Jeru", Artist: "Gerry Mulligan", Price: 17.99, CreatedAt: time.Now().Unix()},
	{ID: "3", Title: "Sarah Vaughan and Clifford Brown", Artist: "Sarah Vaughan", Price: 39.99, CreatedAt: time.Now().Unix()},
}

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()

	// Middleware to validate service key from your Internal API
	router.Use(func(c *gin.Context) {
		serviceKey := c.GetHeader("X-Service-Key")
		if serviceKey != "beheerder-service-key" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid service key"})
			c.Abort()
			return
		}
		c.Next()
	})

	// Albums endpoints
	router.GET("/albums", getAlbums)
	router.GET("/albums/:id", getAlbumByID)
	router.POST("/albums", createAlbum)
	router.PUT("/albums/:id", updateAlbum)
	router.DELETE("/albums/:id", deleteAlbum)

	fmt.Println("🔧 Mock API Beheerder starting on :8081")
	fmt.Println("   🎯 Ready to receive requests from your Internal API")
	router.Run(":8081")
}

func getAlbums(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"albums": mockAlbums,
		"count":  len(mockAlbums),
	})
}

func getAlbumByID(c *gin.Context) {
	id := c.Param("id")

	for _, album := range mockAlbums {
		if album.ID == id {
			c.JSON(http.StatusOK, gin.H{
				"album": album,
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
}

func createAlbum(c *gin.Context) {
	var albumData map[string]interface{}

	if err := c.ShouldBindJSON(&albumData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	// Convert to Album struct
	album := Album{
		ID:        albumData["id"].(string),
		Title:     albumData["title"].(string),
		Artist:    albumData["artist"].(string),
		Price:     albumData["price"].(float64),
		CreatedAt: int64(albumData["createdAt"].(float64)),
	}

	if createdBy, exists := albumData["createdBy"]; exists {
		album.CreatedBy = createdBy.(string)
	}

	// Check for duplicate
	for _, existing := range mockAlbums {
		if existing.ID == album.ID {
			c.JSON(http.StatusConflict, gin.H{"error": "Album already exists"})
			return
		}
	}

	mockAlbums = append(mockAlbums, album)

	c.JSON(http.StatusCreated, gin.H{
		"album":   album,
		"message": "Album created in data service",
	})
}

func updateAlbum(c *gin.Context) {
	id := c.Param("id")
	var albumData map[string]interface{}

	if err := c.ShouldBindJSON(&albumData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	for i, album := range mockAlbums {
		if album.ID == id {
			// Update the album
			mockAlbums[i].Title = albumData["title"].(string)
			mockAlbums[i].Artist = albumData["artist"].(string)
			mockAlbums[i].Price = albumData["price"].(float64)
			mockAlbums[i].UpdatedAt = int64(albumData["updatedAt"].(float64))

			if updatedBy, exists := albumData["updatedBy"]; exists {
				mockAlbums[i].UpdatedBy = updatedBy.(string)
			}

			c.JSON(http.StatusOK, gin.H{
				"album":   mockAlbums[i],
				"message": "Album updated in data service",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
}

func deleteAlbum(c *gin.Context) {
	id := c.Param("id")

	for i, album := range mockAlbums {
		if album.ID == id {
			// Remove the album
			mockAlbums = append(mockAlbums[:i], mockAlbums[i+1:]...)
			c.JSON(http.StatusOK, gin.H{
				"message": "Album deleted from data service",
			})
			return
		}
	}

	c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
}
