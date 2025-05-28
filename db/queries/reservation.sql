-- name: CreateReservation :one
INSERT INTO "reservation" (
    "user_id",
    "ticket_id",
    "payment_id"
) VALUES (
    $1, $2, $3
) RETURNING 
    id, 
    user_id, 
    ticket_id, 
    payment_id, 
    status, 
    EXTRACT(EPOCH FROM duration_time)::bigint as duration_time_seconds,
    created_at;

-- name: UpdateReservation :one
UPDATE "reservation"
SET 
    "status" = $1
WHERE id = $2
RETURNING 
    id, 
    user_id, 
    ticket_id, 
    payment_id, 
    status, 
    EXTRACT(EPOCH FROM duration_time)::bigint as duration_time_seconds,
    created_at;

-- name: GetReservationStatus :one 
SELECT status FROM "reservation" 
WHERE id = $1;