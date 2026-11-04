-- +goose Up
CREATE TABLE job_sections (
 job_id uuid NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
 section_id text NOT NULL,
 position integer NOT NULL,
 record jsonb NOT NULL,
 PRIMARY KEY(job_id, section_id)
);
CREATE TABLE section_cache (cache_key text PRIMARY KEY, record jsonb NOT NULL);
CREATE TABLE provider_calls (
 id bigserial PRIMARY KEY,
 job_id uuid NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
 section_id text NOT NULL,
 provider text NOT NULL,
 model text NOT NULL,
 error_code text NOT NULL,
 input_tokens integer NOT NULL,
 output_tokens integer NOT NULL,
 created_at timestamptz NOT NULL DEFAULT now()
);
UPDATE jobs SET status='queued', stage='generating', attempts=0 WHERE status='input_ready' AND stage='evidence_confirmed';

-- +goose Down
DROP TABLE provider_calls;
DROP TABLE section_cache;
DROP TABLE job_sections;
