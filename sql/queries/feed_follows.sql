-- name: AddFeedFollow :one
WITH insert_feed_follow AS (
        INSERT INTO feed_follows (user_id, feed_id) VALUES ($1, $2) RETURNING *
)
SELECT ff.*, u.name AS user_name, f.name AS feed_name FROM insert_feed_follow ff JOIN users u ON ff.user_id = u.id JOIN feed f ON f.id = ff.feed_id;


-- name: GetFeedFollowsByUserId :many
SELECT
    ff.id, ff.user_id AS user_id, ff.feed_id AS feed_id, u.name AS user_name, f.name AS feed_name
FROM feed_follows ff JOIN users u ON ff.user_id = u.id JOIN feed f ON ff.feed_id = f.id WHERE u.id = $1;