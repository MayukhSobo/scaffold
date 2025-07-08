-- name: GetAdminUsers :many
SELECT *
FROM users
WHERE role = 'admin';

-- name: GetPendingVerificationUsers :many
SELECT *
FROM users
WHERE status = 'pending_verification';

-- name: GetUsers :many
SELECT * FROM users;

-- name: GetUser :one
SELECT * FROM users
WHERE id = ?; 