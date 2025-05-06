package main

import (
"context"
"log"
"net/http"
"os"
"os/signal"
"syscall"
"time"

"github.com/online-order-system/order-service/api"
"github.com/online-order-system/order-service/config"
"github.com/online-order-system/order-service/db"
"github.com/online-order-system/order-service/kafka"
"github.com/online-order-system/order-service/service"
)

func main() {
// Load configuration
cfg := config.LoadConfig()

// Connect to database
database, err := db.NewDatabase(cfg)
if err != nil {
log.Fatalf("Failed to connect to database: %v", err)
}

// Create tables
err = database.CreateTables()
if err != nil {
log.Fatalf("Failed to create tables: %v", err)
}

// Create repository
repository := db.NewOrderRepository(database)

// Create Kafka producer
producer := kafka.NewProducer(cfg)
defer producer.Close()

// Create service
orderService := service.NewOrderService(cfg, repository, producer)

// Set order service instance for direct compensation
service.SetOrderServiceInstance(orderService)

// Create Kafka consumer
consumer := kafka.NewConsumer(cfg, orderService)
paymentConsumer := kafka.NewPaymentConsumer(cfg, orderService)

// Start Kafka consumers
ctx, cancel := context.WithCancel(context.Background())
defer cancel()
consumer.StartConsuming(ctx)
paymentConsumer.StartConsuming(ctx)

// Setup router
// Use the new router setup
router := api.SetupRouter(orderService)

// Start server
srv := &http.Server{
Addr:    ":" + cfg.ServerPort,
Handler: router,
}

// Start server in a goroutine
go func() {
log.Printf("Order Service starting on port %s...", cfg.ServerPort)
if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
log.Fatalf("Failed to start server: %v", err)
}
}()

// Wait for interrupt signal to gracefully shutdown the server
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit
log.Println("Shutting down server...")

// Cancel context to stop Kafka consumer
cancel()

// Shutdown server with a timeout
ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
if err := srv.Shutdown(ctx); err != nil {
log.Fatalf("Server forced to shutdown: %v", err)
}

log.Println("Server exiting")
}
