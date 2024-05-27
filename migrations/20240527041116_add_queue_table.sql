-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied

CREATE TABLE queue (
    id SERIAL PRIMARY KEY,
    message TEXT NOT NULL,
    processed BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

DROP TABLE queue;
