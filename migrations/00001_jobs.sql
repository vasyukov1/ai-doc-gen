-- +goose Up
CREATE TABLE jobs (
    id                  uuid            PRIMARY KEY,
    status              text            NOT NULL CHECK (status IN (
        'queued',
        'running',
        'input_ready',
        'needs_input',
        'retry_wait',
        'ready_for_review',
        'completed',
        'failed',
        'cancelled'
    )),
    stage               text            NOT NULL DEFAULT '',
    profile_id          text            NOT NULL,
    language            text            NOT NULL CHECK (language IN ('en', 'ru')),
    description         text            NOT NULL DEFAULT '',
    attempts            integer         NOT NULL DEFAULT 0,
    max_attempts        integer         NOT NULL DEFAULT 3,
    lease_owner         text,
    lease_until         timestamptz,
    last_heartbeat_at   timestamptz,
    next_retry_at       timestamptz,
    error_code          text            NOT NULL DEFAULT '',
    error_message       text            NOT NULL DEFAULT '',
    required_action     text            NOT NULL DEFAULT '',
    completed_sections  integer         NOT NULL DEFAULT 0,
    total_sections      integer         NOT NULL DEFAULT 0,
    created_at          timestamptz     NOT NULL DEFAULT now(),
    updated_at          timestamptz     NOT NULL DEFAULT now()
);

CREATE TABLE job_inputs (
    job_id          uuid    PRIMARY KEY REFERENCES jobs(id) ON DELETE CASCADE,
    source_type     text    NOT NULL CHECK (source_type IN ('github', 'zip')),
    github_url      text,
    upload_path     text,
    original_name   text    NOT NULL DEFAULT '',
    sha256          text    NOT NULL DEFAULT '',
    size_bytes      bigint  NOT NULL DEFAULT 0,

    CHECK ((github_url IS NOT NULL) <> (upload_path IS NOT NULL))
);

CREATE TABLE job_status_history (
    id          bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    job_id      uuid        NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    old_status  text,
    new_status  text        NOT NULL,
    stage       text        NOT NULL DEFAULT '',
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE job_progress_events (
    id          bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    job_id      uuid        NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    code        text        NOT NULL,
    message     text        NOT NULL,
    created_at  timestamptz NOT NULL DEFAULT now()
);

CREATE TABLE job_artifacts (
    id              bigint      GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    job_id          uuid        NOT NULL REFERENCES jobs(id) ON DELETE CASCADE,
    kind            text        NOT NULL,
    storage_path    text        NOT NULL,
    sha256          text        NOT NULL,
    created_at      timestamptz NOT NULL DEFAULT now()
);

CREATE INDEX jobs_claim_idx ON jobs (created_at)
    WHERE status IN ('queued', 'retry_wait', 'running');

CREATE INDEX job_progress_events_job_idx ON job_progress_events (job_id, id);

-- +goose Down
DROP TABLE job_artifacts;
DROP TABLE job_progress_events;
DROP TABLE job_status_history;
DROP TABLE job_inputs;
DROP TABLE jobs;
