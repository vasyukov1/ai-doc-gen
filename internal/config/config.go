package config

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr, DatabaseURL, DataDir, UploadDir, SourceDir, GitHubToken, LLMEncryptionKey string
	Retention                                                                           time.Duration
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
	root := os.Getenv("DATA_DIR")
	if root == "" {
		root = "./artifacts"
	}
	root, err := filepath.Abs(root)
	if err != nil {
		return Config{}, err
	}
	days := 7
	if value := os.Getenv("RETENTION_DAYS"); value != "" {
		days, err = strconv.Atoi(value)
		if err != nil || days < 1 || days > 365 {
			return Config{}, errors.New("RETENTION_DAYS must be 1 through 365")
		}
	}
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = filepath.Join(root, "uploads")
	}
	uploadDir, err = filepath.Abs(uploadDir)
	if err != nil {
		return Config{}, err
	}
	sourceDir := os.Getenv("SOURCE_DIR")
	if sourceDir == "" {
		sourceDir = filepath.Join(root, "sources")
	}
	sourceDir, err = filepath.Abs(sourceDir)
	if err != nil {
		return Config{}, err
	}
	for _, dir := range []string{uploadDir, sourceDir} {
		rel, err := filepath.Rel(root, dir)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return Config{}, errors.New("UPLOAD_DIR and SOURCE_DIR must be inside DATA_DIR")
		}
	}
	return Config{HTTPAddr: addr, DatabaseURL: databaseURL, DataDir: root, UploadDir: uploadDir, SourceDir: sourceDir, GitHubToken: os.Getenv("GITHUB_TOKEN"), LLMEncryptionKey: os.Getenv("LLM_ENCRYPTION_KEY"), Retention: time.Duration(days) * 24 * time.Hour}, nil
}
