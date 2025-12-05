package config

import (
	"os"
	"time"
)

type Config struct {
	Env         string
	Port        string
	DatabaseURL string
	RedisURL    string
	CacheURL    time.Duration
}

func Load() (*Config, error) {
	cfg := &Config{
		Env:         getEnv("ENV", "development"),
		Port:        getEnv("PORT", "50052"),
		DatabaseURL: getEnv("DATABASE_URL", "postgres://user:Bogdan_20@localhost:5432/user_profile_db?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "redis://localhost:6379/1"),
	}

	// Parse cache TTL
	ttlStr := getEnv("CACHE_TTL", "1h")
	ttl, err := time.ParseDuration(ttlStr)
	if err != nil {
		return nil, err
	}
	cfg.CacheURL = ttl

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}

	return defaultValue
}
