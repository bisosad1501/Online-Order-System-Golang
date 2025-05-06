package config

import (
	"os"
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
}

// LoadConfig loads configuration from environment variables
func LoadConfig() *Config {
	return &Config{
		// Server configuration
		ServerPort: getEnv("PORT", "8086"),

		// Database configuration
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5433"),
		DBUser:     getEnv("DB_USER", "useruser"),
		DBPassword: getEnv("DB_PASSWORD", "userpass"),
		DBName:     getEnv("DB_NAME", "userdb"),

		// Kafka configuration
		KafkaBootstrapServers: getEnv("KAFKA_BOOTSTRAP_SERVERS", "kafka:9092"),
		KafkaTopic:            getEnv("KAFKA_TOPIC", "orders"),
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
