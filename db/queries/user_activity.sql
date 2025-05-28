-- name: CreateUserActivity :one
INSERT INTO "user_activity" (
    "user_id",
    "route_id",
    "vehicle_type"
) VALUES (
    $1, $2, $3
) RETURNING *;

-- name: UpdateUserActivity :one
UPDATE "user_activity"
SET 
    "status" = $1
WHERE id = $2
RETURNING *;
    
