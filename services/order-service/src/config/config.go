package config

import (
"os"
"strconv"
)

// Config holds all configuration for the service
type Config struct {
// Server configuration
ServerPort string

// Database configuration
DBHost     string
DBPort     string
DBUser     string
DBPassword string
DBName     string

// Kafka configuration
KafkaBootstrapServers string
KafkaTopic            string

// External services
InventoryServiceURL      string
PaymentServiceURL        string
ShippingServiceURL       string
RecommendationServiceURL string
NotificationServiceURL   string
UserServiceURL           string
CartServiceURL           string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
return &Config{
// Server configuration
ServerPort: getEnv("PORT", "8081"),

// Database configuration
DBHost:     getEnv("DB_HOST", "localhost"),
DBPort:     getEnv("DB_PORT", "5432"),
DBUser:     getEnv("DB_USER", "postgres"),
DBPassword: getEnv("DB_PASSWORD", "postgres"),
DBName:     getEnv("DB_NAME", "orderdb"),

// Kafka configuration
KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092"),
KafkaTopic:            getEnv("KAFKA_TOPIC", "orders"),

// External services
InventoryServiceURL:      getEnv("INVENTORY_SERVICE_URL", "http://inventory-service:8082"),
PaymentServiceURL:        getEnv("PAYMENT_SERVICE_URL", "http://payment-service:8083"),
ShippingServiceURL:       getEnv("SHIPPING_SERVICE_URL", "http://shipping-service:8084"),
RecommendationServiceURL: getEnv("RECOMMENDATION_SERVICE_URL", "http://recommendation-service:8088"),
NotificationServiceURL:   getEnv("NOTIFICATION_SERVICE_URL", "http://notification-service:8085"),
UserServiceURL:           getEnv("USER_SERVICE_URL", "http://user-service:8086"),
CartServiceURL:           getEnv("CART_SERVICE_URL", "http://cart-service:8087"),
}
}

// Simple helper function to read an environment variable or return a default value
func getEnv(key, defaultValue string) string {
if value, exists := os.LookupEnv(key); exists {
return value
}
return defaultValue
}

// Helper to read an environment variable into an integer or return a default value
func getEnvAsInt(key string, defaultValue int) int {
valueStr := getEnv(key, "")
if value, err := strconv.Atoi(valueStr); err == nil {
return value
}
return defaultValue
}

// Helper to read an environment variable into a boolean or return a default value
func getEnvAsBool(key string, defaultValue bool) bool {
valueStr := getEnv(key, "")
if value, err := strconv.ParseBool(valueStr); err == nil {
return value
}
return defaultValue
}
