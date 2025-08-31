package config

import (
	"os"
)

type Config struct {
	AppPort          string
	JWTSecret        string
	JWTAccessExpiry  string
	JWTRefreshExpiry string
	MySQLHost        string
	MySQLPort        string
	MySQLDatabase    string
	MySQLUsername    string
	MySQLPassword    string
	RedisHost        string
	RedisPort        string
	RedisPassword    string
	RedisDB          string
	GitHubAppID      string
	GitHubPrivateKey string
	GitHubWebhookSecret string
	GitHubToken      string
}

var AppConfig Config

func Load() {
	AppConfig = Config{
		AppPort:          getEnv("APP_PORT", "8080"),
		JWTSecret:        getEnv("JWT_SECRET", "dev-secret"),
		JWTAccessExpiry:  getEnv("JWT_ACCESS_EXPIRY", "24h"),
		JWTRefreshExpiry: getEnv("JWT_REFRESH_EXPIRY", "168h"),
		MySQLHost:        getEnv("MYSQL_HOST", "localhost"),
		MySQLPort:        getEnv("MYSQL_PORT", "3306"),
		MySQLDatabase:    getEnv("MYSQL_DATABASE", "deployment_platform"),
		MySQLUsername:    getEnv("MYSQL_USERNAME", "root"),
		MySQLPassword:    getEnv("MYSQL_PASSWORD", "password"),
		RedisHost:        getEnv("REDIS_HOST", "localhost"),
		RedisPort:        getEnv("REDIS_PORT", "6379"),
		RedisPassword:    getEnv("REDIS_PASSWORD", ""),
		RedisDB:          getEnv("REDIS_DB", "0"),
		GitHubAppID:      os.Getenv("GITHUB_APP_ID"),
		GitHubPrivateKey: os.Getenv("GITHUB_PRIVATE_KEY"),
		GitHubWebhookSecret: os.Getenv("GITHUB_WEBHOOK_SECRET"),
		GitHubToken:      os.Getenv("GITHUB_TOKEN"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}