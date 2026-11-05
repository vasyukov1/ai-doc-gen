package documents

import (
	"archive/zip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/platform"
)

// CopyRetainedSource creates independently disposable source and input files.
// Only retained text with the recorded digest is copied; no source is executed.
func CopyRetainedSource(ctx context.Context, original domain.Inventory, id, profile, uploadDir string) (domain.Inventory, domain.Input, error) {
	if !platform.ValidID(id) || original.Source.WorkspacePath == "" || len(original.Files) > maxArchiveEntries {
		return domain.Inventory{}, domain.Input{}, errors.New("saved source is unavailable")
	}
	info, err := os.Lstat(original.Source.WorkspacePath)
	if err != nil || !info.IsDir() || info.Mode()&os.ModeSymlink != 0 {
		return domain.Inventory{}, domain.Input{}, errors.New("saved source is unavailable")
	}
	root, err := os.OpenRoot(original.Source.WorkspacePath)
	if err != nil {
		return domain.Inventory{}, domain.Input{}, fmt.Errorf("open saved source: %w", err)
	}
	defer root.Close()

	workspace, err := os.MkdirTemp(filepath.Dir(original.Source.WorkspacePath), id+"-")
	if err != nil {
		return domain.Inventory{}, domain.Input{}, fmt.Errorf("create reused workspace: %w", err)
	}
	if err := os.MkdirAll(uploadDir, 0o700); err != nil {
		_ = os.RemoveAll(workspace)
		return domain.Inventory{}, domain.Input{}, fmt.Errorf("create upload directory: %w", err)
	}
	upload, err := os.CreateTemp(uploadDir, id+"-*.zip")
	if err != nil {
		_ = os.RemoveAll(workspace)
		return domain.Inventory{}, domain.Input{}, fmt.Errorf("create reused archive: %w", err)
	}
	keep := false
	defer func() {
		_ = upload.Close()
		if !keep {
			_ = os.RemoveAll(workspace)
			_ = os.Remove(upload.Name())
		}
	}()

	archive := zip.NewWriter(&boundedArchiveWriter{writer: upload, remaining: MaxArchiveBytes})
	inventory := domain.Inventory{Source: original.Source, Files: append([]domain.SourceFile(nil), original.Files...)}
	inventory.Source.JobID = id
	inventory.Source.WorkspacePath = workspace
	seen := map[string]bool{}
	var copied, selected int
	var expanded int64
	for index := range inventory.Files {
		file := &inventory.Files[index]
		if !file.Eligible {
			file.Selected = false
			continue
		}
		if err := validateEntryPath(file.Path); err != nil {
			return domain.Inventory{}, domain.Input{}, err
		}
		if seen[file.Path] || file.SizeBytes < 0 || file.SizeBytes > maxTextFileBytes || expanded+file.SizeBytes > maxExpandedBytes {
			return domain.Inventory{}, domain.Input{}, errors.New("saved source exceeds safe copy limits")
		}
		seen[file.Path] = true
		file.ReasonCode, file.Eligible, file.Selected = classifyFile(file.Path, profile, uint64(file.SizeBytes))
		if !file.Eligible {
			continue
		}
		if err := copyRetainedFile(ctx, root, workspace, archive, *file); err != nil {
			return domain.Inventory{}, domain.Input{}, err
		}
		expanded += file.SizeBytes
		copied++
		if file.Selected {
			selected++
			if selected > 500 {
				file.Selected = false
			}
		}
	}
	if copied == 0 || selected == 0 {
		return domain.Inventory{}, domain.Input{}, errors.New("no retained text is available for this profile")
	}
	if err := archive.Close(); err != nil {
		return domain.Inventory{}, domain.Input{}, fmt.Errorf("finish reused archive: %w", err)
	}
	if _, err := upload.Seek(0, io.SeekStart); err != nil {
		return domain.Inventory{}, domain.Input{}, err
	}
	hash := sha256.New()
	size, err := io.Copy(hash, &contextReader{ctx: ctx, reader: upload})
	if err != nil {
		return domain.Inventory{}, domain.Input{}, err
	}
	if err := upload.Close(); err != nil {
		return domain.Inventory{}, domain.Input{}, err
	}
	keep = true

	return inventory, domain.Input{SourceType: domain.SourceZIP, UploadPath: upload.Name(), OriginalName: "retained-project.zip", SHA256: hex.EncodeToString(hash.Sum(nil)), SizeBytes: size}, nil
}

func copyRetainedFile(ctx context.Context, root *os.Root, workspace string, archive *zip.Writer, file domain.SourceFile) error {
	info, err := root.Lstat(filepath.FromSlash(file.Path))
	if err != nil || !info.Mode().IsRegular() || info.Size() != file.SizeBytes {
		return errors.New("saved source file is missing or changed")
	}
	source, err := root.Open(filepath.FromSlash(file.Path))
	if err != nil {
		return fmt.Errorf("open retained file: %w", err)
	}
	defer source.Close()

	destination := filepath.Join(workspace, filepath.FromSlash(file.Path))
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return err
	}
	output, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return err
	}
	defer output.Close()
	entry, err := archive.Create("project/" + file.Path)
	if err != nil {
		return err
	}

	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(output, entry, hash), io.LimitReader(&contextReader{ctx: ctx, reader: source}, maxTextFileBytes+1))
	if err != nil {
		return fmt.Errorf("copy retained file: %w", err)
	}
	if size != file.SizeBytes || hex.EncodeToString(hash.Sum(nil)) != file.SHA256 {
		return errors.New("saved source file is missing or changed")
	}
	if binary, err := hasBinaryContent(destination); err != nil || binary {
		return errors.New("saved source file is not retained UTF-8 text")
	}

	return output.Close()
}

type boundedArchiveWriter struct {
	writer    io.Writer
	remaining int64
}

func (w *boundedArchiveWriter) Write(data []byte) (int, error) {
	if int64(len(data)) > w.remaining {
		return 0, errors.New("retained archive exceeds the upload limit")
	}
	written, err := w.writer.Write(data)
	w.remaining -= int64(written)

	return written, err
}
