-- name: CreateUsers :one
INSERT INTO users (
  user_name, 
  hassed_password,
  full_name,
  email
) VALUES (
  $1, $2, $3, $4
)
RETURNING *;


-- name: GetUsers :one
SELECT * FROM users
WHERE user_name = $1 LIMIT 1;