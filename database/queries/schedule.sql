-- name: ListSchedule :many
SELECT * FROM Schedule
WHERE id = "Scheduled" AND invocation_timestamp > NOW()
ORDER BY invocation_timestamp DESC
LIMIT $1;

-- name: IncrementFailure :one
UPDATE Schedule
SET 
  retries_no = retries_no + 1,
  failure_reason = $2,
  status = "Failed"
WHERE id = $1
RETURNING *;


-- name: ScheduleSuccss :one
UPDATE Schedule
SET 
  status = "Invoked"
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
  max_retries
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)
RETURNING *;