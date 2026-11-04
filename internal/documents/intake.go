package documents

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

type IntakeError struct {
	Code      string
	Message   string
	Temporary bool
}

func (e *IntakeError) Error() string {
	return e.Message
}

func invalidSource(code, message string) error {
	return &IntakeError{Code: code, Message: message}
}

func InspectSource(ctx context.Context, jobID, profile, sourceDir, token string, input domain.Input) (domain.Inventory, error) {
	if err := os.MkdirAll(sourceDir, 0o700); err != nil {
		return domain.Inventory{}, fmt.Errorf("create source directory: %w", err)
	}

	workDir, err := os.MkdirTemp(sourceDir, jobID+"-")
	if err != nil {
		return domain.Inventory{}, fmt.Errorf("create work directory: %w", err)
	}
	keep := false
	defer func() {
		if !keep {
			_ = os.RemoveAll(workDir)
		}
	}()

	source := domain.Source{JobID: jobID, SourceType: input.SourceType, WorkspacePath: workDir}
	archivePath := input.UploadPath
	if input.SourceType == domain.SourceGitHub {
		archivePath, source.CommitSHA, source.SourceDigest, err = fetchGitHubArchive(ctx, input.GitHubURL, sourceDir, token)
		if err != nil {
			return domain.Inventory{}, err
		}
		defer os.Remove(archivePath)

		source.RepositoryURL = input.GitHubURL
	} else if input.SourceType == domain.SourceZIP {
		if err := checkUpload(ctx, input); err != nil {
			return domain.Inventory{}, err
		}
		source.SourceDigest = input.SHA256
	} else {
		return domain.Inventory{}, invalidSource("unsupported_source", "The source type is not supported.")
	}

	files, err := extractArchive(ctx, archivePath, workDir, profile)
	if err != nil {
		return domain.Inventory{}, err
	}
	if len(files) == 0 {
		return domain.Inventory{}, invalidSource("empty_archive", "No files were found in the project.")
	}

	keep = true

	return domain.Inventory{Source: source, Files: files}, nil
}

func checkUpload(ctx context.Context, input domain.Input) error {
	file, err := os.Open(input.UploadPath)
	if errors.Is(err, os.ErrNotExist) {
		return invalidSource("source_missing", "Uploaded source is unavailable. Submit it again.")
	}
	if err != nil {
		return fmt.Errorf("open upload: %w", err)
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return fmt.Errorf("inspect upload: %w", err)
	}
	if !info.Mode().IsRegular() || info.Size() > MaxArchiveBytes || info.Size() != input.SizeBytes {
		return invalidSource("upload_changed", "The uploaded archive changed. Submit it again.")
	}

	hash := sha256.New()
	if _, err := io.Copy(hash, &contextReader{ctx: ctx, reader: file}); err != nil {
		return fmt.Errorf("read upload: %w", err)
	}
	if hex.EncodeToString(hash.Sum(nil)) != input.SHA256 {
		return invalidSource("upload_changed", "The uploaded archive changed. Submit it again.")
	}

	return nil
}
