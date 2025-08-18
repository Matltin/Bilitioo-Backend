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
    "id", "email", "phone_number", "hashed_password", "email_verified", "phone_verified", "role"
FROM "user"
WHERE "email" = $1 OR "phone_number" = $2;

-- name: InitialProfile :exec
INSERT INTO "profile" (
  "user_id"
) VALUES (
  $1
);

-- name: GetUserByID :one
SELECT 
    "id", "email", "phone_number", "hashed_password", "role"
FROM "user"
WHERE "id" = $1;

-- name: UpdateUserContact :one
UPDATE "user"
SET
  email = COALESCE(sqlc.narg(email), email),
  phone_number = COALESCE(sqlc.narg(phone_number), phone_number),
  hashed_password = COALESCE(sqlc.narg(hashed_password), hashed_password)
WHERE id = $1
RETURNING *;

-- name: UpdateUserEmailVerified :exec
UPDATE "user"
SET email_verified = $2
WHERE id = $1;