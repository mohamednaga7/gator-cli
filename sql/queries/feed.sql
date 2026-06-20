-- name: CreateFeed :one
INSERT INTO feed(name, url, user_id) VALUES ($1, $2, $3) RETURNING *;

-- name: GetAllFeeds :many
SELECT f.id, f.name, f.url, u.name AS user_name FROM feed f JOIN users u ON u.id = f.user_id;

-- name: GetFeedByUrl :one
SELECT * FROM feed WHERE url = $1;