package config

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	HTTPAddr, DatabaseURL, DataDir, UploadDir, SourceDir, StateDir, GitHubToken, LLMEncryptionKey string
	Retention                                                                                     time.Duration
	CookieSecure                                                                                  bool
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
	stateDir := os.Getenv("STATE_DIR")
	if stateDir == "" {
		stateDir = filepath.Join(root, "state")
	}
	stateDir, err = filepath.Abs(stateDir)
	if err != nil {
		return Config{}, err
	}
	for _, dir := range []string{uploadDir, sourceDir, stateDir} {
		rel, err := filepath.Rel(root, dir)
		if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return Config{}, errors.New("UPLOAD_DIR, SOURCE_DIR and STATE_DIR must be inside DATA_DIR")
		}
	}

	cookieSecure := false
	if value := os.Getenv("COOKIE_SECURE"); value != "" {
		cookieSecure, err = strconv.ParseBool(value)
		if err != nil {
			return Config{}, errors.New("COOKIE_SECURE must be a boolean")
		}
	}

	return Config{CookieSecure: cookieSecure, HTTPAddr: addr, DatabaseURL: databaseURL, DataDir: root, UploadDir: uploadDir, SourceDir: sourceDir, StateDir: stateDir, GitHubToken: os.Getenv("GITHUB_TOKEN"), LLMEncryptionKey: os.Getenv("LLM_ENCRYPTION_KEY"), Retention: time.Duration(days) * 24 * time.Hour}, nil
}

// EncryptionKey returns the AES master key for saved provider credentials.
// An explicit LLM_ENCRYPTION_KEY wins. Otherwise a key is generated once and
// stored in the state directory, so local runs work without manual setup. The
// web and worker processes must share this directory.
func (c Config) EncryptionKey() (string, error) {
	if c.LLMEncryptionKey != "" {
		return c.LLMEncryptionKey, nil
	}

	path := filepath.Join(c.StateDir, "llm-encryption.key")
	key, err := os.ReadFile(path)
	if err == nil {
		return strings.TrimSpace(string(key)), nil
	}
	if !errors.Is(err, os.ErrNotExist) {
		return "", fmt.Errorf("read encryption key: %w", err)
	}
	if err := os.MkdirAll(c.StateDir, 0o700); err != nil {
		return "", fmt.Errorf("create state directory: %w", err)
	}

	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", fmt.Errorf("generate encryption key: %w", err)
	}
	generated := base64.StdEncoding.EncodeToString(raw)
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if errors.Is(err, os.ErrExist) {
		key, readErr := os.ReadFile(path)
		if readErr != nil {
			return "", fmt.Errorf("read encryption key: %w", readErr)
		}
		return strings.TrimSpace(string(key)), nil
	}
	if err != nil {
		return "", fmt.Errorf("create encryption key: %w", err)
	}
	if _, err := file.WriteString(generated + "\n"); err != nil {
		_ = file.Close()
		return "", fmt.Errorf("write encryption key: %w", err)
	}
	if err := file.Close(); err != nil {
		return "", fmt.Errorf("close encryption key: %w", err)
	}

	return generated, nil
}
