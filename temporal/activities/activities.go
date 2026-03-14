// Package activities includes all temporal item activities
package activities

import (
	"context"

	"PaginationPlayground/internal/models"
	"PaginationPlayground/internal/service"
)

type ItemActivities interface {
	SearchItemActivity(context.Context, string) (models.SearchActivityResponse, error)
}

type OsrsActivities struct {
	itemService service.ItemService
}

func NewOsrsActivities(service service.ItemService) ItemActivities {
	return &OsrsActivities{service}
}

func (a *OsrsActivities) SearchItemActivity(ctx context.Context, itemName string) (models.SearchActivityResponse, error) {
	items, err := a.itemService.SearchForItems(ctx, itemName)
	if err != nil {
		return models.SearchActivityResponse{}, err
	}
	return models.SearchActivityResponse{Total: len(items), Items: items}, nil
}
