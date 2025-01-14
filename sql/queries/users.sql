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

-- name: FeedByUrl :one
SELECT id FROM feeds
WHERE url = $1;

-- name: FeedFollowByUser :many
SELECT feeds.name AS feed_name, users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1;

-- name: CreateFeedFollow :many
WITH inserted_feed_follows AS(
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *
)
SELECT inserted_feed_follows.*, feeds.name, users.name
FROM inserted_feed_follows
INNER JOIN users ON users.id = inserted_feed_follows.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follows.feed_id;


