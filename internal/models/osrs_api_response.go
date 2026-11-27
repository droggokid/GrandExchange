// Package models contains all models
package models

import "encoding/json"

type SearchActivityResponse struct {
	Total int          `json:"total"`
	Items []SearchItem `json:"items"`
}

type SearchResponse struct {
	Total int          `json:"total"`
	Items []SearchItem `json:"items"`
}

type SearchItem struct {
	Icon        string   `json:"icon"`
	IconLarge   string   `json:"icon_large"`
	ID          int      `json:"id"`
	Type        string   `json:"type"`
	TypeIcon    string   `json:"typeIcon"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Current     PriceBox `json:"current"`
	Today       PriceBox `json:"today"`
	Members     string   `json:"members"`
}

type PriceBox struct {
	Trend string          `json:"trend"`
	Price json.RawMessage `json:"price"`
}
