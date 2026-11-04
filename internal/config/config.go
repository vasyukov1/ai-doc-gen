package config

import (
	"errors"
	"os"
	"path/filepath"
)

type Config struct {
	HTTPAddr         string
	DatabaseURL      string
	UploadDir        string
	SourceDir        string
	GitHubToken      string
	LLMEncryptionKey string
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

	sourceDir := os.Getenv("SOURCE_DIR")
	if sourceDir == "" {
		sourceDir = "./artifacts/sources"
	}
	sourceDir, err = filepath.Abs(sourceDir)
	if err != nil {
		return Config{}, err
	}

	return Config{
		HTTPAddr:         addr,
		DatabaseURL:      databaseURL,
		UploadDir:        uploadDir,
		SourceDir:        sourceDir,
		GitHubToken:      os.Getenv("GITHUB_TOKEN"),
		LLMEncryptionKey: os.Getenv("LLM_ENCRYPTION_KEY"),
	}, nil
}
