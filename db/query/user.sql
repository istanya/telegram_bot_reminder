-- name: CreateEvent :one
INSERT INTO events (
  user_id,
  user_name,
  dt_reminder,
  bot_message_id,
  message,
  state,
  dt_created
) VALUES (
  $1, $2, $3, $4, $5, $6, $7
)RETURNING *;

-- name: GetEvent :one
SELECT * FROM events
WHERE bot_message_id = $1 LIMIT 1;

-- name: ListEventsByTimeAndState :many
SELECT * FROM events
WHERE dt_reminder <= $1 AND state = $2;

-- name: UpdateEventState :exec
UPDATE events
SET state = $2
WHERE id = $1;