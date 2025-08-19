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
    t.departure_time,
    t.status
FROM "ticket" t
INNER JOIN "reservation" r ON r.ticket_id = t.id 
WHERE r.id = $1;

-- name: CancelReservation :exec
UPDATE "ticket"
SET status = 'NOT_RESERVED'
WHERE id = $1 AND status = 'RESERVED';


-- name: GetCompletedUserReservation :many
SELECT 
    re.id AS "reservation_id",
    t.id AS "ticket_id",
    oc.province,
    dc.province
FROM "reservation" re 
INNER JOIN "ticket" t ON re.ticket_id = t.id
INNER JOIN "route" ro ON t.route_id = ro.id
INNER JOIN "city" oc ON oc.id = ro.origin_city_id
INNER JOIN "city" dc ON dc.id = ro.destination_city_id
WHERE re.status = 'RESERVED' AND re.user_id = $1;

-- name: GetAllUserReservation :many
SELECT
    re.id,
    t.id as ticket_id,
    oc.province as origin_province,
    dc.province as destination_province,
    re.status,
    re.payment_id,
    p.amount
FROM "reservation" re
INNER JOIN "ticket" t ON re.ticket_id = t.id
INNER JOIN "route" ro ON t.route_id = ro.id
INNER JOIN "city" oc ON oc.id = ro.origin_city_id
INNER JOIN "city" dc ON dc.id = ro.destination_city_id
INNER JOIN "payment" p ON re.payment_id = p.id -- Join with payment table
WHERE re.user_id = $1;

-- name: MarkExpiredReservations :exec
UPDATE reservation
SET status = 'CANCELED-BY-TIME'
WHERE (created_at + duration_time) < now()
  AND status != 'RESERVED';
