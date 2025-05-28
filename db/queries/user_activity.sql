-- name: CreateUserActivity :one
INSERT INTO "user_activity" (
    "user_id",
    "route_id",
    "vehicle_type"
) VALUES (
    $1, $2, $3
) RETURNING *;
    
