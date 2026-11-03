package config

import (
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	HTTPAddr    string
	DatabaseURL string
	UploadDir   string
}

func Load() (Config, error) {
	addr := os.Getenv("HTTP_ADDR")
	if addr == "" {
		addr = ":8080"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is required")
	}

	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "./artifacts/uploads"
	}
	uploadDir, err := filepath.Abs(uploadDir)
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTPAddr:    addr,
		DatabaseURL: databaseURL,
		UploadDir:   uploadDir,
	}, nil
}
