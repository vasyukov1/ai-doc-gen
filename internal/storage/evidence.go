package storage

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

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
	query, args := evidenceTable.INSERT(evidenceTable.JobID, evidenceTable.Record).VALUES(id, string(payload)).
		ON_CONFLICT(evidenceTable.JobID).DO_UPDATE(postgres.SET(evidenceTable.Record.SET(evidenceTable.EXCLUDED.Record))).Sql()
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
