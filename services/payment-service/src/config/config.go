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

// Payment gateway configuration
PaymentGatewayURL string
PaymentGatewayKey string

// Stripe configuration
StripeSecretKey      string
StripePublishableKey string
StripeWebhookSecret  string

// Payment mode (mock or stripe)
PaymentMode string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
return &Config{
// Server configuration
ServerPort: getEnv("PORT", "8083"),

// Database configuration
DBHost:     getEnv("DB_HOST", "localhost"),
DBPort:     getEnv("DB_PORT", "5432"),
DBUser:     getEnv("DB_USER", "postgres"),
DBPassword: getEnv("DB_PASSWORD", "postgres"),
DBName:     getEnv("DB_NAME", "paymentdb"),

// Kafka configuration
KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092"),
KafkaTopic:            getEnv("KAFKA_TOPIC", "payments"),

// Payment gateway configuration
PaymentGatewayURL: getEnv("PAYMENT_GATEWAY_URL", "https://api.example.com/payments"),
PaymentGatewayKey: getEnv("PAYMENT_GATEWAY_KEY", "test_key"),

// Stripe configuration
StripeSecretKey:      getEnv("STRIPE_SECRET_KEY", "sk_test_your_test_key"),
StripePublishableKey: getEnv("STRIPE_PUBLISHABLE_KEY", "pk_test_your_test_key"),
StripeWebhookSecret:  getEnv("STRIPE_WEBHOOK_SECRET", "whsec_your_webhook_secret"),

// Payment mode (mock or stripe)
PaymentMode: getEnv("PAYMENT_MODE", "mock"),
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
