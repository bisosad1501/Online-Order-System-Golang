package main

import (
"context"
"log"
"os"
"os/signal"
"syscall"
"time"

"github.com/online-order-system/notification-service/api"
"github.com/online-order-system/notification-service/config"
"github.com/online-order-system/notification-service/db"
"github.com/online-order-system/notification-service/kafka"
"github.com/online-order-system/notification-service/service"
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
repository := db.NewNotificationRepository(database)

// Create Kafka producer
producer := kafka.NewProducer(cfg)
defer producer.Close()

// Create service
notificationService := service.NewNotificationService(cfg, repository, producer)

// Create Kafka consumer
consumer := kafka.NewConsumer(cfg, notificationService)

// Start Kafka consumer
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
consumer.StartConsuming(ctx)

// Set up router
router := api.SetupRouter(notificationService)

// Start server
go func() {
if err := router.Run(":" + cfg.ServerPort); err != nil {
log.Fatalf("Failed to start server: %v", err)
}
}()

log.Printf("Notification Service starting on port %s...", cfg.ServerPort)

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
