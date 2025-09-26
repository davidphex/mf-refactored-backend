package handlers

import (
	"net/http"

	"github.com/davidphex/memoryframe-backend/internal/models"
	"github.com/davidphex/memoryframe-backend/internal/services"
	"github.com/gin-gonic/gin"
)

type PageHandler struct {
	service services.PagesService
}

func NewPageHandler(service services.PagesService) *PageHandler {
	return &PageHandler{service: service}
}

func (h *PageHandler) GetAlbumPages(c *gin.Context) {
	albumId := c.Param("id")

	if albumId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Album ID parameter is required"})
		return
	}

	pages, err := h.service.GetAlbumPages(albumId)
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

func (h *PageHandler) GetPage(c *gin.Context) {
	pageId := c.Param("id")

	if pageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page ID parameter is required"})
		return
	}

	page, err := h.service.GetPage(pageId)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if page == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Page not found"})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) InsertPage(c *gin.Context) {
	var page models.AlbumPage

	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.InsertPage(&page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, page)
}

func (h *PageHandler) UpdatePage(c *gin.Context) {
	var page models.AlbumPage

	if err := c.ShouldBindJSON(&page); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.UpdatePage(&page); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, page)
}

func (h *PageHandler) DeletePage(c *gin.Context) {
	pageId := c.Param("id")

	if pageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page ID parameter is required"})
		return
	}

	if err := h.service.DeletePage(pageId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page deleted successfully"})
}

type PageElementsUpdateRequest struct {
	Elements []models.AlbumPageElement `json:"elements" binding:"required"`
}

func (h *PageHandler) UpdatePageElements(c *gin.Context) {
	pageId := c.Param("id")

	if pageId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Page ID parameter is required"})
		return
	}

	var req PageElementsUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if len(req.Elements) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Elements array cannot be empty"})
		return
	}

	if err := h.service.UpdatePageElements(pageId, req.Elements); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Page elements updated successfully"})
}
