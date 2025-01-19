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
SELECT id, name FROM users WHERE name = $1;

-- name: DeleteAllUsers :exec
DELETE FROM users;
DELETE FROM feeds;

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

-- name: FetchFeedId :one
SELECT id FROM feeds where url = $1;

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

-- name: DeleteFeedFollowByFUrlUser :exec
DELETE FROM feed_follows
WHERE feed_follows.user_id = (SELECT id FROM users WHERE users.name = $1)
and feed_follows.feed_id = (SELECT id FROM feeds WHERE feeds.url = $2);

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1, updated_at = $2
WHERE url = $3;

-- name: GetNextFeedToFetch :one
SELECT url FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST
LIMIT 1; 

-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8
)
ON CONFLICT (url) DO NOTHING
RETURNING *;

-- name: GetPostsByUser :many
SELECT posts.* FROM posts
INNER JOIN feed_follows
ON posts.feed_id = feed_follows.feed_id
INNER JOIN users
ON users.id = feed_follows.user_id
WHERE users.name = $1
ORDER BY published_at ASC
LIMIT $2;

