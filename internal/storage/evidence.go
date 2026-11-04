package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

var ErrEvidenceUnavailable = errors.New("evidence review is unavailable")
var ErrInvalidEvidence = errors.New("invalid evidence correction")

func (r *Repository) Description(ctx context.Context, id string) (string, error) {
	jobs := table.Jobs
	query, args := jobs.SELECT(jobs.Description).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	var description string
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&description)

	return description, err
}

func SaveEvidenceTx(ctx context.Context, tx *sql.Tx, id string, evidence domain.Evidence) error {
	payload, err := json.Marshal(evidence)
	if err != nil {
		return err
	}
	evidenceTable := table.JobEvidence
	query, args := evidenceTable.INSERT(evidenceTable.JobID, evidenceTable.Record).VALUES(id, string(payload)).ON_CONFLICT(evidenceTable.JobID).DO_NOTHING().Sql()
	_, err = tx.ExecContext(ctx, query, args...)

	return err
}

func (r *Repository) GetEvidence(ctx context.Context, id string) (domain.Evidence, error) {
	evidenceTable := table.JobEvidence
	query, args := evidenceTable.SELECT(evidenceTable.Record).WHERE(evidenceTable.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	var payload []byte
	if err := r.db.QueryRowContext(ctx, query, args...).Scan(&payload); err != nil {
		return domain.Evidence{}, err
	}
	var evidence domain.Evidence
	err := json.Unmarshal(payload, &evidence)

	return evidence, err
}

func (r *Repository) ConfirmEvidence(ctx context.Context, id, digest string, corrections map[string]string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	jobs := table.Jobs
	query, args := jobs.SELECT(jobs.Status, jobs.Stage).WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).FOR(postgres.UPDATE()).Sql()
	var status domain.Status
	var stage string
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&status, &stage); err != nil {
		return err
	}
	if status != domain.StatusNeedsInput || stage != "reviewing_evidence" {
		return ErrEvidenceUnavailable
	}
	evidenceTable := table.JobEvidence
	query, args = evidenceTable.SELECT(evidenceTable.Record).WHERE(evidenceTable.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	var payload []byte
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&payload); err != nil {
		return err
	}
	var evidence domain.Evidence
	if err := json.Unmarshal(payload, &evidence); err != nil {
		return err
	}
	if digest != evidence.Digest {
		return ErrEvidenceUnavailable
	}
	valid := map[string]bool{}
	for _, fact := range evidence.Facts {
		valid[fact.ID] = true
	}
	for key, value := range corrections {
		if (key != "additional_context" && !valid[key]) || len(value) > 4000 {
			return ErrInvalidEvidence
		}
	}
	evidence.Corrections = corrections
	evidence.Confirmed = true
	payload, err = json.Marshal(evidence)
	if err != nil {
		return err
	}
	query, args = evidenceTable.UPDATE(evidenceTable.Record, evidenceTable.ConfirmedAt).SET(string(payload), postgres.NOW()).WHERE(evidenceTable.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	if err := TransitionTx(ctx, tx, id, status, domain.StatusInputReady, "evidence_confirmed", "evidence_confirmed", "Evidence confirmed. Document drafting is not available yet."); err != nil {
		return err
	}
	query, args = jobs.UPDATE(jobs.RequiredAction).SET("").WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit evidence review: %w", err)
	}

	return nil
}
