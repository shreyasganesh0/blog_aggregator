// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"time"

	"github.com/google/uuid"
)

const checkUser = `-- name: CheckUser :one
SELECT id, name FROM users WHERE name = $1
`

type CheckUserRow struct {
	ID   uuid.UUID
	Name string
}

func (q *Queries) CheckUser(ctx context.Context, name string) (CheckUserRow, error) {
	row := q.db.QueryRowContext(ctx, checkUser, name)
	var i CheckUserRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, created_at, updated_at, name, url, user_id
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	Url       string
	UserID    uuid.UUID
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
		arg.Url,
		arg.UserID,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const createFeedFollow = `-- name: CreateFeedFollow :many
WITH inserted_feed_follows AS(
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING id, created_at, updated_at, feed_id, user_id
)
SELECT inserted_feed_follows.id, inserted_feed_follows.created_at, inserted_feed_follows.updated_at, inserted_feed_follows.feed_id, inserted_feed_follows.user_id, feeds.name, users.name
FROM inserted_feed_follows
INNER JOIN users ON users.id = inserted_feed_follows.user_id
INNER JOIN feeds ON feeds.id = inserted_feed_follows.feed_id
`

type CreateFeedFollowParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

type CreateFeedFollowRow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	FeedID    uuid.UUID
	UserID    uuid.UUID
	Name      string
	Name_2    string
}

func (q *Queries) CreateFeedFollow(ctx context.Context, arg CreateFeedFollowParams) ([]CreateFeedFollowRow, error) {
	rows, err := q.db.QueryContext(ctx, createFeedFollow,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.FeedID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CreateFeedFollowRow
	for rows.Next() {
		var i CreateFeedFollowRow
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.FeedID,
			&i.UserID,
			&i.Name,
			&i.Name_2,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, created_at, updated_at, name
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}

const deleteAllUsers = `-- name: DeleteAllUsers :exec
DELETE FROM users
`

func (q *Queries) DeleteAllUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, deleteAllUsers)
	return err
}

const feedByUrl = `-- name: FeedByUrl :one
SELECT id FROM feeds
WHERE url = $1
`

func (q *Queries) FeedByUrl(ctx context.Context, url string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, feedByUrl, url)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const feedFollowByUser = `-- name: FeedFollowByUser :many
SELECT feeds.name AS feed_name, users.name AS user_name
FROM feed_follows
INNER JOIN feeds ON feed_follows.feed_id = feeds.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1
`

type FeedFollowByUserRow struct {
	FeedName string
	UserName string
}

func (q *Queries) FeedFollowByUser(ctx context.Context, name string) ([]FeedFollowByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, feedFollowByUser, name)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FeedFollowByUserRow
	for rows.Next() {
		var i FeedFollowByUserRow
		if err := rows.Scan(&i.FeedName, &i.UserName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchEntireFeed = `-- name: FetchEntireFeed :many
SELECT f.name, f.url, u.name
FROM feeds as f
JOIN users as u
ON u.id = f.user_id
`

type FetchEntireFeedRow struct {
	Name   string
	Url    string
	Name_2 string
}

func (q *Queries) FetchEntireFeed(ctx context.Context) ([]FetchEntireFeedRow, error) {
	rows, err := q.db.QueryContext(ctx, fetchEntireFeed)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []FetchEntireFeedRow
	for rows.Next() {
		var i FetchEntireFeedRow
		if err := rows.Scan(&i.Name, &i.Url, &i.Name_2); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchUserFeed = `-- name: FetchUserFeed :many
SELECT id, created_at, updated_at, name, url, user_id FROM feeds WHERE user_id = $1
`

func (q *Queries) FetchUserFeed(ctx context.Context, userID uuid.UUID) ([]Feed, error) {
	rows, err := q.db.QueryContext(ctx, fetchUserFeed, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Feed
	for rows.Next() {
		var i Feed
		if err := rows.Scan(
			&i.ID,
			&i.CreatedAt,
			&i.UpdatedAt,
			&i.Name,
			&i.Url,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const fetchUserId = `-- name: FetchUserId :one
SELECT id FROM users WHERE name = $1
`

func (q *Queries) FetchUserId(ctx context.Context, name string) (uuid.UUID, error) {
	row := q.db.QueryRowContext(ctx, fetchUserId, name)
	var id uuid.UUID
	err := row.Scan(&id)
	return id, err
}

const getUsers = `-- name: GetUsers :many
SELECT name FROM users
`

func (q *Queries) GetUsers(ctx context.Context) ([]string, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		items = append(items, name)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}
