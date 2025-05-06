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

// Email configuration
SMTPHost     string
SMTPPort     string
SMTPUsername string
SMTPPassword string
SMTPFrom     string

// SMS configuration
SMSAPIKey      string
SMSAPISecret   string
SMSFromNumber  string

// Push notification configuration
PushAPIKey     string

// Service URLs
OrderServiceURL string
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
return &Config{
// Server configuration
ServerPort: getEnv("PORT", "8085"),

// Database configuration
DBHost:     getEnv("DB_HOST", "localhost"),
DBPort:     getEnv("DB_PORT", "5432"),
DBUser:     getEnv("DB_USER", "postgres"),
DBPassword: getEnv("DB_PASSWORD", "postgres"),
DBName:     getEnv("DB_NAME", "notificationdb"),

// Kafka configuration
KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092"),
KafkaTopic:            getEnv("KAFKA_TOPIC", "orders"),

// Email configuration
SMTPHost:     getEnv("SMTP_HOST", "smtp.example.com"),
SMTPPort:     getEnv("SMTP_PORT", "587"),
SMTPUsername: getEnv("SMTP_USERNAME", "user@example.com"),
SMTPPassword: getEnv("SMTP_PASSWORD", "password"),
SMTPFrom:     getEnv("SMTP_FROM", "noreply@example.com"),

// SMS configuration
SMSAPIKey:      getEnv("SMS_API_KEY", "sms_api_key"),
SMSAPISecret:   getEnv("SMS_API_SECRET", "sms_api_secret"),
SMSFromNumber:  getEnv("SMS_FROM_NUMBER", "+1234567890"),

// Push notification configuration
PushAPIKey:     getEnv("PUSH_API_KEY", "push_api_key"),

// Service URLs
OrderServiceURL: getEnv("ORDER_SERVICE_URL", "http://order-service:8081"),
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
