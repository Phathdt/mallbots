-- name: CreateUser :one
INSERT INTO users (
    email,
    password,
    full_name,
    created_at,
    updated_at
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;
