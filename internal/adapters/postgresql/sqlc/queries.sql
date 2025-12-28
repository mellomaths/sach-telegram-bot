-- name: CreateSAC :one
INSERT INTO sacs (
    user_id, message
) VALUES ($1, $2) RETURNING *;

-- name: CreateUser :one
INSERT INTO users (
    id, user_name, first_name, last_name
) VALUES ($1, $2, $3, $4) RETURNING *;

-- name: FindUserByID :one
SELECT * FROM users WHERE id = $1 LIMIT 1;
