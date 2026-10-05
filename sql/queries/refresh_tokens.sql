-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES (
    $1, NOW(), NOW(), $2, $3
)
RETURNING *;

-- name: GetUserFromRefreshToken :one
SELECT users.* FROM refresh_tokens JOIN users ON user_id = id 
WHERE refresh_tokens.token = $1
AND revoked_at IS NULL
AND expires_at > NOW();

-- name: RevokeToken :one
UPDATE refresh_tokens
SET revoked_at = NOW(), updated_at = NOW()
WHERE token = $1
RETURNING *;

