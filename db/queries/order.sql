-- name: FindAllOrders :many
SELECT * FROM orders ORDER BY created_at;

-- name: InsertOrder :one
INSERT INTO orders (details)
VALUES (
  $1
) RETURNING *;
