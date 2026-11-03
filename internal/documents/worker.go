package documents

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/pgqueue"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
)

type Worker struct {
	queue *pgqueue.Queue
	jobs  *storage.Repository
	owner string
}

func NewWorker(queue *pgqueue.Queue, jobs *storage.Repository, owner string) *Worker {
	return &Worker{queue: queue, jobs: jobs, owner: owner}
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		claim, err := w.queue.Claim(ctx, w.owner)
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}

			log.Printf("claim job: %v", err)
		} else if claim != nil {
			w.process(ctx, claim)
			continue
		}

		select {
		case <-ctx.Done():
			return nil
		case <-time.After(time.Second):
		}
	}
}

func (w *Worker) process(parent context.Context, claim *pgqueue.Claim) {
	ctx, cancel := context.WithTimeout(parent, 20*time.Second)
	defer cancel()

	done := make(chan struct{})
	go w.heartbeat(ctx, claim.ID, done)
	defer close(done)

	input, err := w.jobs.GetInput(ctx, claim.ID)
	if err == nil {
		err = validateInput(ctx, input)
	}

	if err != nil {
		w.reportFailure(parent, claim, err)
		return
	}

	if err := w.queue.InputReady(parent, claim.ID, w.owner); err != nil {
		log.Printf("finish job %s: %v", claim.ID, err)
	}
}

func (w *Worker) heartbeat(ctx context.Context, id string, done <-chan struct{}) {
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-ticker.C:
			if err := w.queue.Heartbeat(ctx, id, w.owner); err != nil {
				log.Printf("heartbeat job %s: %v", id, err)
				return
			}
		}
	}
}

func (w *Worker) reportFailure(ctx context.Context, claim *pgqueue.Claim, err error) {
	if errors.Is(err, os.ErrNotExist) {
		if failErr := w.queue.Fail(ctx, claim.ID, w.owner, "source_missing", "Uploaded source is unavailable. Submit it again."); failErr != nil {
			log.Printf("fail job %s: %v", claim.ID, failErr)
		}
		return
	}

	if claim.Attempts >= 3 {
		if failErr := w.queue.Fail(ctx, claim.ID, w.owner, "input_error", "The input could not be checked. Submit it again."); failErr != nil {
			log.Printf("fail job %s: %v", claim.ID, failErr)
		}
		return
	}

	delay := time.Duration(claim.Attempts*5) * time.Second
	if retryErr := w.queue.Retry(ctx, claim.ID, w.owner, "input_retry", "Input check paused. Retrying shortly.", delay); retryErr != nil {
		log.Printf("retry job %s: %v", claim.ID, retryErr)
	}
}

func validateInput(ctx context.Context, input domain.Input) error {
	if input.SourceType == domain.SourceGitHub {
		if input.GitHubURL == "" {
			return errors.New("GitHub URL is empty")
		}
		return nil
	}
	if input.SourceType != domain.SourceZIP {
		return errors.New("unsupported input type")
	}

	file, err := os.Open(input.UploadPath)
	if err != nil {
		return fmt.Errorf("open upload: %w", err)
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, &contextReader{ctx: ctx, reader: file}); err != nil {
		return fmt.Errorf("read upload: %w", err)
	}
	if hex.EncodeToString(hash.Sum(nil)) != input.SHA256 {
		return errors.New("upload digest changed")
	}

	return nil
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(data []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, err
	}

	return r.reader.Read(data)
}
