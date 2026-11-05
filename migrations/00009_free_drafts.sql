-- +goose Up
CREATE TABLE job_draft_runs (
 job_id uuid PRIMARY KEY REFERENCES jobs(id) ON DELETE CASCADE,
 encrypted_connections text NOT NULL,
 deadline timestamptz NOT NULL DEFAULT (now() + interval '240 seconds'),
 attempts integer NOT NULL DEFAULT 0 CHECK (attempts BETWEEN 0 AND 12),
 reason text NOT NULL DEFAULT '',
 disabled jsonb NOT NULL DEFAULT '{}',
 section_attempts jsonb NOT NULL DEFAULT '{}'
);
CREATE TABLE provider_limits (
 scope text PRIMARY KEY,
 provider text NOT NULL,
 window_start timestamptz NOT NULL DEFAULT now(),
 requests integer NOT NULL DEFAULT 0,
 tokens integer NOT NULL DEFAULT 0,
 next_at timestamptz NOT NULL DEFAULT now(),
 active_until timestamptz NOT NULL DEFAULT now(),
 holder text NOT NULL DEFAULT ''
);
ALTER TABLE provider_calls ADD COLUMN duration_ms bigint NOT NULL DEFAULT 0;
-- +goose Down
ALTER TABLE provider_calls DROP COLUMN duration_ms;
DROP TABLE provider_limits;
DROP TABLE job_draft_runs;
