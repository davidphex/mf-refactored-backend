package handlers

import (
	"net/http"

	"github.com/davidphex/memoryframe-backend/internal/models"
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

func (h *AlbumHandler) InsertAlbum(c *gin.Context) {
	var album models.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.InsertAlbum(&album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, album)
}

func (h *AlbumHandler) UpdateAlbum(c *gin.Context) {
	var album models.Album

	if err := c.ShouldBindJSON(&album); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if album.ID.IsZero() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album ID is required"})
		return
	}

	if err := h.service.UpdateAlbum(&album); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, album)
}

func (h *AlbumHandler) DeleteAlbum(c *gin.Context) {
	id := c.Param("id")

	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID parameter is required"})
		return
	}

	if err := h.service.DeleteAlbum(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
