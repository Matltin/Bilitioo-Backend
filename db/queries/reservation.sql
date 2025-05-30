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

-- name: GetIDReservation :many
SELECT r.id FROM "reservation" r
INNER JOIN "payment" p ON p.id = r.payment_id
WHERE p.id = $1;

-- name: GetReservationStatus :one 
SELECT status FROM "reservation" 
WHERE id = $1;

-- name: GetReservationDetails :one
SELECT
    r.id,
    r.payment_id,
    t.amount, 
    r.user_id,
    t.status
FROM "ticket" t
INNER JOIN "reservation" r ON r.ticket_id = t.id 
WHERE t.id = $1;

-- name: CancelReservation :exec
UPDATE "ticket"
SET status = 'CANCELLED'
WHERE id = $1 AND status = 'RESERVED';
