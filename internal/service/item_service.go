// Package service contains item related business logic
package service

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"GrandExchange/internal/client"
	"GrandExchange/internal/models"
	"GrandExchange/internal/persist"
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
	firstPage, err := strconv.Atoi(page)
	if err != nil || firstPage < 1 {
		return models.SearchResponse{}, fmt.Errorf("invalid catalogue page %q", page)
	}

	items := []models.SearchItem{}
	for currentPage := firstPage; ; currentPage++ {
		resp, err := s.itemClient.FetchOsrsData(ctx, category, alpha, strconv.Itoa(currentPage))
		if err != nil {
			return models.SearchResponse{}, err
		}
		if len(resp.Items) == 0 {
			return models.SearchResponse{Total: len(items), Items: items}, nil
		}
		items = append(items, resp.Items...)
		time.Sleep(200 * time.Millisecond) // sleep to avoid rate-limit
	}
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
