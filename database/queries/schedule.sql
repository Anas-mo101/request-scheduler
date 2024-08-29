-- name: ListSchedule :many
SELECT * FROM Schedule
WHERE status = 'Scheduled'
ORDER BY invocation_timestamp ASC
LIMIT $1;

-- name: ListRegSchedule :many
SELECT * FROM Schedule
WHERE 
  ($3::text IS NULL OR invocation_timestamp = $3) AND
  ($4::text IS NULL OR request_method = $4) AND
  ($5::text IS NULL OR request_url = $5) AND
  ($6::int32 IS NULL OR max_retries = $6) AND
  ($7::text IS NULL OR request_body_type = $7) AND
  ($8::text IS NULL OR status = $8)
ORDER BY invocation_timestamp ASC
LIMIT $1 OFFSET $2;

-- name: IncrementFailure :one
UPDATE Schedule
SET 
  retries_no = retries_no + 1,
  failure_reason = $2,
  status = 'Failed'
WHERE id = $1
RETURNING *;

-- name: UpdateSchedule :one
UPDATE Schedule
SET 
  invocation_timestamp = $2, 
  request_method = $3, 
  request_url = $4, 
  request_body = $5, 
  request_header = $6,
  request_query = $7, 
  max_retries = $8,
  request_body_type = $9,
  status = $10
WHERE id = $1
RETURNING *;

-- name: GetSchedule :one
SELECT * FROM Schedule
WHERE id = $1;

-- name: DeletSchedule :one
DELETE FROM Schedule
WHERE id = $1
RETURNING *;

-- name: ScheduleSuccss :one
UPDATE Schedule
SET 
  status = 'Invoked'
WHERE id = $1
RETURNING *;

-- name: CreateSchedule :one
INSERT INTO Schedule (
  invocation_timestamp, 
  request_method, 
  request_url, 
  request_body, 
  request_header,
  request_query, 
  max_retries,
  request_body_type,
  status
) VALUES (
  $1, $2, $3, $4, $5, $6, $7, $8, 'Scheduled'
)
RETURNING *;