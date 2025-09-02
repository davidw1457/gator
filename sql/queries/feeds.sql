-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetFeeds :many
SELECT f.*, u.name AS user_name
FROM feeds AS f
    INNER JOIN users AS u ON f.user_id = u.id;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT i.*, f.name AS feed_name, u.name AS user_name
FROM inserted_feed_follow AS i
    INNER JOIN feeds AS f ON i.feed_id = f.id
    INNER JOIN users AS u ON i.user_id = u.id;

-- name: GetFeedByURL :one
SELECT *
FROM feeds
WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT f.name AS feed_name, u.name AS user_name
FROM feeds AS f
    INNER JOIN feed_follows AS ff on f.id = ff.feed_id
    INNER JOIN users AS u on ff.user_id = u.id
WHERE ff.user_id = $1;

-- name: DeleteFeedFollow :exec
DELETE
FROM feed_follows
WHERE feed_id = $1 AND user_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched_at = $1,
    updated_at = $2
WHERE id = $3;

-- name: GetNextFeedToFetch :one
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC NULLS FIRST;
