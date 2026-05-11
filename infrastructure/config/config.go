package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL        string
	ClerkWebhookSecret string
	Port               string
	Environment        string
	ClerkSecretKey     string
	ClerkFrontendAPI   string
	RabbitMQURL        string
	RabbitMQSecretKey  string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	return &Config{
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		ClerkWebhookSecret: getEnv("CLERK_WEBHOOK_SECRET", ""),
		Port:               getEnv("PORT", "3000"),
		Environment:        getEnv("GO_ENV", "development"),
		ClerkSecretKey:     getEnv("CLERK_SECRET_KEY", ""),
		ClerkFrontendAPI:   getEnv("CLERK_FRONTEND_API", ""),
		RabbitMQURL:        getEnv("RABBITMQ_URL", ""),
		RabbitMQSecretKey:  getEnv("RABBITMQ_SECRET_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
