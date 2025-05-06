package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/online-order-system/inventory-service/api"
	"github.com/online-order-system/inventory-service/cache"
	"github.com/online-order-system/inventory-service/config"
	"github.com/online-order-system/inventory-service/db"
	"github.com/online-order-system/inventory-service/kafka"
	"github.com/online-order-system/inventory-service/service"
)

func main() {
	// Load configuration
	cfg := config.LoadConfig()

	// Connect to database
	database, err := db.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	// Create tables
	err = database.CreateTables()
	if err != nil {
		log.Fatalf("Failed to create tables: %v", err)
	}

	// Create repository
	repository := db.NewInventoryRepository(database)

	// Create Redis cache
	redisCache, err := cache.NewRedisCache(cfg)
	if err != nil {
		log.Printf("Failed to connect to Redis: %v", err)
		log.Println("Continuing without Redis cache...")
		redisCache = nil
	} else {
		log.Println("Connected to Redis cache")
		defer redisCache.Close()
	}

	// Create Kafka producer
	producer := kafka.NewProducer(cfg)
	defer producer.Close()

	// Create service
	inventoryService := service.NewInventoryService(cfg, repository, producer, redisCache)

	// Create Kafka consumer
	consumer := kafka.NewConsumer(cfg, inventoryService)

	// Start Kafka consumer
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	consumer.StartConsuming(ctx)

	// Set up router
	router := api.SetupRouter(inventoryService)

	// Start cache refresh goroutine
	if redisCache != nil {
		go func() {
			// Refresh cache every 5 minutes
			ticker := time.NewTicker(5 * time.Minute)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					log.Println("Performing scheduled cache refresh...")
					// Delete products:all cache to force refresh on next query
					err := redisCache.Delete(ctx, "products:all")
					if err != nil {
						log.Printf("Failed to refresh products cache: %v", err)
					} else {
						log.Println("Successfully refreshed products cache")
					}
				}
			}
		}()
		log.Println("Started automatic cache refresh (every 5 minutes)")
	}

	// Start server
	go func() {
		if err := router.Run(":" + cfg.ServerPort); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	log.Printf("Inventory Service starting on port %s...", cfg.ServerPort)

	// Wait for interrupt signal to gracefully shut down the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	// Cancel context to stop Kafka consumer
	cancel()

	// Give Kafka consumer some time to stop
	time.Sleep(1 * time.Second)

	log.Println("Server exiting")
}
