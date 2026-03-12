-- name: CreateUserAuth :execresult
INSERT INTO users_auth (login, password_hash) VALUES (?, ?);

-- name: CreateUser :exec
INSERT INTO users (user_id, name, role) VALUES (?, ?, ?);

-- name: GetUserIdentitiesById :one
SELECT 
    a.id, 
    a.login, 
    a.password_hash,
    u.name, 
    u.role
FROM users_auth a
INNER JOIN users u ON a.id = u.user_id
WHERE a.id = ?
LIMIT 1;

-- name: GetUserIdentitiesByLogin :one
SELECT 
    a.id, 
    a.login, 
    a.password_hash,
    u.name, 
    u.role
FROM users_auth a
INNER JOIN users u ON a.id = u.user_id
WHERE a.login = ?
LIMIT 1;

-- name: AddRefreshToken :exec
INSERT INTO refresh_tokens (token, user_id, expired_at, revoked, used)
VALUES (?, ?, ?, ?, ?);

-- name: GetRefreshToken :one
SELECT token, user_id, expired_at, revoked, used
FROM refresh_tokens
WHERE token = ?;

-- name: UpdateRefreshToken :execresult
UPDATE refresh_tokens 
SET 
    expired_at = ?, 
    used = ?, 
    revoked = ?
WHERE token = ?;

-- name: DeleteUserSessions :execresult
DELETE FROM refresh_tokens 
WHERE user_id = ?;

-- name: CleanupUselessTokens :exec
DELETE FROM refresh_tokens 
WHERE expired_at < NOW() 
   OR (used = 1 AND created_at < DATE_SUB(NOW(), INTERVAL 1 DAY))
   OR (revoked = 1 AND created_at < DATE_SUB(NOW(), INTERVAL 1 DAY));
