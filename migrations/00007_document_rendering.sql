-- +goose Up
CREATE TABLE job_document_inputs (
 job_id uuid PRIMARY KEY REFERENCES jobs(id) ON DELETE CASCADE,
 record jsonb NOT NULL,
 capability_hash text NOT NULL DEFAULT '',
 reviewed boolean NOT NULL DEFAULT false,
 created_at timestamptz NOT NULL DEFAULT now()
);

-- +goose Down
DROP TABLE job_document_inputs;
