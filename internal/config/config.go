package config

import (
	"os"
	"strconv"
	"time"
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
	JWT        struct {
		SecretKey string
		ExpiresIn time.Duration
		Issuer    string
	}
	Admin struct {
		Token string
	}
}

func Load() *Config {
	cfg := Config{
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
	cfg.JWT.SecretKey = GetEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production")
	expiresIn, _ := strconv.Atoi(GetEnv("JWT_EXPIRES_IN", "24"))
	cfg.JWT.ExpiresIn = time.Duration(expiresIn) * time.Hour
	cfg.JWT.Issuer = GetEnv("JWT_ISSUER", "pr-reviewer-service")

	cfg.Admin.Token = GetEnv("ADMIN_TOKEN", "admin-secret-token-change-in-production")
	return &cfg
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
