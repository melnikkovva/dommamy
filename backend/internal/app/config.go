package app

import (
	"errors"
	"os"
	"time"
)

type LearnConfig struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
}

func (c LearnConfig) Enabled() bool {
	return c.BaseURL != "" && c.ClientID != "" && c.ClientSecret != ""
}

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	FilesDir    string
	AuthSecret  string
	AuthTTL     time.Duration
	Learn       LearnConfig
}

func LoadConfig() (Config, error) {
	cfg := Config{
		HTTPAddr:    getenv("HTTP_ADDR", ":8080"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		FilesDir:    getenv("FILES_DIR", "./uploads"),
		AuthSecret:  os.Getenv("AUTH_SECRET"),
		AuthTTL:     24 * time.Hour,
		Learn: LearnConfig{
			BaseURL:      os.Getenv("LEARN_API_BASE_URL"),
			ClientID:     os.Getenv("LEARN_API_CLIENT_ID"),
			ClientSecret: os.Getenv("LEARN_API_CLIENT_SECRET"),
		},
	}
	if cfg.DatabaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is not set")
	}
	if len(cfg.AuthSecret) < 32 {
		return Config{}, errors.New("AUTH_SECRET must be set and at least 32 characters long")
	}
	return cfg, nil
}

func getenv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}