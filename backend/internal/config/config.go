package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Host         string
	Port         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	URL             string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	URL string
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	Secret            string
	AccessExpiration  time.Duration
	RefreshExpiration time.Duration
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Host:         getEnv("SERVER_HOST", "0.0.0.0"),
			Port:         getEnv("SERVER_PORT", "8080"),
			ReadTimeout:  getDurationEnv("SERVER_READ_TIMEOUT", 10*time.Second),
			WriteTimeout: getDurationEnv("SERVER_WRITE_TIMEOUT", 10*time.Second),
		},
		Database: DatabaseConfig{
			URL:             getEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/lightshare?sslmode=disable"),
			MaxOpenConns:    getIntEnv("DATABASE_MAX_OPEN_CONNS", 25),
			MaxIdleConns:    getIntEnv("DATABASE_MAX_IDLE_CONNS", 5),
			ConnMaxLifetime: getDurationEnv("DATABASE_CONN_MAX_LIFETIME", 5*time.Minute),
		},
		Redis: RedisConfig{
			URL: getEnv("REDIS_URL", "redis://localhost:6379"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET", "development-secret-change-in-production"),
			AccessExpiration:  getDurationEnv("JWT_ACCESS_EXPIRATION", 1*time.Hour),
			RefreshExpiration: getDurationEnv("JWT_REFRESH_EXPIRATION", 30*24*time.Hour),
		},
	}
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getIntEnv gets an integer environment variable or returns a default value
func getIntEnv(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getDurationEnv gets a duration environment variable or returns a default value
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	if value := os.Getenv(key); value != "" {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
	}
	return defaultValue
}
