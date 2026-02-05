-- name: InsertVote :one
INSERT INTO votes (review, voter, type)
VALUES ( $1, current_setting('app.user_id')::uuid, $2 )
RETURNING *;
