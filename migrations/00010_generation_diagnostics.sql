-- +goose Up
ALTER TABLE provider_calls ADD COLUMN diagnostics jsonb NOT NULL DEFAULT '{}';
ALTER TABLE provider_calls ADD COLUMN rejected_content text NOT NULL DEFAULT '';

-- +goose Down
ALTER TABLE provider_calls DROP COLUMN rejected_content;
ALTER TABLE provider_calls DROP COLUMN diagnostics;
