-- name: IsProductOwner :one
SELECT current_setting('app.user_id')::uuid = (
  SELECT owner FROM products WHERE id = $1
);

-- name: FindProductsPaged :many
SELECT * FROM products ORDER BY created_at DESC
LIMIT $1 OFFSET $2;

-- name: SearchForProducts :many
SELECT *, ts_rank(search_vector, to_tsquery($1)) as relevance
FROM products
WHERE search_vector @@ to_tsquery($1)
ORDER BY relevance DESC;

-- name: InsertProduct :one
INSERT INTO products
  (name, details, price, owner)
VALUES ( $1, $2, $3, current_setting('app.user_id')::uuid)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET
  name = coalesce(sqlc.narg('name'), name),
  details = coalesce(sqlc.narg('details'), details),
  price = coalesce(sqlc.narg('price'), price)
WHERE id = sqlc.arg('id')
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: DeleteProduct :one
DELETE FROM products WHERE id = $1 RETURNING *;
