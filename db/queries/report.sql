-- name: GetReports :many
SELECT * FROM "report";

-- name: AnswerReport :one
UPDATE "report"
SET 
    response_text = $1,
    admin_id = $2
WHERE id = $3
RETURNING *;

-- name: CreateReport :one
INSERT INTO "report" (
    reservation_id,
    user_id,
    request_type,
    request_text,
    response_text
)
VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetUserReport :many
SELECT 
    r.id,
    r.user_id,
    r.reservation_id,
    r.request_type,
    r.request_text,
    r.response_text,
    r.admin_id,
    res.ticket_id,
    res.status AS reservation_status,
    t.departure_time,
    t.arrival_time,
    t.amount AS ticket_amount,
    t.vehicle_type,
    rt.origin_city_id,
    rt.destination_city_id,
    oc.province AS origin_city_name,
    dc.province AS destination_city_name
FROM "report" r
LEFT JOIN reservation res ON r.reservation_id = res.id
LEFT JOIN ticket t ON res.ticket_id = t.id
LEFT JOIN route rt ON t.route_id = rt.id
LEFT JOIN city oc ON rt.origin_city_id = oc.id
LEFT JOIN city dc ON rt.destination_city_id = dc.id
WHERE r.user_id = $1
ORDER BY r.id DESC;


-- name: GetUserReportSimple :many
SELECT 
    id,
    user_id,
    reservation_id,
    request_type,
    request_text,
    response_text,
    admin_id
FROM "report" 
WHERE user_id = $1
ORDER BY id DESC;

-- name: GetUserReportByStatus :many
SELECT 
    r.id,
    r.user_id,
    r.reservation_id,
    r.request_type,
    r.request_text,
    r.response_text,
    r.admin_id,
    CASE 
        WHEN r.response_text IS NULL OR r.response_text = '' THEN 'PENDING'
        ELSE 'ANSWERED'
    END as status
FROM "report" r
WHERE r.user_id = $1
  AND (
        ($2::text = 'PENDING' AND (r.response_text IS NULL OR r.response_text = ''))
     OR ($2::text = 'ANSWERED' AND (r.response_text IS NOT NULL AND r.response_text != ''))
     OR ($2::text NOT IN ('PENDING','ANSWERED'))
  )
ORDER BY r.id DESC;

-- name: GetUserReportByType :many
SELECT 
    id,
    user_id,
    reservation_id,
    request_type,
    request_text,
    response_text,
    admin_id
FROM "report" 
WHERE user_id = $1 
  AND request_type = $2
ORDER BY id DESC;
