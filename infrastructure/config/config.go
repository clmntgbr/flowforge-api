package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL                 string
	ClerkWebhookSecret          string
	Port                        string
	Environment                 string
	ClerkSecretKey              string
	ClerkFrontendAPI            string
	RabbitMQURL                 string
	RabbitMQSecretKey           string
	CORSAllowedOrigins          []string
	CORSAllowCredentials        bool
	CORSAllowMethods            []string
	CORSAllowHeaders            []string
	CORSMaxAge                  int
	RateLimitMax                int
	ExchangeName                string
	ConsumerQueueName           string
	ConsumerRoutingKeyCompleted string
	ConsumerRoutingKeyFailed    string
	RunnerQueueName             string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, using environment variables")
	}

	return &Config{
		DatabaseURL:                 getEnv("DATABASE_URL"),
		ClerkWebhookSecret:          getEnv("CLERK_WEBHOOK_SECRET"),
		Port:                        getEnv("PORT"),
		Environment:                 getEnv("GO_ENV"),
		ClerkSecretKey:              getEnv("CLERK_SECRET_KEY"),
		ClerkFrontendAPI:            getEnv("CLERK_FRONTEND_API"),
		RabbitMQURL:                 getEnv("RABBITMQ_URL"),
		RabbitMQSecretKey:           getEnv("RABBITMQ_SECRET_KEY"),
		CORSAllowedOrigins:          strings.Split(getEnv("CORS_ALLOWED_ORIGINS"), ","),
		CORSAllowCredentials:        getEnvBool("CORS_ALLOW_CREDENTIALS"),
		CORSAllowMethods:            strings.Split(getEnv("CORS_ALLOW_METHODS"), ","),
		CORSAllowHeaders:            strings.Split(getEnv("CORS_ALLOW_HEADERS"), ","),
		CORSMaxAge:                  getEnvInt("CORS_MAX_AGE"),
		RateLimitMax:                getEnvInt("RATE_LIMIT_MAX"),
		ExchangeName:                getEnv("EXCHANGE_NAME"),
		ConsumerQueueName:           getEnv("CONSUMER_QUEUE_NAME"),
		RunnerQueueName:             getEnv("RUNNER_QUEUE_NAME"),
		ConsumerRoutingKeyCompleted: getEnv("CONSUMER_ROUTING_KEY_COMPLETED"),
		ConsumerRoutingKeyFailed:    getEnv("CONSUMER_ROUTING_KEY_FAILED"),
	}
}

func getEnv(key string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	log.Panicf("required environment variable %s is not set", key)
	return ""
}

func getEnvBool(key string) bool {
	value := os.Getenv(key)
	if value == "" {
		return false
	}

	return value == "true"
}

func getEnvInt(key string) int {
	value := os.Getenv(key)
	if value == "" {
		log.Panicf("required environment variable %s is not set", key)
		return 0
	}

	parsedValue, err := strconv.Atoi(value)
	if err != nil {
		log.Panicf("invalid integer for %s: %q", key, value)
		return 0
	}

	return parsedValue
}
