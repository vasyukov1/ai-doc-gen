-- +goose Up
-- Legacy records remain private and are retained until an operator migrates them.
ALTER TABLE jobs ADD COLUMN owner_hash text;
ALTER TABLE jobs ADD COLUMN expires_at timestamptz;
ALTER TABLE jobs ADD COLUMN purge_after timestamptz;
ALTER TABLE jobs ADD COLUMN resume_status text;
ALTER TABLE section_cache ADD COLUMN created_at timestamptz NOT NULL DEFAULT now();
CREATE INDEX jobs_owner_idx ON jobs(owner_hash, created_at DESC);
CREATE INDEX jobs_expiry_idx ON jobs(expires_at) WHERE expires_at IS NOT NULL;

-- +goose Down
ALTER TABLE section_cache DROP COLUMN created_at;
DROP INDEX jobs_expiry_idx;
DROP INDEX jobs_owner_idx;
ALTER TABLE jobs DROP COLUMN resume_status;
ALTER TABLE jobs DROP COLUMN purge_after;
ALTER TABLE jobs DROP COLUMN expires_at;
ALTER TABLE jobs DROP COLUMN owner_hash;
