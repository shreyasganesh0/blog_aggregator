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

-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: FetchUserId :one
SELECT id FROM users WHERE name = $1;

-- name: FetchUserFeed :many
SELECT * FROM feeds WHERE user_id = $1;

-- name: FetchEntireFeed :many
SELECT f.name, f.url, u.name
FROM feeds as f
JOIN users as u
ON u.id = f.user_id;
