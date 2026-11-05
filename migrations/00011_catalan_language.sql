-- +goose Up
ALTER TABLE jobs DROP CONSTRAINT jobs_language_check;
ALTER TABLE jobs ADD CONSTRAINT jobs_language_check CHECK (language IN ('en', 'ru', 'ca'));

-- +goose Down
-- Refuse rollback while Catalan jobs exist rather than changing their document language.
ALTER TABLE jobs DROP CONSTRAINT jobs_language_check;
ALTER TABLE jobs ADD CONSTRAINT jobs_language_check CHECK (language IN ('en', 'ru'));
