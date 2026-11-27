// Package handler consists of item related http layer
package handler

import (
	"fmt"
	"net/http"

	"GrandExchange/internal/models"
	"GrandExchange/internal/service"
	"GrandExchange/temporal"

	"github.com/gin-gonic/gin"
)

type ItemHandler interface {
	fetchItems(*gin.Context, string) (models.SearchResponse, error)
	FetchAndPersistItems(*gin.Context)
	SearchItems(*gin.Context)
}

type OsrsHandler struct {
	itemService    service.ItemService
	cacheService   service.CacheService
	temporalClient temporal.TemporalClient
}

func NewOsrsHandler(service service.ItemService, cacheService service.CacheService, temporalClient temporal.TemporalClient) ItemHandler {
	return &OsrsHandler{service, cacheService, temporalClient}
}

func (h *OsrsHandler) SearchItems(c *gin.Context) {
	itemName := c.Param("name")

	cache, err := h.cacheService.Get(c.Request.Context(), itemName)
	if err != nil {
		fmt.Println("Cache error:", err)
	}
	if err == nil && cache != nil {
		items, err := h.temporalClient.StartSearchWorkflow(c.Request.Context(), itemName)
		if err != nil {
			fmt.Println("Temporal workflow error:", err)
			items, dbErr := h.itemService.SearchForItems(c.Request.Context(), itemName)
			if dbErr != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": dbErr.Error()})
				return
			}
			c.JSON(http.StatusOK, models.SearchResponse{Total: len(items), Items: items})
			return
		}
		fmt.Println("Cache hit")
		c.JSON(http.StatusOK, models.SearchResponse{Total: len(items), Items: items})
	} else {
		h.FetchAndPersistItems(c)
	}
}

func (h *OsrsHandler) FetchAndPersistItems(c *gin.Context) {
	itemName := c.Param("name")

	out, err := h.fetchItems(c, itemName)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := h.itemService.PersistSearchResponse(c.Request.Context(), out); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, out)
}

func (h *OsrsHandler) fetchItems(c *gin.Context, alpha string) (models.SearchResponse, error) {
	category := "1"
	page := "1"

	resp, err := h.itemService.FetchItems(c.Request.Context(), category, alpha, page)
	if err != nil {
		return models.SearchResponse{}, err
	}
	if err := h.cacheService.Set(c.Request.Context(), alpha, resp); err != nil {
		fmt.Println("Cache set error:", err)
	}
	return models.SearchResponse{Total: resp.Total, Items: resp.Items}, nil
}
