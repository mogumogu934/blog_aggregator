-- name: AddFeed :one
INSERT INTO feeds (id, created_at, updated_at, last_fetched_at, name, url, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7
)
RETURNING *;

-- name: GetFeedFromURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: GetFeeds :many
SELECT name, url, user_id
FROM feeds;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = NOW(), updated_at = NOW()
WHERE id = $1;

-- name: GetNextFeedToFetch :one
SELECT * FROM feeds
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;