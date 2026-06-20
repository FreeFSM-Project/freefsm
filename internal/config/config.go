package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

var (
	Version = "dev"
	Commit  = "none"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string

	Addr           string
	LogLevel       string
	LogFile        string
	SessionSecret  string
	SetupToken     string

	UploadDir      string
	MaxUploadSize  int64
}

func Load() (*Config, error) {
	godotenv.Load()

	cfg := &Config{
		DBHost:        getEnv("FREEFSM_DB_HOST", "localhost"),
		DBPort:        getEnvInt("FREEFSM_DB_PORT", 5432),
		DBName:        getEnv("FREEFSM_DB_NAME", "freefsm"),
		DBUser:        getEnv("FREEFSM_DB_USER", "freefsm"),
		DBPassword:    getEnv("FREEFSM_DB_PASSWORD", ""),
		DBSSLMode:     getEnv("FREEFSM_DB_SSLMODE", "disable"),
		Addr:          getEnv("FREEFSM_ADDR", ":3000"),
		LogLevel:      getEnv("FREEFSM_LOG_LEVEL", "info"),
		LogFile:       getEnv("FREEFSM_LOG_FILE", ""),
		SessionSecret: getEnv("FREEFSM_SESSION_SECRET", ""),
		SetupToken:    getEnv("FREEFSM_SETUP_TOKEN", ""),
		UploadDir:     getEnv("FREEFSM_UPLOAD_DIR", "/var/lib/freefsm/uploads"),
		MaxUploadSize: getEnvInt64("FREEFSM_MAX_UPLOAD_SIZE", 26214400),
	}

	if cfg.SessionSecret == "" {
		return nil, fmt.Errorf("FREEFSM_SESSION_SECRET is required")
	}

	return cfg, nil
}

func (c *Config) DSN() string {
	ssl := c.DBSSLMode
	if ssl == "" {
		ssl = "disable"
	}
	return fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName, ssl,
	)
}

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func getEnvInt(key string, def int) int {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.Atoi(v)
		if err == nil {
			return n
		}
	}
	return def
}

func getEnvInt64(key string, def int64) int64 {
	if v := os.Getenv(key); v != "" {
		n, err := strconv.ParseInt(v, 10, 64)
		if err == nil {
			return n
		}
	}
	return def
}
