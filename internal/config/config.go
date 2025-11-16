package config

import (
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	SSLMode    string
	ServerPort string
	AdminToken string
	UserToken  string
}

func Load() *Config {
	return &Config{
		DBHost:     GetEnv("DB_HOST", "localhost"),
		DBPort:     GetEnv("DB_PORT", "5432"),
		DBUser:     GetEnv("DB_USER", "postgres"),
		DBPassword: GetEnv("DB_PASSWORD", "password"),
		DBName:     GetEnv("DB_NAME", "pr_review"),
		SSLMode:    GetEnv("DB_SSLMODE", "disable"),
		ServerPort: GetEnv("SERVER_PORT", "8080"),
		AdminToken: GetEnv("ADMIN_TOKEN", "admin-secret-token"),
		UserToken:  GetEnv("USER_TOKEN", "user-secret-token"),
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func GetEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
