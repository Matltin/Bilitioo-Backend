-- name: CreateUser :one
INSERT INTO "user" (
  "email",
  "phone_number",
  "hashed_password",
  "email_verified",
  "phone_verified"
) VALUES (
  $1, $2, $3, $4, $5
) RETURNING *;