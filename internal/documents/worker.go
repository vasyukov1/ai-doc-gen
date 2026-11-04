package documents

import (
	"context"
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
	queue     *pgqueue.Queue
	jobs      *storage.Repository
	owner     string
	sourceDir string
	token     string
}

func NewWorker(queue *pgqueue.Queue, jobs *storage.Repository, owner, sourceDir, token string) *Worker {
	return &Worker{queue: queue, jobs: jobs, owner: owner, sourceDir: sourceDir, token: token}
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
	ctx, cancel := context.WithTimeout(parent, 2*time.Minute)
	defer cancel()

	done := make(chan struct{})
	go w.heartbeat(ctx, claim.ID, done)
	defer close(done)

	input, err := w.jobs.GetInput(ctx, claim.ID)
	if err != nil {
		w.reportFailure(parent, claim, err)
		return
	}

	stage, message := "extracting", "Checking the uploaded archive."
	if input.SourceType == domain.SourceGitHub {
		stage, message = "fetching", "Fetching the repository at a pinned commit."
	}
	if err := w.queue.Stage(ctx, claim.ID, w.owner, stage, stage, message); err != nil {
		log.Printf("update job stage %s: %v", claim.ID, err)
		w.reportFailure(parent, claim, err)
		return
	}

	inventory, err := InspectSource(ctx, claim.ID, claim.ProfileID, w.sourceDir, w.token, input)
	if err != nil {
		w.reportFailure(parent, claim, err)
		return
	}

	if err := w.queue.Stage(ctx, claim.ID, w.owner, "saving_inventory", "inventory", "Saving the file inventory."); err != nil {
		_ = os.RemoveAll(inventory.Source.WorkspacePath)
		log.Printf("update job stage %s: %v", claim.ID, err)
		w.reportFailure(parent, claim, err)
		return
	}
	if err := w.queue.InventoryReady(ctx, claim.ID, w.owner, inventory); err != nil {
		if !errors.Is(err, pgqueue.ErrCommitUnknown) {
			_ = os.RemoveAll(inventory.Source.WorkspacePath)
			w.reportFailure(parent, claim, err)
		}
		log.Printf("save inventory %s: %v", claim.ID, err)
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
	var intakeErr *IntakeError
	if errors.As(err, &intakeErr) && !intakeErr.Temporary {
		if failErr := w.queue.Fail(ctx, claim.ID, w.owner, intakeErr.Code, intakeErr.Message); failErr != nil {
			log.Printf("fail job %s: %v", claim.ID, failErr)
		}
		return
	}

	if claim.Attempts >= 3 {
		if failErr := w.queue.Fail(ctx, claim.ID, w.owner, "source_unavailable", "The source could not be inspected. Submit it again."); failErr != nil {
			log.Printf("fail job %s: %v", claim.ID, failErr)
		}
		return
	}

	delay := time.Duration(claim.Attempts*5) * time.Second
	if retryErr := w.queue.Retry(ctx, claim.ID, w.owner, "source_retry", "Source inspection paused. Retrying shortly.", delay); retryErr != nil {
		log.Printf("retry job %s: %v", claim.ID, retryErr)
	}
	log.Printf("inspect source %s: %v", claim.ID, err)
}

type contextReader struct {
	ctx    context.Context
	reader io.Reader
}

func (r *contextReader) Read(data []byte) (int, error) {
	if err := r.ctx.Err(); err != nil {
		return 0, fmt.Errorf("read source: %w", err)
	}

	return r.reader.Read(data)
}
