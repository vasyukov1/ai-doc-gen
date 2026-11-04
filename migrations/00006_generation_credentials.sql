-- +goose Up
CREATE TABLE job_generation_credentials (
 job_id uuid PRIMARY KEY REFERENCES jobs(id) ON DELETE CASCADE,
 model text NOT NULL,
 encrypted_key text NOT NULL
);
-- +goose Down
DROP TABLE job_generation_credentials;
