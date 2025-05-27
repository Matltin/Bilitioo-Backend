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

-- name: GetUser :one
SELECT 
    "id", "email", "phone_number", "hashed_password"
FROM "user"
WHERE "email" = $1 OR "phone_number" = $2;

-- name: InitialProfile :exec
INSERT INTO "profile" (
  "user_id"
) VALUES (
  $1
);