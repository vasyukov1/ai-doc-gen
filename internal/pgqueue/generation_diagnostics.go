package pgqueue

import (
	"context"
	"encoding/json"
	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

// SaveCallDiagnostics retains private repair data under the job's existing retention policy.
func (q *Queue) SaveCallDiagnostics(ctx context.Context, id, owner string, callID int64, diagnostic any, rejected string) error {
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = checkOwner(ctx, tx, id, owner); err != nil {
		return err
	}
	raw, err := json.Marshal(diagnostic)
	if err != nil {
		return err
	}
	t := table.ProviderCalls
	query, args := t.UPDATE(t.Diagnostics, t.RejectedContent).SET(string(raw), rejected).WHERE(t.ID.EQ(postgres.Int64(callID)).AND(t.JobID.EQ(draftID(id)))).Sql()
	_, err = tx.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}
	return tx.Commit()
}
