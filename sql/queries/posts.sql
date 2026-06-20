-- name: CreatePost :one
INSERT INTO posts(title, url, description, feed_id, published_at) VALUES ($1, $2, $3, $4, $5) RETURNING *;


-- name: GetPostsForUser :many
SELECT p.* FROM posts p JOIN feed f ON p.feed_id = f.id WHERE f.user_id = $1 ORDER BY p.published_at DESC LIMIT $2;