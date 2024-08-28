-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TYPE "method" AS ENUM ('GET', 'POST');

-- Create the enum type for status
CREATE TYPE "status" AS ENUM ('Scheduled', 'Invoked', 'Failed');

CREATE TYPE "bodyType" AS ENUM ('TEXT', 'JSON');

-- Create the Schedule table using the enums
CREATE TABLE Schedule (
    id SERIAL PRIMARY KEY,
    invocation_timestamp TIMESTAMPTZ NOT NULL,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    request_method "method" NOT NULL,
    request_body_type "bodyType",
    request_body TEXT,
    request_url TEXT NOT NULL,
    request_header JSONB,
    request_query JSONB,
    status "status" NOT NULL,
    retries_no INT DEFAULT 0,
    max_retries INT DEFAULT 1,
    failure_reason TEXT
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE Schedule;
DROP TYPE "method";
DROP TYPE "bodyType";
DROP TYPE "status";
SELECT 'down SQL query';
-- +goose StatementEnd
