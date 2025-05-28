-- name: CreateUserActivity :one
INSERT INTO "user_activity" (
    "user_id",
    "route_id",
    "vehicle_type"
) VALUES (
    $1, $2, $3
) RETURNING 
    id, 
    user_id, 
    route_id, 
    vehicle_type, 
    status, 
    EXTRACT(EPOCH FROM duration_time)::bigint as duration_time_seconds,
    created_at;

-- name: UpdateUserActivity :one
UPDATE "user_activity"
SET 
    "status" = $1
WHERE id = $2
RETURNING 
    id, 
    user_id, 
    route_id, 
    vehicle_type, 
    status, 
    EXTRACT(EPOCH FROM duration_time)::bigint as duration_time_seconds,
    created_at;