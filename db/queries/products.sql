-- name: FindAllProducts :many
SELECT * FROM products ORDER BY created_at;

-- name: InsertProduct :one
INSERT INTO products
  (name, details, price)
VALUES ( $1, $2, $3 )
RETURNING *;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1 LIMIT 1;

-- name: DeleteProduct :one
DELETE FROM products WHERE id = $1 RETURNING *;
