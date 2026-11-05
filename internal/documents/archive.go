package documents

import (
	"archive/zip"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
)

const MaxArchiveBytes = 100 << 20
const maxExpandedBytes = 500 << 20
const maxTextFileBytes = 2 << 20
const maxArchiveEntries = 20000
const maxPathBytes = 512
const maxPathDepth = 20

func extractArchive(ctx context.Context, archivePath, workDir, profile string) ([]domain.SourceFile, error) {
	archive, err := zip.OpenReader(archivePath)
	if err != nil {
		return nil, invalidSource("invalid_archive", "The ZIP archive could not be opened.")
	}
	defer archive.Close()

	entries, prefix, err := inspectEntries(archive.File)
	if err != nil {
		return nil, err
	}

	files := make([]domain.SourceFile, 0, len(entries))
	seen := make(map[string]bool, len(entries))
	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return nil, err
		}

		name := strings.TrimPrefix(entry.Name, prefix)
		if entry.FileInfo().IsDir() || name == "" {
			continue
		}
		key := strings.ToLower(name)
		if entry.Mode()&os.ModeSymlink != 0 || !entry.Mode().IsRegular() {
			files = append(files, domain.SourceFile{
				Path: name, SizeBytes: int64(entry.UncompressedSize64),
				ReasonCode: "unsafe_entry", Eligible: false, Selected: false,
			})
			continue
		}
		if seen[key] {
			files = append(files, domain.SourceFile{
				Path: name, SizeBytes: int64(entry.UncompressedSize64),
				ReasonCode: "duplicate_path", Eligible: false, Selected: false,
			})
			continue
		}
		seen[key] = true

		reason, eligible, selected := classifyFile(name, profile, entry.UncompressedSize64)
		file := domain.SourceFile{
			Path: name, SizeBytes: int64(entry.UncompressedSize64),
			ReasonCode: reason, Eligible: eligible, Selected: selected,
		}
		if eligible {
			digest, binary, err := extractTextFile(ctx, entry, filepath.Join(workDir, filepath.FromSlash(name)))
			if err != nil {
				return nil, err
			}
			if binary {
				file.ReasonCode = "binary_content"
				file.Eligible = false
				file.Selected = false
			} else {
				file.SHA256 = digest
			}
		}

		files = append(files, file)
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Path < files[j].Path })
	selected := 0
	for index := range files {
		if !files[index].Selected {
			continue
		}
		selected++
		if selected > 500 {
			files[index].Selected = false
		}
	}

	return files, nil
}

func inspectEntries(entries []*zip.File) ([]*zip.File, string, error) {
	if len(entries) > maxArchiveEntries {
		return nil, "", invalidSource("too_many_files", "The archive has too many entries.")
	}

	seen := make(map[string]bool, len(entries))
	regular := make(map[string]bool, len(entries))
	var expanded uint64
	var firstRoot string
	sharedRoot := true

	for _, entry := range entries {
		name := strings.TrimSuffix(entry.Name, "/")
		if err := validateEntryPath(name); err != nil {
			return nil, "", err
		}
		// Links, special files and duplicate names are excluded from the
		// inventory instead of rejecting the whole project.
		if entry.Mode()&os.ModeSymlink != 0 || (!entry.FileInfo().IsDir() && !entry.Mode().IsRegular()) {
			continue
		}
		seen[strings.ToLower(name)] = true

		if entry.FileInfo().IsDir() {
			continue
		}
		regular[strings.ToLower(name)] = true
		if entry.UncompressedSize64 > maxExpandedBytes-expanded {
			return nil, "", invalidSource("expanded_size", "The archive expands beyond the size limit.")
		}
		expanded += entry.UncompressedSize64

		root, rest, hasRoot := strings.Cut(name, "/")
		if !hasRoot || rest == "" {
			sharedRoot = false
		}
		if firstRoot == "" {
			firstRoot = root
		} else if firstRoot != root {
			sharedRoot = false
		}
	}
	for name := range seen {
		parent := path.Dir(name)
		for parent != "." {
			if regular[parent] {
				return nil, "", invalidSource("unsafe_path", "The archive contains conflicting paths at "+clipText(name, 200)+".")
			}
			parent = path.Dir(parent)
		}
	}

	if sharedRoot && firstRoot != "" {
		return entries, firstRoot + "/", nil
	}

	return entries, "", nil
}

func validateEntryPath(name string) error {
	if name == "" || name == "." || len(name) > maxPathBytes || !utf8.ValidString(name) ||
		strings.ContainsAny(name, "\\\x00") || strings.HasPrefix(name, "/") ||
		path.Clean(name) != name || strings.HasPrefix(name, "../") || name == ".." ||
		strings.Contains(name, ":") || strings.Count(name, "/") >= maxPathDepth {
		return invalidSource("unsafe_path", "The archive contains an unsafe path: "+clipText(name, 200)+".")
	}

	return nil
}

// ValidateArchive checks the central directory without unpacking files, so an
// upload can be rejected with the exact reason before a job is created.
func ValidateArchive(path string) error {
	archive, err := zip.OpenReader(path)
	if err != nil {
		return invalidSource("invalid_archive", "The ZIP archive could not be opened.")
	}
	defer archive.Close()

	_, _, err = inspectEntries(archive.File)

	return err
}

func extractTextFile(ctx context.Context, entry *zip.File, destination string) (string, bool, error) {
	if err := os.MkdirAll(filepath.Dir(destination), 0o700); err != nil {
		return "", false, fmt.Errorf("create source directory: %w", err)
	}

	source, err := entry.Open()
	if err != nil {
		return "", false, invalidSource("invalid_archive", "The ZIP archive contains an unreadable file.")
	}
	defer source.Close()

	file, err := os.OpenFile(destination, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o600)
	if err != nil {
		return "", false, fmt.Errorf("create source file: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	size, err := io.Copy(io.MultiWriter(file, hash), io.LimitReader(&contextReader{ctx: ctx, reader: source}, maxTextFileBytes+1))
	if err != nil {
		return "", false, invalidSource("invalid_archive", "The ZIP archive contains an unreadable file.")
	}
	if size > maxTextFileBytes || uint64(size) != entry.UncompressedSize64 {
		return "", false, invalidSource("expanded_size", "A file exceeds the archive size limit.")
	}
	if err := file.Close(); err != nil {
		return "", false, fmt.Errorf("close source file: %w", err)
	}

	binary, err := hasBinaryContent(destination)
	if err != nil {
		return "", false, fmt.Errorf("inspect source file: %w", err)
	}
	if binary {
		if err := os.Remove(destination); err != nil {
			return "", false, fmt.Errorf("remove binary file: %w", err)
		}
		return "", true, nil
	}

	return hex.EncodeToString(hash.Sum(nil)), false, nil
}

func hasBinaryContent(name string) (bool, error) {
	data, err := os.ReadFile(name)
	if err != nil {
		return false, err
	}

	return bytes.IndexByte(data, 0) >= 0 || !utf8.Valid(data), nil
}
