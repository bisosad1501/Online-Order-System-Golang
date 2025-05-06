package config

import (
"os"
"strconv"
"time"
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
OrderServiceURL string

// Redis configuration
RedisHost     string
RedisPort     string
RedisPassword string
RedisDB       int
RedisCacheTTL time.Duration
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
return &Config{
// Server configuration
ServerPort: getEnv("PORT", "8084"),

// Database configuration
DBHost:     getEnv("DB_HOST", "localhost"),
DBPort:     getEnv("DB_PORT", "5432"),
DBUser:     getEnv("DB_USER", "postgres"),
DBPassword: getEnv("DB_PASSWORD", "postgres"),
DBName:     getEnv("DB_NAME", "shippingdb"),

// Kafka configuration
KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092"),
KafkaTopic:            getEnv("KAFKA_TOPIC", "shipments"),

// External services
OrderServiceURL: getEnv("ORDER_SERVICE_URL", "http://order-service:8081"),

// Redis configuration
RedisHost:     getEnv("REDIS_HOST", "localhost"),
RedisPort:     getEnv("REDIS_PORT", "6379"),
RedisPassword: getEnv("REDIS_PASSWORD", ""),
RedisDB:       getEnvAsInt("REDIS_DB", 2),
RedisCacheTTL: time.Duration(getEnvAsInt("REDIS_CACHE_TTL", 3600)) * time.Second,
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
