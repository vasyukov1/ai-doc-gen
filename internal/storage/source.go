package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/go-jet/jet/v2/postgres"
	"github.com/vasyukov1/ai-doc-gen/internal/domain"
	"github.com/vasyukov1/ai-doc-gen/internal/storage/generated/ai_doc_gen/public/table"
)

var ErrSelectionUnavailable = errors.New("file selection is not available")
var ErrInvalidSelection = errors.New("invalid file selection")

func SaveInventoryTx(ctx context.Context, tx *sql.Tx, inventory domain.Inventory) error {
	source := inventory.Source
	sources := table.JobSources
	stmt := sources.INSERT(sources.JobID, sources.SourceType, sources.RepositoryURL,
		sources.CommitSha, sources.SourceDigest, sources.WorkspacePath).
		VALUES(source.JobID, source.SourceType, source.RepositoryURL, source.CommitSHA,
			source.SourceDigest, source.WorkspacePath)
	query, args := stmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("save source: %w", err)
	}

	files := table.JobFiles
	for start := 0; start < len(inventory.Files); start += 200 {
		end := min(start+200, len(inventory.Files))
		insert := files.INSERT(files.JobID, files.Path, files.SizeBytes, files.Sha256,
			files.ReasonCode, files.Eligible, files.Selected)

		for _, file := range inventory.Files[start:end] {
			insert = insert.VALUES(source.JobID, file.Path, file.SizeBytes, file.SHA256,
				file.ReasonCode, file.Eligible, file.Selected)
		}

		query, args := insert.Sql()
		if _, err := tx.ExecContext(ctx, query, args...); err != nil {
			return fmt.Errorf("save source files: %w", err)
		}
	}

	return nil
}

func (r *Repository) GetInventory(ctx context.Context, id string) (domain.Inventory, error) {
	sources := table.JobSources
	stmt := sources.SELECT(sources.SourceType, sources.RepositoryURL, sources.CommitSha,
		sources.SourceDigest, sources.WorkspacePath).
		WHERE(sources.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args := stmt.Sql()

	inventory := domain.Inventory{Source: domain.Source{JobID: id}}
	err := r.db.QueryRowContext(ctx, query, args...).Scan(&inventory.Source.SourceType,
		&inventory.Source.RepositoryURL, &inventory.Source.CommitSHA,
		&inventory.Source.SourceDigest, &inventory.Source.WorkspacePath)
	if err != nil {
		return domain.Inventory{}, fmt.Errorf("get source: %w", err)
	}

	files := table.JobFiles
	fileStmt := files.SELECT(files.Path, files.SizeBytes, files.Sha256,
		files.ReasonCode, files.Eligible, files.Selected).
		WHERE(files.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).
		ORDER_BY(files.Path)
	query, args = fileStmt.Sql()
	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return domain.Inventory{}, fmt.Errorf("list source files: %w", err)
	}
	defer rows.Close()

	inventory.Files = make([]domain.SourceFile, 0)
	for rows.Next() {
		var file domain.SourceFile
		if err := rows.Scan(&file.Path, &file.SizeBytes, &file.SHA256,
			&file.ReasonCode, &file.Eligible, &file.Selected); err != nil {
			return domain.Inventory{}, fmt.Errorf("read source file: %w", err)
		}
		inventory.Files = append(inventory.Files, file)
	}
	if err := rows.Err(); err != nil {
		return domain.Inventory{}, fmt.Errorf("list source files: %w", err)
	}

	return inventory, nil
}

func (r *Repository) ConfirmSelection(ctx context.Context, id string, paths []string) error {
	if len(paths) == 0 || len(paths) > 500 {
		return ErrInvalidSelection
	}

	selected := make(map[string]bool, len(paths))
	for _, path := range paths {
		if selected[path] {
			return ErrInvalidSelection
		}
		selected[path] = true
	}

	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("start selection: %w", err)
	}
	defer tx.Rollback()

	jobs := table.Jobs
	stmt := jobs.SELECT(jobs.Status, jobs.Stage).
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID())).
		FOR(postgres.UPDATE())
	query, args := stmt.Sql()
	var status domain.Status
	var stage string
	if err := tx.QueryRowContext(ctx, query, args...).Scan(&status, &stage); err != nil {
		return fmt.Errorf("get selection state: %w", err)
	}
	if status != domain.StatusNeedsInput || stage != "selecting_files" {
		return ErrSelectionUnavailable
	}

	files := table.JobFiles
	fileStmt := files.SELECT(files.Path).
		WHERE(files.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()).
			AND(files.Eligible.EQ(postgres.Bool(true))))
	query, args = fileStmt.Sql()
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("list eligible files: %w", err)
	}

	eligible := make(map[string]bool)
	for rows.Next() {
		var path string
		if err := rows.Scan(&path); err != nil {
			rows.Close()
			return fmt.Errorf("read eligible file: %w", err)
		}
		eligible[path] = true
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return fmt.Errorf("list eligible files: %w", err)
	}
	rows.Close()

	for path := range selected {
		if !eligible[path] {
			return ErrInvalidSelection
		}
	}

	clearStmt := files.UPDATE(files.Selected).SET(false).
		WHERE(files.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args = clearStmt.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("clear file selection: %w", err)
	}

	pathValues := make([]postgres.Expression, 0, len(paths))
	for _, path := range paths {
		pathValues = append(pathValues, postgres.Text(path))
	}
	selectStmt := files.UPDATE(files.Selected).SET(true).
		WHERE(files.JobID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()).
			AND(files.Path.IN(pathValues...)))
	query, args = selectStmt.Sql()
	result, err := tx.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("save file selection: %w", err)
	}
	count, err := result.RowsAffected()
	if err != nil || count != int64(len(paths)) {
		return ErrInvalidSelection
	}

	if err := TransitionTx(ctx, tx, id, domain.StatusNeedsInput, domain.StatusInputReady,
		"source_selected", "files_selected", "File selection saved. Evidence extraction is the next step."); err != nil {
		return err
	}
	update := jobs.UPDATE(jobs.RequiredAction).SET("").
		WHERE(jobs.ID.EQ(postgres.CAST(postgres.Text(id)).AS_UUID()))
	query, args = update.Sql()
	if _, err := tx.ExecContext(ctx, query, args...); err != nil {
		return fmt.Errorf("clear required action: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit selection: %w", err)
	}

	return nil
}
