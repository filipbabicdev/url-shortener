package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	ServerPort        	string
	DatabaseURL       	string
	Env               	string
	DBMaxConns    		int32
	DBMinConns    		int32
	DBConnMaxLifetime 	time.Duration
}

func Load() (*Config, error) {
	return &Config{
		ServerPort:        	getEnv("PORT", "8090"),
		DatabaseURL:       	getEnv("DATABASE_URL", "postgres://url:url@localhost:5432/url?sslmode=disable"),
		Env:               	getEnv("ENV", "development"),
		DBMaxConns:    		getEnvAsInt32("DB_MAX_CONNS", 25),
		DBMinConns:    		getEnvAsInt32("DB_MIN_CONNS", 5),
		DBConnMaxLifetime: 	getEnvAsDuration("DB_CONN_MAX_LIFETIME_MIN", 5*time.Minute),
	}, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt32(key string, fallback int32) int32 {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.ParseInt(value, 10, 32); err == nil {
			return int32(i)
		}
	}
	return fallback
}

func getEnvAsDuration(key string, fallback time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if i, err := strconv.ParseInt(value, 10, 64); err == nil {
			return time.Duration(i) * time.Minute
		}
	}
	return fallback
}