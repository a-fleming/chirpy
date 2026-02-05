-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, user_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetRefreshTokenDetails :one
SELECT *
FROM refresh_tokens
WHERE token = $1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET updated_at = NOW(), revoked_at = NOW()
WHERE token = $1;
