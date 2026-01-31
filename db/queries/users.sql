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
