-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: ListFeeds :many
SELECT * from feeds;

-- name: GetNextFeedsToFetch :many
SELECT * from feeds
ORDER BY latest_fetched_at ASC NULLS FIRST
LIMIT $1;

-- name: MakrFeedAsFetched :one
UPDATE feeds SET latest_fetched_at = NOW(), updated_at = NOW()
WHERE ID = $1
RETURNING *;
