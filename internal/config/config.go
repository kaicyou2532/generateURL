package config

import (
	"fmt"
	"os"
	"strconv"
)

const (
	defaultPort        = 8080
	defaultUploadDir   = "uploads"
	defaultMaxFileSize = int64(10 << 20) // 10 MiB
)

// Config bundles runtime configuration for the HTTP server.
type Config struct {
	Port        int
	UploadDir   string
	MaxFileSize int64
	BaseURL     string
}

// Load constructs a Config using environment variables with sensible defaults.
func Load() (Config, error) {
	cfg := Config{
		Port:        defaultPort,
		UploadDir:   defaultUploadDir,
		MaxFileSize: defaultMaxFileSize,
	}

	if v := os.Getenv("PORT"); v != "" {
		p, err := strconv.Atoi(v)
		if err != nil {
			return Config{}, fmt.Errorf("invalid PORT value: %w", err)
		}
		cfg.Port = p
	}

	if v := os.Getenv("UPLOAD_DIR"); v != "" {
		cfg.UploadDir = v
	}

	if v := os.Getenv("MAX_FILE_SIZE"); v != "" {
		sz, err := strconv.ParseInt(v, 10, 64)
		if err != nil {
			return Config{}, fmt.Errorf("invalid MAX_FILE_SIZE value: %w", err)
		}
		cfg.MaxFileSize = sz
	}

	cfg.BaseURL = os.Getenv("BASE_URL")

	return cfg, nil
}
