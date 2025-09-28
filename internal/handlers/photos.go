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

	err = h.service.UploadPhoto(fileHeader, albumId)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to upload photo: " + err.Error()})
		return
	}

	c.JSON(200, gin.H{"message": "Photo uploaded successfully"})

}
