-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetUsers :many
SELECT name FROM users;

-- name: CheckUser :one
SELECT name FROM users WHERE name = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;

