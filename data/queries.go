package data

const createEventTable = `-- create table for heroku
CREATE TABLE IF NOT EXISTS event (
    id SERIAL PRIMARY KEY,
    what TEXT,
    loc TEXT,
    "when" TIMESTAMP
)
`

const createEvent = `-- name: CreateEvent :one
INSERT INTO event (
  what, loc, "when"
) VALUES (
  $1, $2, $3
)
RETURNING id, what, loc, "when"
`

const deleteEvent = `-- name: DeleteEvent :exec
DELETE FROM event
WHERE id = $1
`

const getEvent = `-- name: GetEvent :one
SELECT id, what, loc, "when" FROM event
WHERE id = $1 LIMIT 1
`

const listEvents = `-- name: ListEvents :many
SELECT id, what, loc, "when" FROM event
ORDER BY "when" ASC
`

const updateEvent = `-- name: UpdateEvent :one
UPDATE event
SET what = $1,
    loc = $2,
    "when" = $3
WHERE id = $4
RETURNING *
`
