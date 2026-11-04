-- +goose Up
CREATE TABLE job_sources (
    job_id          uuid PRIMARY KEY REFERENCES jobs(id) ON DELETE CASCADE,
    source_type     text NOT NULL CHECK (source_type IN ('github', 'zip')),
    repository_url  text NOT NULL DEFAULT '',
    commit_sha      text NOT NULL DEFAULT '',
    source_digest   text NOT NULL,
    workspace_path  text NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE job_files (
    job_id       uuid NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    path         text NOT NULL,
    size_bytes   bigint NOT NULL,
    sha256       text NOT NULL DEFAULT '',
    reason_code  text NOT NULL,
    eligible     boolean NOT NULL,
    selected     boolean NOT NULL,
    created_at   timestamptz NOT NULL DEFAULT now(),
    PRIMARY KEY (job_id, path),
    CHECK (NOT selected OR eligible)
);

INSERT INTO job_status_history (job_id, old_status, new_status, stage)
SELECT id, 'input_ready', 'queued', '' FROM jobs
WHERE status = 'input_ready' AND stage = '';

INSERT INTO job_progress_events (job_id, code, message)
SELECT id, 'source_queued', 'Source inspection queued.' FROM jobs
WHERE status = 'input_ready' AND stage = '';

UPDATE jobs SET status = 'queued', updated_at = now()
WHERE status = 'input_ready' AND stage = '';

-- +goose Down
-- Requeued legacy jobs keep their new status after rollback.
DROP TABLE job_files;
DROP TABLE job_sources;
