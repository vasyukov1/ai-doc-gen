package documents

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/platform"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
)

type Cleaner struct {
	jobs                   *storage.Repository
	root, uploads, sources string
}

func NewCleaner(jobs *storage.Repository, root, uploads, sources string) *Cleaner {
	return &Cleaner{jobs, root, uploads, sources}
}
func (c *Cleaner) Run(ctx context.Context) {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()
	for {
		if err := c.Once(ctx); err != nil && ctx.Err() == nil {
			log.Printf("cleanup: %v", err)
		}
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
		}
	}
}
func (c *Cleaner) Once(ctx context.Context) error {
	if err := c.jobs.ExpireJobs(ctx); err != nil {
		return err
	}
	if err := c.jobs.PurgeJobs(ctx, c.removeJob); err != nil {
		return err
	}
	if err := c.jobs.CleanupCache(ctx); err != nil {
		return err
	}
	return c.orphans(ctx)
}

// Never follow a symlink in a deletion path, even below a configured root.
func (c *Cleaner) safeRemove(path string) error {
	root, err := filepath.Abs(c.root)
	if err != nil {
		return err
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return err
	}
	rel, err := filepath.Rel(root, path)
	if err != nil || rel == "." || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return fmt.Errorf("cleanup path outside data root")
	}
	current := root
	info, err := os.Lstat(root)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return fmt.Errorf("data root is a symlink")
	}
	parts := strings.Split(rel, string(filepath.Separator))
	for i, part := range parts {
		current = filepath.Join(current, part)
		info, err = os.Lstat(current)
		if os.IsNotExist(err) {
			return nil
		}
		if err != nil {
			return err
		}
		if info.Mode()&os.ModeSymlink != 0 && i < len(parts)-1 {
			return fmt.Errorf("cleanup parent is a symlink")
		}
	}
	return os.RemoveAll(path)
}
func (c *Cleaner) removeJob(id string) error {
	if !platform.ValidID(id) {
		return fmt.Errorf("invalid cleanup ID")
	}
	for _, dir := range []string{c.uploads, c.sources} {
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), id+"-") {
				if err = c.safeRemove(filepath.Join(dir, entry.Name())); err != nil {
					return err
				}
			}
		}
	}
	return c.safeRemove(filepath.Join(c.sources, "artifacts", id))
}

// Sweep only old unreferenced entries. Every live or legacy job protects its work.
func (c *Cleaner) orphans(ctx context.Context) error {
	paths, active, err := c.jobs.ProtectedPaths(ctx)
	if err != nil {
		return err
	}
	protected := func(path string) bool {
		for saved := range paths {
			if saved == path || strings.HasPrefix(saved, path+string(filepath.Separator)) {
				return true
			}
		}
		return false
	}
	removed := 0
	for _, dir := range []string{c.uploads, c.sources, filepath.Join(c.sources, "artifacts")} {
		entries, err := os.ReadDir(dir)
		if os.IsNotExist(err) {
			continue
		}
		if err != nil {
			return err
		}
		for _, entry := range entries {
			if removed >= 50 {
				return nil
			}
			if ctx.Err() != nil {
				return ctx.Err()
			}
			name := entry.Name()
			path := filepath.Join(dir, name)
			if name == "artifacts" || protected(path) {
				continue
			}
			id := name
			if len(id) > 36 {
				id = id[:36]
			}
			if active[id] {
				continue
			}
			info, err := entry.Info()
			if err != nil {
				return err
			}
			if time.Since(info.ModTime()) < 24*time.Hour {
				continue
			}
			// Keep known artifact versions and sweep superseded versions separately.
			if err = c.safeRemove(path); err != nil {
				return err
			}
			removed++
		}
	}
	// Artifact roots with current publications still need their superseded renders removed.
	roots, _ := os.ReadDir(filepath.Join(c.sources, "artifacts"))
	for _, job := range roots {
		if removed >= 50 {
			return nil
		}
		if !platform.ValidID(job.Name()) || active[job.Name()] {
			continue
		}
		dir := filepath.Join(c.sources, "artifacts", job.Name())
		versions, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, version := range versions {
			path := filepath.Join(dir, version.Name())
			if protected(path) {
				continue
			}
			info, err := version.Info()
			if err != nil {
				return err
			}
			if time.Since(info.ModTime()) < 24*time.Hour {
				continue
			}
			if err = c.safeRemove(path); err != nil {
				return err
			}
			removed++
			if removed >= 50 {
				return nil
			}
		}
	}
	return nil
}
