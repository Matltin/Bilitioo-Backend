-- name: UpdateProfile :one
UPDATE "profile"
SET 
    "pic_dir" = COALESCE(sqlc.narg(pic_dir), "pic_dir"),
    "first_name" = COALESCE(sqlc.narg(first_name), "first_name"),
    "last_name" = COALESCE(sqlc.narg(last_name), "last_name"),
    "city_id" = COALESCE(sqlc.narg(city_id), "city_id"),
    "national_code" = COALESCE(sqlc.narg(national_code), "national_code")
WHERE "user_id" = $1
RETURNING *;


-- name: GetUserProfile :one
SELECT 
  u.id AS user_id,
  u.email,
  u.phone_number,
  u.role,
  u.status,
  u.phone_verified,
  u.email_verified,
  u.created_at,
  p.pic_dir,
  p.first_name,
  p.last_name,
  p.city_id,
  p.national_code
FROM "user" u
JOIN "profile" p ON u.id = p.user_id
WHERE u.id = $1;

-- name: AddToUserWallet :exec
UPDATE "profile"
SET wallet = wallet + $1
WHERE user_id = $2;
