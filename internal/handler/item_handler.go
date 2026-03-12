// Package handler consists of item related http layer
package handler

import (
	"net/http"

	"PaginationPlayground/internal/models"
	"PaginationPlayground/internal/service"

	"github.com/gin-gonic/gin"
)

type ItemHandler interface {
	FetchItems(*gin.Context) (models.SearchResponse, error)
	FetchAndPersistItems(*gin.Context)
	SearchItems(*gin.Context)
}

type OsrsHandler struct {
	itemService service.ItemService
}

func NewOsrsHandler(service service.ItemService) ItemHandler {
	return &OsrsHandler{service}
}

func (h *OsrsHandler) SearchItems(c *gin.Context) {
	itemName := c.Param("name")
	items, err := h.itemService.SearchForItems(itemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, models.SearchResponse{Total: len(items), Items: items})
}

func (h *OsrsHandler) FetchItems(c *gin.Context) (models.SearchResponse, error) {
	category := c.DefaultQuery("category", "1")
	alpha := c.DefaultQuery("alpha", "c")
	page := c.DefaultQuery("page", "1")

	out, err := h.itemService.FetchItems(category, alpha, page)
	if err != nil {
		return models.SearchResponse{}, err
	}

	return models.SearchResponse{Total: len(out.Items), Items: out.Items}, nil
}

func (h *OsrsHandler) FetchAndPersistItems(c *gin.Context) {
	out, err := h.FetchItems(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.itemService.PersistSearchResponse(out); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}
