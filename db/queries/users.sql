-- name: Is :one
SELECT current_setting('app.user_id')::uuid = sqlc.arg('user_id')::uuid;

-- name: SetConfig :one
SELECT set_config('app.user_id', sqlc.arg('user_id'), false);

-- name: UnsafeGetUserByEmail :one
SELECT * FROM users WHERE email = $1;

-- name: UnsafeGetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: InsertUser :one
INSERT INTO users (
	email, username, name, password_hash
) VALUES ( $1, $2, $3, crypt($4, gen_salt('bf')) )
RETURNING id, email, username, name, balance;

-- name: VerifyAuth :one
SELECT EXISTS (
  SELECT 1
  FROM users
  WHERE email = $1
  AND password_hash = crypt($2, password_hash)
) AS valid;

-- name: UpdateUserInfo :one
UPDATE users
SET
  email = coalesce(sqlc.narg('email'), email),
  username = coalesce(sqlc.narg('username'), username),
  name = coalesce(sqlc.narg('name'), name)
WHERE id = $1
RETURNING id, email, username, name, balance;

-- name: DeleteUser :one
DELETE FROM users
WHERE id = $1
RETURNING id, email, username, role, name, balance, created_at;
