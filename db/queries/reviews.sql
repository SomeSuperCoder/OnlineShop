-- name: IsAuthor :one
SELECT EXISTS (
  SELECT 1 FROM reviews
  WHERE id = $1 AND author = current_setting('app.user_id')::UUID
);

-- name: IsOwnerOrAuthor :one
SELECT EXISTS (
  SELECT 1 FROM reviews r
  LEFT JOIN products p ON r.product = p.id
  WHERE r.id = $1 AND (
    current_setting('app.user_id')::UUID IN (
      r.author,
      p.owner
    )
  )
);

-- name: GetReviewsForProduct :many
SELECT * FROM reviews WHERE product = $1 ORDER BY created_at DESC
LIMIT $2 OFFSET $3;

-- name: InsertReview :one
INSERT INTO reviews
( product, comment, stars, author )
VALUES ( $1, $2, $3, current_setting('app.user_id')::UUID )
RETURNING *;

-- name: UpdateReview :one
UPDATE reviews
SET
  comment = coalesce(sqlc.narg('comment'), comment),
  stars = coalesce(sqlc.narg('stars'), stars)
WHERE id = $1
RETURNING *;

-- name: DeleteReview :one
DELETE FROM reviews WHERE id = $1 RETURNING *;
