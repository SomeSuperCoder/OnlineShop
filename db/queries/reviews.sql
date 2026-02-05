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
SELECT
  r.*,
  coalesce(sum(CASE WHEN v.type = 'upvote' THEN 1 ELSE 0 END), 0) AS upvotes,
  coalesce(sum(CASE WHEN v.type = 'downvote' THEN 1 ELSE 0 END), 0) AS downvotes,
  coalesce(sum(CASE WHEN v.type = 'upvote' THEN 1 WHEN v.type = 'downvote' THEN -1 ELSE 0 END)) AS rating
FROM reviews r
LEFT JOIN votes v ON r.id = v.review
WHERE r.product = $1
GROUP BY r.id
ORDER BY r.created_at DESC
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
