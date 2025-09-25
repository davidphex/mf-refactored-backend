package handlers

import (
	"net/http"

	"github.com/davidphex/memoryframe-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type AlbumHandler struct {
	service services.AlbumService
}

func NewAlbumHandler(service services.AlbumService) *AlbumHandler {
	return &AlbumHandler{service: service}
}

func (h *AlbumHandler) GetAllAlbums(c *gin.Context) {
	albums, err := h.service.GetAllAlbums()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, albums)
}

func (h *AlbumHandler) GetAlbumByID(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	album, err := h.service.GetAlbumByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if album == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Album not found"})
		return
	}

	c.JSON(http.StatusOK, album)

}
