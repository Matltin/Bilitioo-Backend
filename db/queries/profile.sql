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