package documents

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/pgqueue"
	"github.com/vasyukov1/ai-doc-gen/internal/platform"
	"github.com/vasyukov1/ai-doc-gen/internal/storage"
)

type Worker struct {
	queue         *pgqueue.Queue
	jobs          *storage.Repository
	owner         string
	sourceDir     string
	token         string
	provider      Provider
	encryptionKey string
}

func NewWorker(queue *pgqueue.Queue, jobs *storage.Repository, owner, sourceDir, token string, provider Provider, encryptionKey string) *Worker {
	return &Worker{queue: queue, jobs: jobs, owner: owner, sourceDir: sourceDir, token: token, provider: provider, encryptionKey: encryptionKey}
}

func (w *Worker) Run(ctx context.Context) error {
	for {
		owner, err := platform.NewID()
		if err != nil {
			return err
		}
		w.owner = owner
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
	release, err := w.jobs.GuardJob(parent, claim.ID)
	if err != nil {
		log.Printf("guard job %s: %v", claim.ID, err)
		return
	}
	defer release()
	ctx, cancel := context.WithTimeout(parent, 60*time.Minute)
	defer cancel()

	done := make(chan struct{})
	heartbeatDone := make(chan struct{})
	go func() { defer close(heartbeatDone); w.heartbeat(ctx, claim.ID, done, cancel) }()
	defer func() { close(done); <-heartbeatDone }()

	renderInput, renderErr := w.jobs.DocumentInput(ctx, claim.ID)
	if renderErr == nil && renderInput.Reviewed {
		w.render(ctx, parent, claim, renderInput)
		return
	}
	if renderErr != nil && !errors.Is(renderErr, sql.ErrNoRows) {
		w.reportFailure(parent, claim, renderErr)
		return
	}

	evidence, evidenceErr := w.jobs.GetEvidence(ctx, claim.ID)
	if evidenceErr == nil && evidence.Confirmed {
		w.generate(ctx, parent, claim, evidence)
		return
	}
	if evidenceErr != nil && !errors.Is(evidenceErr, sql.ErrNoRows) {
		w.reportFailure(parent, claim, evidenceErr)
		return
	}

	inventory, sourceErr := w.jobs.GetInventory(ctx, claim.ID)
	if sourceErr == nil {
		w.analyze(ctx, parent, claim, inventory)
		return
	}
	if !errors.Is(sourceErr, sql.ErrNoRows) {
		w.reportFailure(parent, claim, sourceErr)
		return
	}

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

	inventory, err = InspectSource(ctx, claim.ID, claim.ProfileID, w.sourceDir, w.token, input)
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
		return
	}

	w.autoSelect(parent, claim.ID, inventory)
}

// autoSelect confirms the default file selection so intake advances to evidence
// review without a manual step. The inventory stays readable for transparency.
// If confirmation fails, the job waits at selecting_files for a manual choice.
func (w *Worker) autoSelect(ctx context.Context, id string, inventory domain.Inventory) {
	paths := make([]string, 0, len(inventory.Files))
	for _, file := range inventory.Files {
		if file.Selected {
			paths = append(paths, file.Path)
		}
	}
	if len(paths) == 0 {
		return
	}

	confirmCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	if err := w.jobs.ConfirmSelection(confirmCtx, id, paths); err != nil {
		log.Printf("auto-select files %s: %v", id, err)
	}
}

func (w *Worker) heartbeat(ctx context.Context, id string, done <-chan struct{}, cancel context.CancelFunc) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-done:
			return
		case <-ticker.C:
			heartbeatCtx, finish := context.WithTimeout(ctx, 3*time.Second)
			err := w.queue.Heartbeat(heartbeatCtx, id, w.owner)
			finish()
			if err != nil {
				cancel()
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

func (w *Worker) analyze(ctx, parent context.Context, claim *pgqueue.Claim, inventory domain.Inventory) {
	if err := w.queue.Stage(ctx, claim.ID, w.owner, "analyzing", "analyzing", "Extracting source evidence."); err != nil {
		w.reportFailure(parent, claim, err)
		return
	}
	description, err := w.jobs.Description(ctx, claim.ID)
	if err != nil {
		w.reportFailure(parent, claim, err)
		return
	}
	evidence, err := ExtractEvidence(ctx, inventory, description)
	if err != nil {
		w.reportFailure(parent, claim, err)
		return
	}
	if err := w.queue.EvidenceReadyForGeneration(ctx, claim.ID, w.owner, evidence); err != nil {
		w.reportFailure(parent, claim, err)
	}
}
