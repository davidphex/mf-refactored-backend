package handlers

import (
	"github.com/davidphex/memoryframe-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type PhotoHandler struct {
	service services.PhotoService
}

func NewPhotoHandler(service services.PhotoService) *PhotoHandler {
	return &PhotoHandler{service: service}
}

func (h *PhotoHandler) UploadPhoto(c *gin.Context) {
	fileHeader, err := c.FormFile("photo")
	if err != nil {
		c.JSON(400, gin.H{"error": "Photo file is required"})
		return
	}

	albumId := c.Param("id")
	if albumId == "" {
		c.JSON(400, gin.H{"error": "Album ID is required"})
		return
	}

	photoId, err := h.service.UploadPhoto(fileHeader, albumId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload photo: " + err.Error()})
		return
	}

	photo, err := h.service.GetPhoto(photoId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch uploaded photo: " + err.Error()})
		return
	}

	// Return the uploaded photo details
	c.JSON(200, photo)
}

func (h *PhotoHandler) GetAlbumPhotos(c *gin.Context) {
	albumId := c.Param("id")
	if albumId == "" {
		c.JSON(400, gin.H{"error": "Album ID is required"})
		return
	}

	photos, err := h.service.GetPhotosByAlbumId(albumId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch photos: " + err.Error()})
		return
	}

	c.JSON(200, photos)
}

func (h *PhotoHandler) GetPhoto(c *gin.Context) {
	photoId := c.Param("id")
	if photoId == "" {
		c.JSON(400, gin.H{"error": "Photo ID is required"})
		return
	}

	photo, err := h.service.GetPhoto(photoId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch photo: " + err.Error()})
		return
	}
	if photo == nil {
		c.JSON(404, gin.H{"error": "Photo not found"})
		return
	}

	c.JSON(200, photo)
}
