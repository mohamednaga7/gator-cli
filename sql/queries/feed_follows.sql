-- name: AddFeedFollow :one
WITH insert_feed_follow AS (
        INSERT INTO feed_follows (user_id, feed_id) VALUES ($1, $2) RETURNING *
)
SELECT ff.*, u.name AS user_name, f.name AS feed_name FROM insert_feed_follow ff JOIN users u ON ff.user_id = u.id JOIN feed f ON f.id = ff.feed_id;