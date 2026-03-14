// Package client contains all item http logic
package client

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"time"

	"PaginationPlayground/internal/models"
)

type ItemClient interface {
	FetchOsrsData(context.Context, string, string, string) (models.SearchResponse, error)
}

type OsrsClient struct {
	client http.Client
}

func NewOsrsClient() ItemClient {
	return &OsrsClient{
		client: http.Client{
			Timeout: time.Second * 10,
		},
	}
}

func (h *OsrsClient) FetchOsrsData(ctx context.Context, category, alpha, page string) (models.SearchResponse, error) {
	u, _ := url.Parse("https://secure.runescape.com/m=itemdb_oldschool/api/catalogue/items.json")
	q := u.Query()
	q.Set("category", category)
	q.Set("alpha", alpha)
	q.Set("page", page)
	u.RawQuery = q.Encode()

	req, _ := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	req.Header.Set("User-Agent", "PaginationPlayground/1.0")
	req.Header.Set("Accept", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return models.SearchResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return models.SearchResponse{}, nil
	}

	var out models.SearchResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return models.SearchResponse{}, err
	}

	return out, nil
}
