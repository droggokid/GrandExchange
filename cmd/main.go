package main

import (
	"context"
	"log"
	"os"

	osrsClient "PaginationPlayground/internal/client"
	"PaginationPlayground/internal/handler"
	"PaginationPlayground/internal/persist"
	"PaginationPlayground/internal/service"
	"PaginationPlayground/temporal"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	dbContext := persist.NewDatabaseContext(ctx, os.Getenv("DB_URL"))
	itemRepo := persist.NewItemRepository(dbContext)

	itemClient := osrsClient.NewOsrsClient()

	itemService := service.NewOsrsService(itemRepo, itemClient)

	temporalClient, err := temporal.NewTemporalClient(itemService)
	if err != nil {
		log.Fatal("Unable to create Temporal client", err)
	}

	itemHandler := handler.NewOsrsHandler(itemService, temporalClient)

	r := gin.Default()

	r.GET("/fetch-osrs-data", itemHandler.FetchAndPersistItems)

	r.GET("/search-item/:name", itemHandler.SearchItems)

	_ = r.Run(":8080")
}
