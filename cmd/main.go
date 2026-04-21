package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	osrsClient "GrandExchange/internal/client"
	"GrandExchange/internal/config"
	"GrandExchange/internal/handler"
	"GrandExchange/internal/persist"
	"GrandExchange/internal/service"
	"GrandExchange/temporal"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	ctx := context.Background()

	dbContext := persist.NewDatabaseContext(ctx, config.DBURL)
	defer dbContext.Conn.Close()

	redisContext := persist.NewRedisContext()
	defer redisContext.Rdb.Close()

	itemRepo := persist.NewItemRepository(dbContext)

	itemClient := osrsClient.NewOsrsClient()

	cacheService := service.NewCacheService(redisContext)
	itemService := service.NewOsrsService(itemRepo, itemClient)

	temporalClient, worker, err := temporal.NewTemporalClient(itemService)
	if err != nil {
		log.Fatal("Unable to create Temporal client", err)
	}
	defer worker.Stop()

	itemHandler := handler.NewOsrsHandler(itemService, cacheService, temporalClient)

	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8000"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type"},
		AllowCredentials: true,
	}))

	r.GET("/fetch-osrs-data", itemHandler.FetchAndPersistItems)
	r.GET("/search-item/:name", itemHandler.SearchItems)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = srv.Shutdown(ctx)
}
