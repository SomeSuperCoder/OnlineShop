-- name: GetReviewsForProduct :many
SELECT * FROM reviews WHERE product = $1 ORDER BY created_at DESC;

-- name: InsertReview :one
INSERT INTO reviews
( product, comment, stars )
VALUES ( $1, $2, $3 )
RETURNING *;

-- name: DeleteReview :one
DELETE FROM reviews WHERE id = $1 RETURNING *;
