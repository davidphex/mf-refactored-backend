package handlers

import (
	"net/http"

	"github.com/davidphex/memoryframe-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	service services.PagesService
}

func NewPageHandler(service services.PagesService) *PageHandler {
	return &PageHandler{service: service}
}

func (h *PageHandler) GetPages(c *gin.Context) {
	albumId := c.Param("id")

	if albumId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album ID parameter is required"})
		return
	}

	pages, err := h.service.GetPages(albumId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if pages == nil || len(*pages) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "No pages found for the specified album"})
		return
	}

	c.JSON(http.StatusOK, pages)
}
