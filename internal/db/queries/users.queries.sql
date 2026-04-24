-- name: CreateUserOnSignup :one
INSERT INTO users(
    nationality,
    phone_number,
    name,
    email,
    password_hash
)
VALUES(
    $1, $2, $3, $4, $5
)
RETURNING *;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email = $1;