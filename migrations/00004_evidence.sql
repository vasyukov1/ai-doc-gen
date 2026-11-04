-- +goose Up
CREATE TABLE job_evidence (
    job_id uuid PRIMARY KEY REFERENCES jobs(id) ON DELETE CASCADE,
    record jsonb NOT NULL,
    confirmed_at timestamptz
);
UPDATE jobs SET status = 'queued', stage = 'analyzing', attempts = 0 WHERE status = 'input_ready' AND stage = 'source_selected';

-- +goose Down
DROP TABLE job_evidence;
