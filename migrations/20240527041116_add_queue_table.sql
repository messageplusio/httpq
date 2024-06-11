-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE jobs (
    id UUID PRIMARY KEY,
    job_type TEXT NOT NULL,
    job JSONB NOT NULL,
    retry_count INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL,
    delay_until TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    CONSTRAINT chk_delay_until CHECK (delay_until >= created_at)
);


CREATE INDEX idx_jobs_status ON jobs (status);
CREATE INDEX idx_jobs_created_at ON jobs (created_at);


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE jobs;