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
  (name, details, price)
VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: DeleteProduct :one
DELETE FROM products WHERE id = $1 RETURNING *;
