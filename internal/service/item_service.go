// Package service contains item related business logic
package service

import (
	"context"

	"PaginationPlayground/internal/client"
	"PaginationPlayground/internal/models"
	"PaginationPlayground/internal/persist"
)

type ItemService interface {
	FetchItems(context.Context, string, string, string) (models.SearchResponse, error)
	SearchForItems(context.Context, string) ([]models.SearchItem, error)
	PersistSearchResponse(context.Context, models.SearchResponse) error
}

type OsrsService struct {
	itemRepo   persist.ItemRepository
	itemClient client.ItemClient
}

func NewOsrsService(repo persist.ItemRepository, client client.ItemClient) ItemService {
	return &OsrsService{repo, client}
}

func (s *OsrsService) FetchItems(ctx context.Context, category string, alpha string, page string) (models.SearchResponse, error) {
	resp, err := s.itemClient.FetchOsrsData(ctx, category, alpha, page)
	if err != nil {
		return models.SearchResponse{}, err
	}
	return resp, nil
}

func (s *OsrsService) SearchForItems(ctx context.Context, itemName string) ([]models.SearchItem, error) {
	items, err := s.itemRepo.GetItem(ctx, itemName)
	if err != nil {
		return nil, err
	}
	return items, nil
}

func (s *OsrsService) PersistSearchResponse(ctx context.Context, response models.SearchResponse) error {
	err := s.itemRepo.SaveItems(ctx, response.Items)
	if err != nil {
		return err
	}
	return nil
}
